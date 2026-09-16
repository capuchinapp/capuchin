package restapi

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func testRunRequest(
	t *testing.T,
	app *fiber.App,
	req *http.Request,
	wantCode int,
	wantBody string,
) {
	t.Helper()

	resp, err := app.Test(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	t.Logf("StatusCode: %d\n", resp.StatusCode)
	assert.Equal(t, wantCode, resp.StatusCode)

	bodyBytes, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	t.Logf("Content-Type: %s\n", resp.Header.Get("Content-Type"))
	if resp.Header.Get("Content-Type") == "application/json" {
		var (
			expected any
			actual   any
		)

		err := json.Unmarshal([]byte(wantBody), &expected)
		assert.NoError(t, err)

		err = json.Unmarshal(bodyBytes, &actual)
		assert.NoError(t, err)

		assert.Equal(t, expected, actual)

		return
	}

	assert.Equal(t, wantBody, string(bodyBytes))
}
