package server

import (
	"net/http"
	"net/http/httptest"
	"testing"


	"github.com/labstack/echo/v4"
	"github.com/meteedev/assessment-tax/tax/handler"
	"github.com/stretchr/testify/assert"
)

func TestRegisterRoutes(t *testing.T) {
	// Create a new instance of echo.Echo
	e := echo.New()


	// Create a new TaxHandler instance (you may need to mock it if necessary)
	handler := &handler.TaxHandler{}

	// Register the routes
	registerRoutes(e, handler)

	// Perform assertions to ensure that the routes are registered correctly
	assert.NotNil(t, e)
	assert.NotNil(t, e.Group("/tax"))

}

func TestTaxRoutes(t *testing.T) {
	// Create a new instance of echo.Echo
	e := echo.New()

	// Create a new TaxHandler instance (you may need to mock it if necessary)
	handler := &handler.TaxHandler{}

	// Register the routes
	registerRoutes(e, handler)

	// Create a request to test the /tax routes
	req := httptest.NewRequest(http.MethodPost, "/tax/calculations", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	// Assert that the response status code is OK
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestAdminRoutes(t *testing.T) {
	// Create a new instance of echo.Echo
	e := echo.New()

	// Create a new TaxHandler instance (you may need to mock it if necessary)
	handler := &handler.TaxHandler{}

	// Register the routes
	registerRoutes(e, handler)

	// Create a request to test the /admin routes
	req := httptest.NewRequest(http.MethodPost, "/admin/deductions/personal", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	// Assert that the response status code is OK
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

