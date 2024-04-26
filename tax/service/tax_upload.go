package service

import (
	"encoding/csv"
	"mime/multipart"
	"strconv"
	"fmt"

	"github.com/meteedev/assessment-tax/apperrs"
	"github.com/meteedev/assessment-tax/constant"
)




func (t *TaxService) UploadCalculationTax(file *multipart.FileHeader)(*TaxUploadResponse,error) {
	
	taxRequests , err  := t.csvToTaxRequest(file)
	

	if err != nil{
		t.logger.Debug().Msg(err.Error())
		return nil, apperrs.NewInternalServerError(constant.MSG_BU_GENERAL_ERROR)
	}

	var taxUploads []TaxUpload
	for _ , taxRequest := range *taxRequests{
		taxResponse , err := t.CalculateTax(&taxRequest)

		if err != nil{
			t.logger.Debug().Msg(err.Error())
			return nil, apperrs.NewInternalServerError(constant.MSG_BU_GENERAL_ERROR)
		}

		taxUpload := getTaxUpload(&taxRequest,taxResponse)
		taxUploads = append(taxUploads, taxUpload)
	}

	uploadTaxResponse := TaxUploadResponse{
		Taxes:taxUploads,
	}
	
	return &uploadTaxResponse , nil 
}



func (t *TaxService) csvToTaxRequest(file *multipart.FileHeader)(*[]TaxRequest,error){
	
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


	var taxRequests []TaxRequest
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

		allowance := Allowance{
			AllowanceType:"donation",
			Amount:donation,
		}

		var allowances []Allowance
		allowances = append(allowances, allowance)

		taxRequest := TaxRequest{
			TotalIncome:totalIncome,
			WHT:wht,
			Allowances:allowances,
		}

		taxRequests = append(taxRequests, taxRequest)

	}
	
	return &taxRequests , nil 

}