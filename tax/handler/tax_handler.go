package handler

import (
	"encoding/csv"
	"fmt"
	"mime/multipart"
	"net/http"
	"strconv"

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
		fmt.Println(err.Error())
		return err
	}

	uploadTaxResponse , err := h.csvToTaxRequest(file)

	if err != nil {
		return err
	}

	fmt.Println(uploadTaxResponse)

	return c.JSON(http.StatusOK, uploadTaxResponse)

}


func (h *TaxHandler) csvToTaxRequest(file *multipart.FileHeader)([]service.TaxRequest,error){
	
	src, err := file.Open()
	if err != nil {
		fmt.Println(err.Error())
		return nil,err
	}
	defer src.Close()

	reader := csv.NewReader(src)

	//Skip the first row
	_, err = reader.Read()
    if err != nil {
		fmt.Println(err.Error())
    }


	var taxRequests []service.TaxRequest
	var totalIncome float64
	var wht float64 
	var donation float64

	for {
		record, err := reader.Read()
		if err != nil {
			break // End of file
		}
		// Process each CSV record as needed
		fmt.Println(record)

		totalIncome,err = strconv.ParseFloat(record[0], 64)
		if err != nil {
			break // End of file
		}

		wht,err = strconv.ParseFloat(record[1], 64)
		if err != nil {
			break // End of file
		}
		
		donation,err = strconv.ParseFloat(record[2], 64)
		if err != nil {
			break // End of file
		}

		allowance := service.Allowance{
			AllowanceType:"donation",
			Amount:donation,
		}

		var allowances []service.Allowance
		allowances = append(allowances, allowance)

		taxRequest := service.TaxRequest{
			TotalIncome:totalIncome,
			WHT:wht,
			Allowances:allowances,
		}

		taxRequests = append(taxRequests, taxRequest)

	}
	
	return taxRequests , nil 
}