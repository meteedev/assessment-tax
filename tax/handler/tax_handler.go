package handler

import (
	"encoding/csv"
	"fmt"
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

func (h *TaxHandler) TaxCalculation(c echo.Context) error {
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
	req := new(service.UpdateDeductRequest)
	
	if err := c.Bind(req); err != nil {
		return err
	}

	updateResponse , err :=h.service.UpdatePersonalAllowance(req)

	if err != nil{
		return err
	}


	return c.JSON(http.StatusOK, updateResponse)
}

func (h *TaxHandler) DeductionsKreceipt(c echo.Context) error {	
	req := new(service.UpdateDeductRequest)
	
	if err := c.Bind(req); err != nil {
		return err
	}

	updateResponse , err :=h.service.UpdateKreceiptAllowance(req)

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

	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	reader := csv.NewReader(src)

	for {
		record, err := reader.Read()
		if err != nil {
			break // End of file
		}
		// Process each CSV record as needed
		fmt.Println(record)
	}
	
	
	return c.String(http.StatusOK, "CSV data received successfully")

}