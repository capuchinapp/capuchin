package cloudbackend

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/contrib/otelfiber"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/hashicorp/golang-lru/v2/expirable"
	"github.com/wneessen/go-mail"
	"go.uber.org/zap"

	"capuchin/internal/cookiemanager"
	"capuchin/internal/domain"
	"capuchin/internal/mailer"
	"capuchin/internal/metrics"
	"capuchin/internal/middleware"
	"capuchin/internal/restapi"
	"capuchin/internal/sanitizer"
	"capuchin/internal/service"
	"capuchin/internal/sessionupdater"
	"capuchin/internal/sqlite"
	"capuchin/internal/validator"
)

const (
	defaultIdleTimeout = 5 * time.Second

	defaultLimiterMax        = 50
	defaultLimiterExpiration = 1 * time.Second

	defaultSessionCacheSize = 100
	defaultSessionCacheTTL  = 10 * time.Minute
)

// Start запускает приложение.
func Start(appVersion string) error { //nolint:gocognit,maintidx // Всё в порядке
	conf, err := LoadConfiguration()
	if err != nil {
		return fmt.Errorf("load configuration: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := newDB(ctx, conf)
	if err != nil {
		return fmt.Errorf("create database connection: %v", err)
	}

	logger, err := newLogger(conf.Log.Level)
	if err != nil {
		return fmt.Errorf("create logger: %v", err)
	}
	defer logger.Sync() //nolint:errcheck // все в порядке

	mailClient, err := mail.NewClient(
		conf.Mailer.Host,
		mail.WithPort(conf.Mailer.Port),
		mail.WithSSL(),
		mail.WithTLSPortPolicy(mail.TLSMandatory),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(conf.Mailer.Username),
		mail.WithPassword(conf.Mailer.Password),
	)
	if err != nil {
		return fmt.Errorf("create mail client: %v", err)
	}

	tp := newTracer(ctx)
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	app := fiber.New(fiber.Config{
		Immutable:               true,
		IdleTimeout:             defaultIdleTimeout,
		EnableTrustedProxyCheck: true,
		TrustedProxies:          conf.HTTP.TrustedProxies,
		ProxyHeader:             fiber.HeaderXForwardedFor,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			msg := fiber.ErrInternalServerError.Message

			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
				msg = e.Message
			}

			if code >= fiber.StatusInternalServerError {
				logger.Error("Internal error",
					zap.String("method", c.Method()),
					zap.String("url", c.OriginalURL()),
					zap.Error(err),
				)
			} else {
				logger.Warn("Client error",
					zap.String("method", c.Method()),
					zap.String("url", c.OriginalURL()),
					zap.Error(err),
				)
			}

			c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)

			err = c.Status(code).SendString(msg)
			if err != nil {
				return fmt.Errorf("send string: %v", err)
			}

			return nil
		},
	})

	app.Use(recover.New())

	// Health
	app.Get("/api/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	app.Use(otelfiber.Middleware(
		otelfiber.WithTracerProvider(tp),
	))

	if conf.HTTP.CORS.Enabled {
		corsConfig := cors.Config{
			AllowOrigins:     conf.HTTP.CORS.AllowOrigins,
			AllowMethods:     conf.HTTP.CORS.AllowMethods,
			AllowCredentials: conf.HTTP.CORS.AllowCredentials,
		}
		corsMiddleware := cors.New(corsConfig)

		app.Use(corsMiddleware)
	}

	cookieMan := cookiemanager.New(
		conf.HTTP.Cookie.Name,
		conf.HTTP.Cookie.ExpiresDays,
		conf.HTTP.Cookie.Secure,
		conf.HTTP.Cookie.Domain,
	)
	sessionRepo := sqlite.NewSessionRepository(db)
	sessionUpdater := sessionupdater.New(sessionRepo, logger)

	app.Use(middleware.NewAuthentication(middleware.AuthenticationConfig{
		CookieManager:  cookieMan,
		SessionRepo:    sessionRepo,
		SessionCache:   expirable.NewLRU[string, domain.Session](defaultSessionCacheSize, nil, defaultSessionCacheTTL),
		SessionUpdater: sessionUpdater,
		Logger:         logger.Named("authentication"),
		NowFunc:        time.Now,
		Paths: []string{
			"/api/clients",
			"/api/projects",
			"/api/sessions",
			"/api/settings",
			"/api/timelogs",
			"/api/tasks",
			"/api/favorites",
		},
	}))

	metricsInstance := metrics.New(logger)
	app.Use(middleware.NewRequestLog(logger.Named("request_log"), metricsInstance))

	mailService := mailer.New(conf.Mailer.From, mailClient)
	validate := validator.New()
	sanitize := sanitizer.New()

	clientRepo := sqlite.NewClientRepository(db)
	projectRepo := sqlite.NewProjectRepository(db)
	timelogRepo := sqlite.NewTimelogRepository(db)
	taskRepo := sqlite.NewTaskRepository(db)
	favoriteRepo := sqlite.NewFavoriteRepository(db)
	settingRepo := sqlite.NewSettingRepository(db)
	userRepo := sqlite.NewUserRepository(db)
	authCodeRepo := sqlite.NewAuthCodeRepository(db)
	auditLogRepo := sqlite.NewAuditLogRepository()

	txManager := sqlite.NewTransactionManager(db)

	authService := service.NewAuthService(mailService, authCodeRepo, sessionRepo, userRepo, txManager, logger)
	sessionService := service.NewSessionService(sessionRepo)
	clientService := service.NewClientService(clientRepo, auditLogRepo, txManager)
	projectService := service.NewProjectService(projectRepo, clientRepo, auditLogRepo, txManager)
	timelogService := service.NewTimelogService(timelogRepo, clientRepo, projectRepo, taskRepo, auditLogRepo, txManager)
	taskService := service.NewTaskService(taskRepo, clientRepo, projectRepo, timelogRepo, auditLogRepo, txManager)
	favoriteService := service.NewFavoriteService(favoriteRepo, clientRepo, projectRepo, taskRepo, auditLogRepo, txManager)

	api := app.Group("/api")

	restapi.RegisterAuthEndpoints(api, validate, authService, cookieMan, sessionRepo, logger)

	app.Use(limiter.New(limiter.Config{
		Max:        defaultLimiterMax,
		Expiration: defaultLimiterExpiration,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusTooManyRequests)
		},
	}))

	restapi.RegisterIndexEndpoints(api, appVersion, cookieMan, sessionRepo, timelogRepo, logger)
	restapi.RegisterClientEndpoints(api, validate, sanitize, clientService, clientRepo, logger)
	restapi.RegisterProjectEndpoints(api, validate, sanitize, projectService, projectRepo, logger)
	restapi.RegisterTimelogEndpoints(api, validate, sanitize, timelogService, timelogRepo, logger)
	restapi.RegisterTaskEndpoints(api, validate, sanitize, taskService, taskRepo, logger)
	restapi.RegisterFavoriteEndpoints(api, validate, sanitize, favoriteService, favoriteRepo, logger)
	restapi.RegisterSettingEndpoints(api, validate, sanitize, settingRepo, logger)
	restapi.RegisterSessionEndpoints(api, sessionService, sessionRepo, cookieMan)

	// Статика фронтенда из embed FS с SPA-fallback. API и метрики раздаются
	// отдельными серверами, поэтому исключаются из статического обработчика.
	app.Use(staticHandler(webDist(), []string{"/api", "/metrics"}))

	// 404 Handler
	app.Use(func(c *fiber.Ctx) error {
		c.Append("X-RouteNotFound", "yes")

		return c.SendStatus(http.StatusNotFound)
	})

	logger.Info("Start capuchin API service")

	// Слушаем порт приложения
	appAddr := fmt.Sprintf(":%d", conf.HTTP.Port)
	go func() {
		if err := app.Listen(appAddr); err != nil {
			logger.Error("App listen",
				zap.Error(err),
			)
			panic(err)
		}
	}()
	logger.Info("Server started",
		zap.String("address", appAddr),
	)

	startMetrics(conf.Metrics.Port, db.DB, conf.Sqlite.DBPath, metricsInstance, logger)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c // Блокируем основной поток до тех пор, пока не будет получено прерывание
	logger.Info("Gracefully shutting down...")
	err = app.Shutdown()
	if err != nil {
		logger.Error("Shutdown",
			zap.Error(err),
		)
		panic(err)
	}

	logger.Info("Running cleanup tasks...")

	// Здесь задачи по очистке
	sessionUpdater.Stop()

	logger.Info("Stop capuchin API service")

	logger.Info("Capuchin was successful shutdown.")

	return nil
}
