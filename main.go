package main

import (
	endpoint_http "capuchin/app/endpoint/http"
	"capuchin/app/middleware/session"
	"capuchin/app/repository/sqlite"
	"capuchin/app/util/cookiemanager"
	"capuchin/app/util/db"
	"embed"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/limiter"

	_ "modernc.org/sqlite"
)

const appName string = "Capuchin"
const dbPath string = "./capuchin.sqlite3"
const envAuth string = "CAPUCHIN_AUTH"
const envPwd string = "CAPUCHIN_PASSWORD"
const cookieKey string = "CAPUCHIN_SID"
const cookieExpiresDays uint16 = 365

// Embed a directory
//
//go:embed _ui/dist/*
var uiFS embed.FS

var withAuthPaths []string = []string{"/api/clients", "/api/projects", "/api/sessions", "/api/settings", "/api/timelogs"}

var errorHandler = func(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	msg := fiber.ErrInternalServerError.Message

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		msg = e.Message
	}

	fmt.Printf("[%s] %s: %s\n", c.Method(), c.OriginalURL(), err.Error())

	c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)

	return c.Status(code).SendString(msg)
}

func main() {
	fmt.Printf("===[ %s ]===\n", appName)

	var err error
	var pwd string = ""

	debugMode := flag.Bool("debug", false, "Debug mode")
	port := flag.Int64("p", 8090, "Port")
	flag.Parse()

	fmt.Printf("🐵 debugMode=%t\n", *debugMode)

	auth, ok := os.LookupEnv(envAuth)
	if ok && auth == "1" {
		fmt.Printf("🐵 %s=1\n", envAuth)

		pwd = os.Getenv(envPwd)
		if pwd == "" {
			log.Fatalf("❌ [init] \"%s\" is empty\n", envPwd)
		}
	} else {
		fmt.Printf("🐵 %s=0\n", envAuth)
	}

	db, err := db.New(dbPath)
	if err != nil {
		log.Fatal("❌ [init] Database connection failed\n", err)
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		s := <-ch
		fmt.Printf("👽 Received signal: %s\n", s)
		time.Sleep(time.Second)
		fmt.Println("💥 Server shutdown...")
		os.Exit(1)
	}()

	validate := validator.New()
	cm := cookiemanager.New(cookieKey, cookieExpiresDays)

	cr := sqlite.NewClientRepository(db)
	pr := sqlite.NewProjectRepository(db)
	tr := sqlite.NewTimelogRepository(db)
	sr := sqlite.NewSettingRepository(db)
	ssr := sqlite.NewSessionRepository(db)

	app := fiber.New(fiber.Config{
		ErrorHandler: errorHandler,
	})

	app.Use(cors.New())

	app.Use(session.New(session.Config{
		Paths:             withAuthPaths,
		CookieManager:     cm,
		SessionRepository: ssr,
	}, debugMode, auth))

	app.Use("/api/auth/login", limiter.New(limiter.Config{
		Max:        20,
		Expiration: 30 * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusTooManyRequests)
		},
		LimiterMiddleware: limiter.SlidingWindow{},
	}))

	api := app.Group("/api")

	endpoint_http.RegisterIndexEndpoints(app, cm, ssr, auth, debugMode)
	endpoint_http.RegisterClientEndpoints(api, validate, *cr)
	endpoint_http.RegisterProjectEndpoints(api, validate, *pr)
	endpoint_http.RegisterTimelogEndpoints(api, validate, *tr, *pr)
	endpoint_http.RegisterSettingEndpoints(api, validate, *sr)

	if auth == "1" {
		endpoint_http.RegisterAuthEndpoints(api, validate, *ssr, cm, pwd)
		endpoint_http.RegisterSessionEndpoints(api, validate, *ssr, cm, debugMode)
	}

	// Обязательно после всех остальных маршрутов
	app.Use("/", filesystem.New(filesystem.Config{
		Root:       http.FS(uiFS),
		PathPrefix: "/_ui/dist",
		MaxAge:     3600,
	}))

	fmt.Printf("🌍 Server started: http://127.0.0.1:%d\n", *port)
	startServer(app, port)
}
