package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/meteedev/assessment-tax/tax/service"
)

type TaxHandler struct {
	service service.TaxServicePort
}

func NewTaxHandler(service service.TaxServicePort) *TaxHandler {
	return &TaxHandler{service: service}
}

func (h *TaxHandler) TaxCalculationsHandler(c echo.Context) error {
	req := new(service.TaxRequest)
	
	if err := c.Bind(req); err != nil {
		return err
	}
	
	taxResponse , err :=h.service.CalculationTax(req)

	if err != nil{
		return err
	}

	return c.JSON(http.StatusOK, taxResponse)
}

func (h *TaxHandler) DeductionsPersonal(c echo.Context) error {	
	return c.JSON(http.StatusOK, "DeductionsPersonal")
}