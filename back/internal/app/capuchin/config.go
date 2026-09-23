package capuchin

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"go.uber.org/zap/zapcore"
)

// Configuration представляет конфигурацию приложения.
type Configuration struct {
	// Environment - среда выполнения.
	// Возможные значения: "local", "stage", "prod".
	Environment string `env:"ENVIRONMENT" envDefault:"prod"`

	Log     Log
	Metrics Metrics
	HTTP    HTTPServer
	Sqlite  Sqlite
	Mailer  Mailer
}

// Log конфигурация логирования.
type Log struct {
	// Уровень логирования.
	// Возможные значения: "debug", "info", "warn", "error".
	Level zapcore.Level `env:"LOG_LEVEL" envDefault:"info"`
}

// Metrics конфигурация метрик.
type Metrics struct {
	Port int `env:"METRICS_PORT" envDefault:"8080"`
}

// HTTPServer конфигурация HTTP сервера.
type HTTPServer struct {
	Port   int `env:"HTTP_PORT" envDefault:"3000"`
	CORS   HTTPCORS
	Cookie Cookie

	// TrustedProxies - список доверенных прокси-серверов.
	// Пример: "127.0.0.1,172.25.0.0/16" - localhost и диапазон Docker сети
	TrustedProxies []string `env:"HTTP_TRUSTED_PROXIES" envDefault:"127.0.0.1,172.25.0.0/16"`
}

// HTTPCORS конфигурация CORS для HTTP сервера.
type HTTPCORS struct {
	Enabled bool `env:"HTTP_CORS_ENABLED" envDefault:"true"`

	// AllowOrigins - разрешенные домены.
	// Пример: "http://localhost:5173", "https://app.capuchin.ru"
	AllowOrigins string `env:"HTTP_CORS_ALLOW_ORIGINS,required"`

	// AllowMethods - разрешенные HTTP-методы.
	AllowMethods string `env:"HTTP_CORS_ALLOW_METHODS" envDefault:"GET,POST,HEAD,PUT,DELETE,PATCH"`

	// AllowCredentials - разрешены куки.
	AllowCredentials bool `env:"HTTP_CORS_ALLOW_CREDENTIALS" envDefault:"true"`
}

// Cookie конфигурация cookie для HTTP сервера.
type Cookie struct {
	// Name - имя куки.
	Name string `env:"COOKIE_NAME" envDefault:"CAPUCHIN_SID"`

	// ExpiresDays - количество дней до истечения срока действия куки.
	ExpiresDays int `env:"COOKIE_EXPIRES_DAYS" envDefault:"365"`

	// Secure - куки доступны только через HTTPS.
	Secure bool `env:"COOKIE_SECURE" envDefault:"true"`

	// Domain - домен, на котором куки доступны.
	// Пример: "localhost", "app.capuchin.ru"
	Domain string `env:"COOKIE_DOMAIN,required"`
}

// Sqlite конфигурация файлового хранилища SQLite.
type Sqlite struct {
	// DBPath - путь к файлу базы данных SQLite.
	DBPath string `env:"SQLITE_DB_PATH" envDefault:"capuchin.db"`
}

// Mailer конфигурация почтовика.
// Поддерживается только TLS.
type Mailer struct {
	// Host - хост почтового сервера.
	// Пример: "smtp.timeweb.ru"
	Host string `env:"MAILER_HOST,required"`

	// Port - порт почтового сервера.
	// Пример: 465
	Port int `env:"MAILER_PORT,required"`

	// Username - имя пользователя почтового сервера.
	// Пример: "no-reply@capuchin.ru"
	Username string `env:"MAILER_USERNAME,required"`

	// Password - пароль пользователя почтового сервера.
	Password string `env:"MAILER_PASSWORD,required"`

	// From - адрес отправителя.
	// Пример: "Capuchin <no-reply@capuchin.ru>"
	From string `env:"MAILER_FROM,required"`
}

// LoadConfiguration возвращает новую конфигурацию приложения на основе переменных среды.
func LoadConfiguration() (*Configuration, error) {
	var config Configuration
	if err := env.Parse(&config); err != nil {
		return nil, fmt.Errorf("parse configuration: %v", err)
	}

	return &config, nil
}
