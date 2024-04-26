package service

import (
	"encoding/csv"
	"fmt"
	"io"
	"mime/multipart"
	"strconv"

	"github.com/meteedev/assessment-tax/apperrs"
	"github.com/meteedev/assessment-tax/constant"
)




func (t *TaxService) UploadCalculationTax(file *multipart.FileHeader)(*TaxUploadResponse,error) {
	
	t.logger.Debug().Msg(file.Filename)

	taxRequests , err  := t.csvToTaxRequest(file)
	

	if err != nil{
		t.logger.Debug().Msg(err.Error())
		return nil, err
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
	
	t.logger.Debug().Msg("csvToTaxRequest")

	src, err := file.Open()

	t.logger.Debug().Msg("after open ")
	if err != nil {
		t.logger.Debug().Msg(err.Error())
		fmt.Println(err.Error())
		return nil, apperrs.NewBadRequestError(err.Error())
	}
	defer src.Close()

	reader := csv.NewReader(src)

	//Skip the first row
	_, err = reader.Read()
    if err != nil {
		return nil, apperrs.NewBadRequestError(err.Error())
    }

	t.logger.Debug().Msg("after skip row ")

	var taxRequests []TaxRequest
	var totalIncome float64
	var wht float64 
	var donation float64

	for {
		record, err := reader.Read()
		if err == io.EOF { // Check for end of file
			break
		} else if err != nil {
			fmt.Println("Error:", err)
			return nil, apperrs.NewBadRequestError(err.Error())
		}

		// Process each CSV record as needed
		fmt.Println(record)

		err = ValidateUploadTaxCsvRecord(record)
		
		if err != nil {
			return nil, apperrs.NewBadRequestError(err.Error())
		}

		totalIncome,err = strconv.ParseFloat(record[0], 64)
		if err != nil {
			return nil, apperrs.NewBadRequestError(err.Error())
		}

		wht,err = strconv.ParseFloat(record[1], 64)
		if err != nil {
			return nil, apperrs.NewBadRequestError(err.Error())
		}
		
		donation,err = strconv.ParseFloat(record[2], 64)
		if err != nil {
			return nil, apperrs.NewBadRequestError(err.Error())
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