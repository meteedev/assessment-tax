package tax

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) TaxCalculationsHandler(c echo.Context) error {
	req := new(TaxRequest)
	
	if err := c.Bind(req); err != nil {
		return err
	}
	
	taxResponse , err :=h.service.CalculationTax(req)

	if err != nil{
		return err
	}

	return c.JSON(http.StatusOK, taxResponse)
}

func (h *Handler) DeductionsPersonal(c echo.Context) error {	
	return c.JSON(http.StatusOK, "DeductionsPersonal")
}