package service

import (
	

	"github.com/rs/zerolog"

	"github.com/meteedev/assessment-tax/apperrs"
	"github.com/meteedev/assessment-tax/constant"
	"github.com/meteedev/assessment-tax/tax/repository"
)

type TaxService struct {
	logger *zerolog.Logger
	DeductRepo repository.TaxDeductConfigPort
}

func NewTaxService(logger *zerolog.Logger, deductRepo repository.TaxDeductConfigPort) TaxServicePort {
	return &TaxService{
		logger: logger,
		DeductRepo: deductRepo,
	}
}




func (t *TaxService) CalculationTax(incomeDetail *TaxRequest)(*TaxResponse,error) {
	
	err := ValidateTaxRequest(incomeDetail)
	
	if err != nil {
		return nil, apperrs.NewBadRequestError(err.Error())
	}
	
	// Calculate tax
	taxResponse, err := t.CalculateTax(incomeDetail)
	if err != nil {
		t.logger.Error().Err(err).Msgf("Error occurred during tax calculation: %v", err)
		return nil, err
	}

	// Log the calculation details with formatting
	t.logger.Info().Msgf("Income: %.2f, Tax Amount: %.2f", incomeDetail.TotalIncome, taxResponse.Tax)
	return taxResponse, nil
}

func (t *TaxService) CalculateTax(incomeDetail *TaxRequest) (*TaxResponse, error) {
	
	income := incomeDetail.TotalIncome
	allowances := incomeDetail.Allowances
	wht := incomeDetail.WHT
	

	// Log the income for calculation
	t.logger.Debug().Msgf("Calculating tax for income: %.2f", income)

	taxedIncome, err := t.deductPersonalAllowance(income)
	if err != nil {
		return nil, apperrs.NewInternalServerError(constant.MSG_BU_GERNERAL_ERROR)
	}
	t.logger.Debug().Msgf("Taxed income (%.2f) after deductPersonalAllowance", taxedIncome)

	taxedIncome, err = t.deductAllowance(taxedIncome,allowances)
	if err != nil {
		return nil, apperrs.NewInternalServerError(constant.MSG_BU_GERNERAL_ERROR)
	}
	t.logger.Debug().Msgf("Taxed income (%.2f) after deductAllowance", taxedIncome)

	
	taxResponse , err := t.calculateWithTaxTable(taxedIncome)
	if err != nil {
		return nil, apperrs.NewInternalServerError(constant.MSG_BU_GERNERAL_ERROR)
	}

	taxDiff := t.deductWht(taxResponse.Tax,wht)

	setUpTaxRefund(taxDiff,taxResponse)

	return taxResponse, nil
}


func setUpTaxRefund(taxDiff float64, taxResponse *TaxResponse){
	if taxDiff < 0{
		taxResponse.TaxRefund = taxDiff * (-1)
		taxResponse.Tax = 0
	}else{
		taxResponse.TaxRefund = 0 
		taxResponse.Tax = taxDiff
	}
}

func (t *TaxService) calculateWithTaxTable(taxedIncome float64,)(*TaxResponse,error){
	
	brackets, err := t.getTaxTable()
	if err != nil {
		return nil,apperrs.NewInternalServerError(constant.MSG_BU_GERNERAL_ERROR)
	}

	var totalTax float64

	for i := range brackets {
		// Calculate taxable amount based on bracket boundaries
		taxableAmount := min(taxedIncome, brackets[i].UpperBound) - brackets[i].LowerBound + 1

		// If taxable amount is positive, calculate tax amount and append to tax levels
		if taxableAmount > 0 {
			taxAmount := taxableAmount * brackets[i].TaxRate
			brackets[i].Tax = taxAmount
			totalTax  += taxAmount
			t.logger.Debug().Msgf("bracket: %+v taxableAmount: %.2f taxAmount:%.2f totalTax:%.2f", brackets[i],taxableAmount,taxAmount,totalTax)
		}

		// If taxed income is within the bracket, break the loop
		if taxedIncome <= brackets[i].UpperBound {
			break
		}
	}

	t.logger.Debug().Msgf("bracket: %+v ", brackets)

	taxResponse := TaxResponse{
		Tax: totalTax,
		TaxBracket: brackets,
	}

	return &taxResponse,nil
}

func (t *TaxService) UpdatePersonalAllowance(updateReq *UpdateDeductRequest)(*UpdateDeductResponse,error){
	
	amount := updateReq.Amount

	err := ValidatePersonaAllowance(amount)
	
	if err != nil {
		return nil, apperrs.NewBadRequestError(err.Error())
	}

	deductId := constant.DEDUCT_PERSONAL_ID	



	updateRow , err := t.DeductRepo.UpdateById(deductId,amount)
	


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

 	updDeductResponse :=UpdateDeductResponse{
		Amount: d.Amount,
	} 

	return &updDeductResponse, nil
	
}


func (t *TaxService) UpdateKreceiptAllowance(updateReq *UpdateDeductRequest)(*UpdateDeductResponse,error){
	
	amount := updateReq.Amount

	err := ValidateKreceiptAllowance(amount)
	
	if err != nil {
		return nil, apperrs.NewBadRequestError(err.Error())
	}

	deductId := constant.DEDUCT_K_RECEIPT_ID

	updateRow , err := t.DeductRepo.UpdateById(deductId,amount)
	
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

 	updDeductResponse :=UpdateDeductResponse{
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


func (t *TaxService) adjustMaximumKreceiptAllowanceDeduct(allowance float64) (float64 ,error){
	
	kreciptAllowanceConfig , err := t.getKreceiptAllowance()

	if err !=nil {
		return 0, err
	}
	
	if allowance > kreciptAllowanceConfig {
		return kreciptAllowanceConfig,nil
	}
	return allowance,nil
}



func (t *TaxService) deductPersonalAllowance(income float64) (float64, error) {
	personalAllowance, err := t.getPersonalAllowance()
	if err != nil {
		return 0, err
	}
	taxedIncome := income - personalAllowance
	return taxedIncome, nil
}



// func deductAllowance(income float64, allowances []Allowance) (float64, error) {
// 	totalAllowance := 0.0
// 	for _, allowance := range allowances {
// 		totalAllowance += adjustMaximumDonationAllowanceDeduct(allowance.Amount)
// 	}
// 	taxedIncome := income - totalAllowance 
// 	return taxedIncome, nil
// }


func (t *TaxService) deductAllowance(income float64, allowances []Allowance) (float64, error) {
	totalAllowance := 0.0
	for _, allowance := range allowances {
		
		if constant.DEDUCT_DONATION_ID == allowance.AllowanceType{
			totalAllowance += t.adjustMaximumDonationAllowanceDeduct(allowance.Amount)
		}

		if constant.DEDUCT_K_RECEIPT_ID == allowance.AllowanceType{
			
			krecieptAdjust , err := t.adjustMaximumKreceiptAllowanceDeduct(allowance.Amount)

			if err != nil {
				return 0,err
			}

			totalAllowance += krecieptAdjust
		}
		
	}
	taxedIncome := income - totalAllowance 
	return taxedIncome, nil
}


func (t *TaxService) deductWht(taxAmount float64, wht float64) (float64) {
	taxDiff := taxAmount - wht
	return taxDiff
}





func (t *TaxService) getTaxTable() ([]TaxBracket, error) {
	brackets := []TaxBracket{
		{Level:"0-150,000",LowerBound: 0, UpperBound: 150000, TaxRate: 0.00, Tax: 0.0 }, // Adjust the tax rate as needed
		{Level:"150,001-500,000",LowerBound: 150001, UpperBound: 500000, TaxRate: 0.10, Tax: 0.0},
		{Level:"500,001-1,000,000",LowerBound: 500001, UpperBound: 1000000, TaxRate: 0.15, Tax: 0.0},
		{Level:"1,000,001-2,000,000",LowerBound: 1000001, UpperBound: 2000000, TaxRate: 0.20, Tax: 0.0},
		{Level:"2,000,001 ขึ้นไป",LowerBound: 2000001, UpperBound: 0, TaxRate: 0.35,Tax: 0.0},
	}

	return brackets, nil
}


func (t *TaxService) adjustLowerBound(lower float64) float64 {
	if lower < 0 {
		return 0
	}
	return lower
}


func (t *TaxService) adjustMaximumDonationAllowanceDeduct(allowance float64) float64 {
	if allowance > constant.MAX_ALLOWANCE_DONATION_DEDUCT {
		return constant.MAX_ALLOWANCE_DONATION_DEDUCT
	}
	return allowance
}


