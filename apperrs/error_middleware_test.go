package apperrs_test

import (
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/labstack/echo/v4"
    "github.com/stretchr/testify/assert"
    "github.com/meteedev/assessment-tax/apperrs"
)

func TestCustomErrorMiddleware(t *testing.T) {
    // Create a new Echo instance
    e := echo.New()

    // Define a test handler that always returns an error
    testHandler := func(c echo.Context) error {
        return echo.NewHTTPError(http.StatusBadRequest, "Bad Request")
    }

    // Register the middleware
    e.Use(apperrs.CustomErrorMiddleware)

    // Register the test handler
    e.GET("/", testHandler)

    // Create a request with the test handler
    req := httptest.NewRequest(http.MethodGet, "/", nil)
    rec := httptest.NewRecorder()
    _ = e.NewContext(req, rec)

    // Perform the request
    e.ServeHTTP(rec, req)

    // Check the response
    assert.Equal(t, http.StatusBadRequest, rec.Code)
    assert.JSONEq(t, `{"code":400,"message":"Bad Request"}`, rec.Body.String())
}
