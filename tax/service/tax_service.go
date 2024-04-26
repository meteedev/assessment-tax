package service

import (
	"github.com/rs/zerolog"
	"github.com/meteedev/assessment-tax/apperrs"
	"github.com/meteedev/assessment-tax/constant"
	"github.com/meteedev/assessment-tax/tax/repository"
)

type TaxService struct {
	logger     *zerolog.Logger
	DeductRepo repository.TaxDeductConfigPort
}

func NewTaxService(logger *zerolog.Logger, deductRepo repository.TaxDeductConfigPort) TaxServicePort {
	return &TaxService{
		logger:     logger,
		DeductRepo: deductRepo,
	}
}

func (t *TaxService) CalculationTax(incomeDetail *TaxRequest) (*TaxResponse, error) {
	err := ValidateTaxRequest(incomeDetail)
	if err != nil {
		return nil, apperrs.NewBadRequestError(err.Error())
	}

	taxResponse, err := t.CalculateTax(incomeDetail)
	if err != nil {
		t.logger.Error().Err(err).Msgf("Error occurred during tax calculation: %v", err)
		return nil, err
	}

	t.logger.Info().Msgf("Income: %.2f, Tax Amount: %.2f", incomeDetail.TotalIncome, taxResponse.Tax)
	return taxResponse, nil
}

func (t *TaxService) CalculateTax(incomeDetail *TaxRequest) (*TaxResponse, error) {
	income := incomeDetail.TotalIncome
	allowances := incomeDetail.Allowances
	wht := incomeDetail.WHT

	t.logger.Debug().Msgf("Calculating tax for income: %.2f", income)

	taxedIncome, err := t.deductPersonalAllowance(income)
	if err != nil {
		return nil, apperrs.NewInternalServerError(constant.MSG_BU_GERNERAL_ERROR)
	}
	t.logger.Debug().Msgf("Taxed income (%.2f) after deductPersonalAllowance", taxedIncome)

	taxedIncome, err = t.deductAllowance(taxedIncome, allowances)
	if err != nil {
		return nil, apperrs.NewInternalServerError(constant.MSG_BU_GERNERAL_ERROR)
	}
	t.logger.Debug().Msgf("Taxed income (%.2f) after deductAllowance", taxedIncome)

	// calculate tax table
	taxStep , totalTax := t.calculateTaxTable(taxedIncome)
	
	taxDiff := t.deductWht(totalTax, wht)
	
	taxResponse := getTaxResponse(taxDiff,taxStep)

	return &taxResponse, nil
}

func getTaxResponse(taxDiff float64,taxStep []TaxStep) TaxResponse {
	
	taxResponse := TaxResponse{
		TaxStep: taxStep,
	}

	if taxDiff < 0 {
		taxResponse.TaxRefund = taxDiff * (-1)
		taxResponse.Tax = 0
	} else {
		taxResponse.TaxRefund = 0
		taxResponse.Tax = taxDiff
	}

	return taxResponse
}


func (t *TaxService) UpdatePersonalAllowance(updateReq *UpdateDeductRequest) (*UpdateDeductResponse, error) {
	amount := updateReq.Amount

	err := ValidatePersonaAllowance(amount)
	if err != nil {
		return nil, apperrs.NewBadRequestError(err.Error())
	}

	deductId := constant.DEDUCT_PERSONAL_ID
	updateRow, err := t.DeductRepo.UpdateById(deductId, amount)
	if err != nil {
		t.logger.Error().Msg(err.Error())
		return nil, apperrs.NewInternalServerError(constant.MSG_BU_DEDUCT_UPD_PERSONAL_FAILED)
	}

	if updateRow == 0 {
		return nil, apperrs.NewUnprocessableEntity(constant.MSG_BU_DEDUCT_UPD_PERSONAL_FAILED)
	}

	d, err := t.DeductRepo.FindById(deductId)
	if err != nil {
		t.logger.Error().Msg(err.Error())
		return nil, apperrs.NewInternalServerError(constant.MSG_BU_DEDUCT_UPD_PERSONAL_FAILED)
	}

	updDeductResponse := UpdateDeductResponse{
		Amount: d.Amount,
	}

	return &updDeductResponse, nil
}

func (t *TaxService) UpdateKreceiptAllowance(updateReq *UpdateDeductRequest) (*UpdateDeductResponse, error) {
	amount := updateReq.Amount

	err := ValidateKreceiptAllowance(amount)
	if err != nil {
		return nil, apperrs.NewBadRequestError(err.Error())
	}

	deductId := constant.DEDUCT_K_RECEIPT_ID
	updateRow, err := t.DeductRepo.UpdateById(deductId, amount)
	if err != nil {
		t.logger.Error().Msg(err.Error())
		return nil, apperrs.NewInternalServerError(constant.MSG_BU_DEDUCT_UPD_PERSONAL_FAILED)
	}

	if updateRow == 0 {
		return nil, apperrs.NewUnprocessableEntity(constant.MSG_BU_DEDUCT_UPD_PERSONAL_FAILED)
	}

	d, err := t.DeductRepo.FindById(deductId)
	if err != nil {
		t.logger.Error().Msg(err.Error())
		return nil, apperrs.NewInternalServerError(constant.MSG_BU_DEDUCT_UPD_PERSONAL_FAILED)
	}

	updDeductResponse := UpdateDeductResponse{
		Amount: d.Amount,
	}

	return &updDeductResponse, nil
}

func (t *TaxService) getPersonalAllowance() (float64, error) {
	personAllowance, err := t.DeductRepo.FindById(constant.DEDUCT_PERSONAL_ID)
	if err != nil {
		return 0, apperrs.NewInternalServerError(constant.MSG_BU_DEDUCT_PERSONAL_CONFIG_NOT_FOUND)
	}
	return personAllowance.Amount, nil
}

func (t *TaxService) getKreceiptAllowance() (float64, error) {
	kreceiptAllowance, err := t.DeductRepo.FindById(constant.DEDUCT_K_RECEIPT_ID)
	if err != nil {
		return 0, apperrs.NewInternalServerError(constant.MSG_BU_DEDUCT_K_RECEIPT_CONFIG_NOT_FOUND)
	}
	return kreceiptAllowance.Amount, nil
}

func (t *TaxService) adjustMaximumKreceiptAllowanceDeduct(allowance float64) (float64, error) {
	kreciptAllowanceConfig, err := t.getKreceiptAllowance()
	if err != nil {
		return 0, err
	}

	if allowance > kreciptAllowanceConfig {
		return kreciptAllowanceConfig, nil
	}
	return allowance, nil
}

func (t *TaxService) deductPersonalAllowance(income float64) (float64, error) {
	personalAllowance, err := t.getPersonalAllowance()
	if err != nil {
		return 0, err
	}
	taxedIncome := income - personalAllowance
	return taxedIncome, nil
}

func (t *TaxService) deductAllowance(income float64, allowances []Allowance) (float64, error) {
	totalAllowance := 0.0
	for _, allowance := range allowances {
		switch allowance.AllowanceType {
		case constant.DEDUCT_DONATION_ID:
			totalAllowance += t.adjustMaximumDonationAllowanceDeduct(allowance.Amount)
		case constant.DEDUCT_K_RECEIPT_ID:
			krecieptAdjust, err := t.adjustMaximumKreceiptAllowanceDeduct(allowance.Amount)
			if err != nil {
				return 0, err
			}
			totalAllowance += krecieptAdjust
		}
	}
	taxedIncome := income - totalAllowance
	return taxedIncome, nil
}

func (t *TaxService) deductWht(taxAmount float64, wht float64) float64 {
	taxDiff := taxAmount - wht
	return taxDiff
}



func (t *TaxService) adjustMaximumDonationAllowanceDeduct(allowance float64) float64 {
	if allowance > constant.MAX_ALLOWANCE_DONATION_DEDUCT {
		return constant.MAX_ALLOWANCE_DONATION_DEDUCT
	}
	return allowance
}

func (t *TaxService) UploadCalculationTax(taxRequests *[]TaxRequest)(*TaxUploadResponse,error){
	
	var taxUploads []TaxUpload
	for _ , taxRequest := range *taxRequests{
		taxResponse , err := t.CalculateTax(&taxRequest)

		if err != nil{
			return nil, apperrs.NewInternalServerError(constant.MSG_BU_GERNERAL_ERROR)
		}

		taxUpload := getTaxUpload(&taxRequest,taxResponse)
		taxUploads = append(taxUploads, taxUpload)
	}

	uploadTaxResponse := TaxUploadResponse{
		Taxes:taxUploads,
	}
	
	return &uploadTaxResponse , nil 
}


func getTaxUpload(taxRequest *TaxRequest, taxResponse *TaxResponse) TaxUpload {
	
	taxUpload := TaxUpload{
		TotalIncome: taxRequest.TotalIncome,
		Tax: taxResponse.Tax,
		TaxRefund: taxResponse.TaxRefund,
	}

	return taxUpload
}
