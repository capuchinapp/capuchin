package middleware

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap/zapcore"

	"capuchin/internal/middleware/mocks"
)

func TestNewRequestLog(t *testing.T) {
	logger := mocks.NewLoggerMock(t)
	logger.EXPECT().
		Info(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		RunAndReturn(func(msg string, fields ...zapcore.Field) {
			assert.Equal(t, "Request", msg)

			for _, field := range fields {
				switch field.Key {
				case "ip":
					assert.Equal(t, "0.0.0.0", field.String)
				case "method":
					assert.Equal(t, http.MethodGet, field.String)
				case "path":
					assert.Equal(t, "/some/path", field.String)
				case "protocol":
					assert.Equal(t, "http", field.String)
				case "status":
					assert.Equal(t, int64(http.StatusOK), field.Integer)
				case "duration":
					assert.GreaterOrEqual(t, field.Integer, int64(0))
				case "user-agent":
					assert.Equal(t, "golang_test", field.String)
				}
			}
		})

	metricsInstance := mocks.NewMetricsMock(t)
	metricsInstance.EXPECT().HTTPRequestsCurrentInc()
	metricsInstance.EXPECT().HTTPRequestsCurrentDec()
	metricsInstance.EXPECT().HTTPRequestsTotalInc(http.MethodGet, "/some/path", "200")
	metricsInstance.EXPECT().HTTPRequestDurationSecondsObserve(http.MethodGet, "/some/path", mock.Anything)

	app := fiber.New()

	app.Use(NewRequestLog(logger, metricsInstance))

	app.Get("/some/path", func(c *fiber.Ctx) error {
		return c.SendStatus(http.StatusOK)
	})

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/some/path", nil)
	req.Header.Set("User-Agent", "golang_test")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	t.Logf("StatusCode: %d\n", resp.StatusCode)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	bodyBytes, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, "OK", string(bodyBytes))
}

func Test_cleanPath(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want string
	}{
		{
			name: "Already clean",
			s:    "/api/timelogs",
			want: "/api/timelogs",
		},
		{
			name: "Simple case without suffix",
			s:    "/api/timelogs/01KKDT647DBC7E18YWWB0EXRS4",
			want: "/api/timelogs/:id",
		},
		{
			name: "Simple case with suffix",
			s:    "/api/timelogs/01KKDRZE2Z8B2Q2YMQ0PMDRGM5/stop",
			want: "/api/timelogs/:id/stop",
		},
		{
			name: "Extended case without suffix",
			s:    "/api/timelogs/01KKDT647DBC7E18YWWB0EXRS4/01KKDD6XJ0E45XT90P9A08HWKE",
			want: "/api/timelogs/:id/:id",
		},
		{
			name: "Extended case with suffix",
			s:    "/api/clients/01KKDD6XJ0E45XT90P9A08HWKE/timelogs/01KKDRZE2Z8B2Q2YMQ0PMDRGM5/stop",
			want: "/api/clients/:id/timelogs/:id/stop",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cleanPath(tt.s)
			assert.Equal(t, tt.want, got)
		})
	}
}
