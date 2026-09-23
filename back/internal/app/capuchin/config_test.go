package capuchin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var allEnvList = []string{ //nolint:gochecknoglobals // it's ok
	"ENVIRONMENT",
	"LOG_LEVEL",
	"METRICS_PORT",
	"HTTP_PORT",
	"HTTP_TRUSTED_PROXIES",
	"HTTP_CORS_ENABLED",
	"HTTP_CORS_ALLOW_ORIGINS",
	"HTTP_CORS_ALLOW_METHODS",
	"HTTP_CORS_ALLOW_CREDENTIALS",
	"COOKIE_NAME",
	"COOKIE_EXPIRES_DAYS",
	"COOKIE_SECURE",
	"COOKIE_DOMAIN",
	"SQLITE_DB_PATH",
	"MAILER_HOST",
	"MAILER_PORT",
	"MAILER_USERNAME",
	"MAILER_PASSWORD",
	"MAILER_FROM",
}

func TestLoadConfiguration(t *testing.T) {
	tests := []struct {
		name       string
		setEnvFunc func()
		want       *Configuration
		wantErr    error
	}{
		{
			name: "Should load configuration from env",
			setEnvFunc: func() {
				t.Setenv("ENVIRONMENT", "local")
				t.Setenv("LOG_LEVEL", "debug")

				t.Setenv("METRICS_PORT", "9090")

				t.Setenv("HTTP_PORT", "3001")

				t.Setenv("HTTP_CORS_ENABLED", "false")
				t.Setenv("HTTP_CORS_ALLOW_ORIGINS", "http://localhost:5173")
				t.Setenv("HTTP_CORS_ALLOW_METHODS", "GET,POST")
				t.Setenv("HTTP_CORS_ALLOW_CREDENTIALS", "false")

				t.Setenv("COOKIE_NAME", "CAPUCHIN_SID_TEST")
				t.Setenv("COOKIE_EXPIRES_DAYS", "100")
				t.Setenv("COOKIE_SECURE", "false")
				t.Setenv("COOKIE_DOMAIN", "localhost")

				t.Setenv("HTTP_TRUSTED_PROXIES", "127.0.0.1,192.168.0.0/16")

				t.Setenv("SQLITE_DB_PATH", "/path/to/capuchin.db")

				t.Setenv("MAILER_HOST", "localhost")
				t.Setenv("MAILER_PORT", "25")
				t.Setenv("MAILER_USERNAME", "username")
				t.Setenv("MAILER_PASSWORD", "password")
				t.Setenv("MAILER_FROM", "from@localhost.tld")
			},
			want: &Configuration{
				Environment: "local",
				Log: Log{
					Level: -1,
				},
				Metrics: Metrics{
					Port: 9090,
				},
				HTTP: HTTPServer{
					Port: 3001,
					CORS: HTTPCORS{
						Enabled:          false,
						AllowOrigins:     "http://localhost:5173",
						AllowMethods:     "GET,POST",
						AllowCredentials: false,
					},
					Cookie: Cookie{
						Name:        "CAPUCHIN_SID_TEST",
						ExpiresDays: 100,
						Secure:      false,
						Domain:      "localhost",
					},
					TrustedProxies: []string{
						"127.0.0.1",
						"192.168.0.0/16",
					},
				},
				Sqlite: Sqlite{
					DBPath: "/path/to/capuchin.db",
				},
				Mailer: Mailer{
					Host:     "localhost",
					Port:     25,
					Username: "username",
					Password: "password",
					From:     "from@localhost.tld",
				},
			},
			wantErr: nil,
		},
		{
			name:       "Should load default configuration",
			setEnvFunc: func() {},
			want: &Configuration{
				Environment: "prod",
				Metrics: Metrics{
					Port: 8080,
				},
				HTTP: HTTPServer{
					Port: 3000,
					CORS: HTTPCORS{
						Enabled:          true,
						AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH",
						AllowCredentials: true,
					},
					Cookie: Cookie{
						Name:        "CAPUCHIN_SID",
						ExpiresDays: 365,
						Secure:      true,
					},
					TrustedProxies: []string{
						"127.0.0.1",
						"172.25.0.0/16",
					},
				},
				Sqlite: Sqlite{
					DBPath: "capuchin.db",
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetConfigEnv(t)
			tt.setEnvFunc()
			got, err := LoadConfiguration()
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

// resetConfigEnv сбрасывает переменные окружения конфига, которые могут быть унаследованы
// из внешней среды (например, директивы из .envrc), чтобы тесты LoadConfiguration не зависели
// от состояния окружения. Пустая строка для caarlos0/env v11 означает «переменная не задана»,
// поэтому применяются значения envDefault из config.go.
func resetConfigEnv(t *testing.T) {
	t.Helper()

	for _, name := range allEnvList {
		t.Setenv(name, "")
	}
}
