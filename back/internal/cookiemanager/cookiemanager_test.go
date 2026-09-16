package cookiemanager

import (
	"regexp"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

func TestCookieManager_Create(t *testing.T) {
	tests := []struct {
		name    string
		srvFunc func(t *testing.T) *CookieManager
		v       string
		want    string
	}{
		{
			name: "Should create cookie",
			srvFunc: func(t *testing.T) *CookieManager {
				t.Helper()

				return &CookieManager{
					name:        "CAPUCHIN_CID",
					domain:      "localhost",
					expiresDays: 1,
					secure:      true,
					nowFunc:     func() time.Time { return time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC) },
				}
			},
			v:    "01JHHXTVAMTYVE7Q3H6VRP7206",
			want: "CAPUCHIN_CID=01JHHXTVAMTYVE7Q3H6VRP7206; expires=Thu, 19 Dec 2024 14:15:00 GMT; domain=localhost; path=/; HttpOnly; secure; SameSite=Strict",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()

			ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
			defer app.ReleaseCtx(ctx)

			tt.srvFunc(t).Create(ctx, tt.v)
			assert.Equal(t, tt.want, string(ctx.Response().Header.Peek("Set-Cookie")))
		})
	}
}

func TestCookieManager_Get(t *testing.T) {
	tests := []struct {
		name    string
		srvFunc func(t *testing.T) *CookieManager
		cookie  string
		want    string
	}{
		{
			name: "Should get cookie",
			srvFunc: func(t *testing.T) *CookieManager {
				t.Helper()

				return &CookieManager{
					re: regexp.MustCompile(`(?m)CAPUCHIN_CID=(\w{26})`),
				}
			},
			cookie: "01JHHXTVAMTYVE7Q3H6VRP7206",
			want:   "01JHHXTVAMTYVE7Q3H6VRP7206",
		},
		{
			name: "Should not get cookie",
			srvFunc: func(t *testing.T) *CookieManager {
				t.Helper()

				return &CookieManager{
					re: regexp.MustCompile(`(?m)CAPUCHIN_CID=(\w{26})`),
				}
			},
			cookie: "",
			want:   "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()

			req := &fasthttp.RequestCtx{}
			req.Request.Header.SetCookie("CAPUCHIN_CID", tt.cookie)

			ctx := app.AcquireCtx(req)
			defer app.ReleaseCtx(ctx)

			srv := tt.srvFunc(t)

			got := srv.Get(ctx)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCookieManager_GetExpiresDays(t *testing.T) {
	tests := []struct {
		name    string
		srvFunc func(t *testing.T) *CookieManager
		want    int
	}{
		{
			name: "Should get expires days",
			srvFunc: func(t *testing.T) *CookieManager {
				t.Helper()

				return &CookieManager{
					expiresDays: 30,
				}
			},
			want: 30,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := tt.srvFunc(t)

			got := srv.GetExpiresDays()
			assert.Equal(t, tt.want, got)
		})
	}
}
