package capuchin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
				t.Setenv("LOG_LEVEL", "debug")

				t.Setenv("METRICS_PORT", "8080")

				t.Setenv("HTTP_PORT", "3001")

				t.Setenv("HTTP_CORS_ENABLED", "false")
				t.Setenv("HTTP_CORS_ALLOW_ORIGINS", "http://localhost:5173")
				t.Setenv("HTTP_CORS_ALLOW_METHODS", "GET,POST")
				t.Setenv("HTTP_CORS_ALLOW_CREDENTIALS", "false")

				t.Setenv("HTTP_TRUSTED_PROXIES", "127.0.0.1,172.25.0.0/16")

				t.Setenv("COOKIE_NAME", "CAPUCHIN_SID_TEST")
				t.Setenv("COOKIE_EXPIRES_DAYS", "100")
				t.Setenv("COOKIE_SECURE", "false")
				t.Setenv("COOKIE_DOMAIN", "localhost")

				t.Setenv("SQLITE_DB_PATH", "capuchin.db")

				t.Setenv("MAILER_HOST", "localhost")
				t.Setenv("MAILER_PORT", "25")
				t.Setenv("MAILER_USERNAME", "username")
				t.Setenv("MAILER_PASSWORD", "password")
				t.Setenv("MAILER_FROM", "from@localhost.tld")
			},
			want: &Configuration{
				Log: Log{
					Level: -1,
				},
				Metrics: Metrics{
					Port: 8080,
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
						"172.25.0.0/16",
					},
				},
				Sqlite: Sqlite{
					DBPath: "capuchin.db",
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setEnvFunc()
			got, err := LoadConfiguration()
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
