package apperrs

import (
	"net/http"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestNewNotFoundError(t *testing.T) {
	expectedMessage := "Not found"
	expectedCode := http.StatusNotFound

	err := NewNotFoundError(expectedMessage)
	echoErr, ok := err.(*echo.HTTPError)

	assert.True(t, ok, "error should be an echo.HTTPError")
	assert.Equal(t, expectedCode, echoErr.Code, "HTTP status code should match")
	assert.Equal(t, expectedMessage, echoErr.Message, "Message should match")
}

func TestNewInternalServerError(t *testing.T) {
	expectedMessage := "Internal server error"
	expectedCode := http.StatusInternalServerError

	err := NewInternalServerError(expectedMessage)
	echoErr, ok := err.(*echo.HTTPError)

	assert.True(t, ok, "error should be an echo.HTTPError")
	assert.Equal(t, expectedCode, echoErr.Code, "HTTP status code should match")
	assert.Equal(t, expectedMessage, echoErr.Message, "Message should match")
}

func TestNewUnprocessableEntity(t *testing.T) {
	expectedMessage := "Unprocessable entity"
	expectedCode := http.StatusUnprocessableEntity

	err := NewUnprocessableEntity(expectedMessage)
	echoErr, ok := err.(*echo.HTTPError)

	assert.True(t, ok, "error should be an echo.HTTPError")
	assert.Equal(t, expectedCode, echoErr.Code, "HTTP status code should match")
	assert.Equal(t, expectedMessage, echoErr.Message, "Message should match")
}

func TestNewBadRequestError(t *testing.T) {
	expectedMessage := "Bad request"
	expectedCode := http.StatusBadRequest

	err := NewBadRequestError(expectedMessage)
	echoErr, ok := err.(*echo.HTTPError)

	assert.True(t, ok, "error should be an echo.HTTPError")
	assert.Equal(t, expectedCode, echoErr.Code, "HTTP status code should match")
	assert.Equal(t, expectedMessage, echoErr.Message, "Message should match")
}
