package apperrs

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type AppError struct {
	Code    int    `json:"error_code"`
	Message string `json:"message"`
}

func (e AppError) Error() string {
	return e.Message
}

func NewNotFoundError(message string) error {
	return echo.NewHTTPError(http.StatusNotFound, message)
}

func NewInternalServerError(message string) error {
	return echo.NewHTTPError(http.StatusInternalServerError, message)
}

func NewUnprocessableEntity(message string) error {
	return echo.NewHTTPError(http.StatusUnprocessableEntity, message)
}

func NewBadRequestError(message string) error {
	return echo.NewHTTPError(http.StatusBadRequest, message)
}