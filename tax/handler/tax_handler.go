package handler

import (
	"net/http"
	"encoding/json"

	"github.com/labstack/echo/v4"
	"github.com/meteedev/assessment-tax/tax/service"

)

type TaxHandler struct {
	service service.TaxServicePort
}

func NewTaxHandler(service service.TaxServicePort) *TaxHandler {
	return &TaxHandler{service: service}
}

func (h *TaxHandler) TaxCalculation(c echo.Context) error {
	
	body , err := h.validateSchema(c, TAX_REQUEST_SCHEMA) 
	if err != nil {
		return err
	}

	// Define a variable of type MyJSONData
	var taxRequest service.TaxRequest

	// Unmarshal the JSON data into the struct
	if err := json.Unmarshal(body, &taxRequest); err != nil {
		return err
	}

	taxResponse , err :=h.service.CalculationTax(&taxRequest)
	if err != nil{
		return err
	}

	return c.JSON(http.StatusOK, taxResponse)
}

func (h *TaxHandler) DeductionsPersonal(c echo.Context) error {	

	body , err := h.validateSchema(c, UPDATE_DEDUCT_REQUEST) 
	if err != nil {
		return err
	}

	var deductRequest service.UpdateDeductRequest
	if err := json.Unmarshal(body, &deductRequest); err != nil {
		return err
	}

	updateResponse , err :=h.service.UpdatePersonalAllowance(&deductRequest)
	if err != nil{
		return err
	}

	return c.JSON(http.StatusOK, updateResponse)
}

func (h *TaxHandler) DeductionsKreceipt(c echo.Context) error {	
	body , err := h.validateSchema(c, UPDATE_DEDUCT_REQUEST) 
	if err != nil {
		return err
	}

	var deductRequest service.UpdateDeductRequest
	if err := json.Unmarshal(body, &deductRequest); err != nil {
		return err
	}

	updateResponse , err :=h.service.UpdateKreceiptAllowance(&deductRequest)
	if err != nil{
		return err
	}
	
	return c.JSON(http.StatusOK, updateResponse)
}

func (h *TaxHandler) TaxUploadCalculation(c echo.Context) error {	
	
	file, err := c.FormFile("taxFile")
	if err != nil {
		return err
	}

	uploadTaxResponse , err := h.service.UploadCalculationTax(file)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, uploadTaxResponse)

}

