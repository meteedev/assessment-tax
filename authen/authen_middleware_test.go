package authen

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
	// Mock environment variables
	os.Setenv("ADMIN_USERNAME", "user")
	os.Setenv("ADMIN_PASSWORD", "pass")

	// Create a new Echo instance
	e := echo.New()

	// Create a new request with basic authentication
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.SetBasicAuth("user", "pass")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Call the AuthMiddleware function
	result, err := AuthMiddleware("user", "pass", c)

	// Check the result
	assert.NoError(t, err)
	assert.True(t, result)

	// Test with incorrect username/password
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req.SetBasicAuth("user", "pass")
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)

	result, err = AuthMiddleware("u", "p", c)
	assert.NoError(t, err)
	assert.False(t, result)
}
