package capuchin

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/fstest"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStaticHandler(t *testing.T) {
	distFS := fstest.MapFS{
		"index.html": &fstest.MapFile{
			Data: []byte("<html><body>Capuchin</body></html>"),
		},
		"assets/app.js": &fstest.MapFile{
			Data: []byte("console.log('capuchin');"),
		},
	}
	excludedPaths := []string{"/api", "/metrics"}

	tests := []struct {
		name     string
		method   string
		target   string
		wantCode int
		wantBody string
	}{
		{
			name:     "API route is not intercepted by static handler",
			method:   http.MethodGet,
			target:   "/api/health",
			wantCode: http.StatusOK,
			wantBody: "OK",
		},
		{
			name:     "Unknown /api path returns 404 instead of SPA fallback",
			method:   http.MethodGet,
			target:   "/api/nonexistent",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Metrics path is excluded from static handler",
			method:   http.MethodGet,
			target:   "/metrics",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Metrics subpath is excluded from static handler",
			method:   http.MethodGet,
			target:   "/metrics/foo",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "POST to static path is not intercepted",
			method:   http.MethodPost,
			target:   "/",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Real static file is served from embed FS",
			method:   http.MethodGet,
			target:   "/assets/app.js",
			wantCode: http.StatusOK,
			wantBody: "console.log('capuchin');",
		},
		{
			name:     "Root serves index.html",
			method:   http.MethodGet,
			target:   "/",
			wantCode: http.StatusOK,
			wantBody: "<html><body>Capuchin</body></html>",
		},
		{
			name:     "Unknown non-API path serves index.html (SPA fallback)",
			method:   http.MethodGet,
			target:   "/reports/2026-01",
			wantCode: http.StatusOK,
			wantBody: "<html><body>Capuchin</body></html>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			app.Get("/api/health", func(c *fiber.Ctx) error {
				return c.SendString("OK")
			})
			app.Use(staticHandler(distFS, excludedPaths))
			app.Use(func(c *fiber.Ctx) error {
				c.Append("X-RouteNotFound", "yes")

				return c.SendStatus(http.StatusNotFound)
			})

			req := httptest.NewRequestWithContext(t.Context(), tt.method, tt.target, nil)
			resp, err := app.Test(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tt.wantCode, resp.StatusCode)
			if tt.wantBody == "" {
				return
			}

			body := new(bytes.Buffer)
			_, err = body.ReadFrom(resp.Body)
			require.NoError(t, err)
			assert.Equal(t, tt.wantBody, body.String())
		})
	}
}

func TestStaticHandlerServesCacheControlForStaticFiles(t *testing.T) {
	distFS := fstest.MapFS{
		"index.html": &fstest.MapFile{
			Data: []byte("<html><body>Capuchin</body></html>"),
		},
	}

	app := fiber.New()
	app.Use(staticHandler(distFS, []string{"/api", "/metrics"}))

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, staticCacheControl, resp.Header.Get("Cache-Control"))
}
