package service

import (
	"encoding/csv"
	"io"
	"strconv"

	"github.com/meteedev/assessment-tax/apperrs"
	"github.com/meteedev/assessment-tax/constant"
)

type CSVParserImpl struct{}

func (c *CSVParserImpl) ParseCSVToTaxRequest(file io.Reader) (*[]TaxRequest, error) {
	reader := csv.NewReader(file)
	_, err := reader.Read() // Skip the header row
	if err != nil {
		return nil, apperrs.NewBadRequestError(err.Error())
	}

	var taxRequests []TaxRequest
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, apperrs.NewBadRequestError(err.Error())
		}

		if len(record)!= constant.CSV_UPLOAD_COLUMN{
			return nil, apperrs.NewBadRequestError(constant.MSG_UPLOAD_CSV_WRONG_FORMAT)
		}

		taxRequest, err := parseTaxRequestRecord(record)
		if err != nil {
			return nil, apperrs.NewBadRequestError(err.Error())
		}

		taxRequests = append(taxRequests, *taxRequest)
	}

	return &taxRequests, nil
}

func parseTaxRequestRecord(record []string) (*TaxRequest, error) {
	
	totalIncome, err := strconv.ParseFloat(record[0], 64)
	if err != nil {
		return nil, err
	}

	wht, err := strconv.ParseFloat(record[1], 64)
	if err != nil {
		return nil, err
	}

	donation, err := strconv.ParseFloat(record[2], 64)
	if err != nil {
		return nil, err
	}

	allowance := Allowance{
		AllowanceType: "donation",
		Amount:        donation,
	}

	taxRequest := TaxRequest{
		TotalIncome: totalIncome,
		WHT:         wht,
		Allowances:  []Allowance{allowance},
	}

	return &taxRequest, nil
}


func (t *TaxService) UploadCalculationTax(file io.Reader) (*TaxUploadResponse, error) {
	//t.logger.Debug().Msg("Uploading calculation tax from reader")

	taxRequests, err := t.csvParser.ParseCSVToTaxRequest(file)
	if err != nil {
		t.logger.Debug().Msg(err.Error())
		return nil, err
	}

	var taxUploads []TaxUpload
	for _, taxRequest := range *taxRequests {
		if err := ValidateTaxRequest(&taxRequest); err != nil {
			return nil, apperrs.NewBadRequestError(err.Error())
		}

		taxResponse, err := t.CalculateTax(&taxRequest)
		if err != nil {
			t.logger.Debug().Msg(err.Error())
			return nil, apperrs.NewInternalServerError(constant.MSG_BU_GENERAL_ERROR)
		}

		taxUpload := getTaxUpload(&taxRequest, taxResponse)
		taxUploads = append(taxUploads, taxUpload)
	}

	return &TaxUploadResponse{Taxes: taxUploads}, nil
}



