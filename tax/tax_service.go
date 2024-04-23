package tax

import (
	"github.com/rs/zerolog"

	"github.com/meteedev/assessment-tax/apperrs"
	"github.com/meteedev/assessment-tax/constant"
)

type TaxService struct {
	logger *zerolog.Logger
}

func NewTaxService(logger *zerolog.Logger) *TaxService {
	return &TaxService{
		logger: logger,
	}
}

func (t *TaxService) CalculationTax(incomeDetail *TaxRequest) (*TaxResponse, error) {
	
	err := ValidateTaxRequest(incomeDetail)
	
	if err != nil {
		return nil, apperrs.NewBadRequestError(err.Error())
	}
	
	// Calculate tax
	taxResponse, err := t.calculateTax(incomeDetail)
	if err != nil {
		t.logger.Error().Err(err).Msgf("Error occurred during tax calculation: %v", err)
		return nil, err
	}

	// Log the calculation details with formatting
	t.logger.Info().Msgf("Income: %.2f, Tax Amount: %.2f", incomeDetail.TotalIncome, taxResponse.Tax)
	return taxResponse, nil
}

func (t *TaxService) calculateTax(incomeDetail *TaxRequest) (*TaxResponse, error) {
	
	income := incomeDetail.TotalIncome
	allowances := incomeDetail.Allowances
	wht := incomeDetail.WHT
	

	// Log the income for calculation
	t.logger.Debug().Msgf("Calculating tax for income: %.2f", income)

	taxedIncome, err := deductPersonalAllowance(income)
	if err != nil {
		return nil, apperrs.NewInternalServerError(constant.MSG_BU_GERNERAL_ERROR)
	}
	t.logger.Debug().Msgf("Taxed income (%.2f) after deductPersonalAllowance", taxedIncome)

	taxedIncome, err = deductAllowance(taxedIncome,allowances)
	if err != nil {
		return nil, apperrs.NewInternalServerError(constant.MSG_BU_GERNERAL_ERROR)
	}
	t.logger.Debug().Msgf("Taxed income (%.2f) after deductAllowance", taxedIncome)

	
	taxResponse , err := t.calculateWithTaxTable(taxedIncome)
	if err != nil {
		return nil, apperrs.NewInternalServerError(constant.MSG_BU_GERNERAL_ERROR)
	}

	taxResponse.Tax = deductWht(taxResponse.Tax,wht)

	return taxResponse, nil
}

// func (t *TaxService) calculateWithTaxTable(taxedIncome float64)(float64,error){
	
// 	brackets, err := getTaxTable()
// 	if err != nil {
// 		return 0, apperrs.NewInternalServerError(constant.MSG_BU_GERNERAL_ERROR)
// 	}

// 	var taxAmount float64

// 	for _, bracket := range brackets {
// 		t.logger.Debug().Msgf("Processing bracket: %+v", bracket)

// 		// Calculate taxable amount based on bracket boundaries
// 		taxableAmount := min(taxedIncome, bracket.UpperBound) - bracket.LowerBound + 1

// 		// If taxable amount is positive, calculate tax amount and append to tax levels
// 		if taxableAmount > 0 {
// 			taxAmount = taxableAmount * bracket.TaxRate
// 		}

// 		// If taxed income is within the bracket, break the loop
// 		if taxedIncome <= bracket.UpperBound {
// 			break
// 		}
// 	}
// 	return taxAmount,nil
// }


func (t *TaxService) calculateWithTaxTable(taxedIncome float64,)(*TaxResponse,error){
	
	brackets, err := getTaxTable()
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



func deductPersonalAllowance(income float64) (float64, error) {
	personalAllowance, err := getPersonalAllowance()
	if err != nil {
		return 0, apperrs.NewInternalServerError(constant.MSG_BU_GERNERAL_ERROR)
	}
	taxedIncome := income - personalAllowance
	return taxedIncome, nil
}


func deductAllowance(income float64, allowances []Allowance) (float64, error) {
	totalAllowance := 0.0
	for _, allowance := range allowances {
		totalAllowance += adjustMaximumAllowanceDeduct(allowance.Amount)
	}
	taxedIncome := income - totalAllowance 
	return taxedIncome, nil
}



func deductWht(taxAmount float64, wht float64) (float64) {
	return adjustLowerBound(taxAmount - wht)
}



func getPersonalAllowance() (float64, error) {
	return 60000.0, nil
}

func getTaxTable() ([]TaxBracket, error) {
	brackets := []TaxBracket{
		{Level:"0-150,000",LowerBound: 0, UpperBound: 150000, TaxRate: 0.00, Tax: 0.0 }, // Adjust the tax rate as needed
		{Level:"150,001-500,000",LowerBound: 150001, UpperBound: 500000, TaxRate: 0.10, Tax: 0.0},
		{Level:"500,001-1,000,000",LowerBound: 500001, UpperBound: 1000000, TaxRate: 0.15, Tax: 0.0},
		{Level:"1,000,001-2,000,000",LowerBound: 1000001, UpperBound: 2000000, TaxRate: 0.20, Tax: 0.0},
		{Level:"2,000,001 ขึ้นไป",LowerBound: 2000001, UpperBound: 0, TaxRate: 0.35,Tax: 0.0},
	}

	return brackets, nil
}


func adjustLowerBound(lower float64) float64 {
	if lower < 0 {
		return 0
	}
	return lower
}


func adjustMaximumAllowanceDeduct(allowance float64) float64 {
	if allowance > constant.MAX_ALLOWANCE_DEDUCT {
		return constant.MAX_ALLOWANCE_DEDUCT
	}
	return allowance
}