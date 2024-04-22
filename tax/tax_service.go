package tax

import (
	"github.com/rs/zerolog"

	"github.com/meteedev/assessment-tax/apperrs"
	"github.com/meteedev/assessment-tax/constant"
	taxrepo "github.com/meteedev/assessment-tax/repository"
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
	taxAmount, err := t.calculateTax(incomeDetail)
	if err != nil {
		t.logger.Error().Err(err).Msgf("Error occurred during tax calculation: %v", err)
		return nil, err
	}

	// Log the calculation details with formatting
	t.logger.Info().Msgf("Income: %.2f, Tax Amount: %.2f", incomeDetail.TotalIncome, taxAmount)

	taxResponse := TaxResponse{
		Tax: taxAmount,
	}

	return &taxResponse, nil
}

func (t *TaxService) calculateTax(incomeDetail *TaxRequest) (float64, error) {
	
	income := incomeDetail.TotalIncome
	allowances := incomeDetail.Allowances
	wht := incomeDetail.WHT
	
	var taxAmount float64

	// Log the income for calculation
	t.logger.Debug().Msgf("Calculating tax for income: %.2f", income)

	taxedIncome, err := deductPersonalAllowance(income)
	if err != nil {
		return 0, apperrs.NewInternalServerError(constant.MSG_BU_GERNERAL_ERROR)
	}
	t.logger.Debug().Msgf("Taxed income (%.2f) after deductPersonalAllowance", taxedIncome)

	taxedIncome, err = deductAllowance(taxedIncome,allowances)
	if err != nil {
		return 0, apperrs.NewInternalServerError(constant.MSG_BU_GERNERAL_ERROR)
	}
	t.logger.Debug().Msgf("Taxed income (%.2f) after deductAllowance", taxedIncome)

	// for _, bracket := range brackets {
	// 	// Log processing of each bracket
	// 	t.logger.Debug().Msgf("Processing bracket: %+v", bracket)

	// 	if taxedIncome <= bracket.UpperBound || bracket.UpperBound == 0 {
	// 		// Log when taxed income is within the bracket upper bound
	// 		t.logger.Debug().Msgf("Taxed income (%.2f) is within the bracket upper bound (%.2f)", taxedIncome, bracket.UpperBound)

	// 		taxableAmount := taxedIncome - adjustLowerBound(bracket.LowerBound-1)
	// 		if taxableAmount > 0 {
	// 			// Log taxable amount and tax amount after applying rate
	// 			t.logger.Debug().Msgf("Taxable amount: %.2f", taxableAmount)
	// 			taxAmount += taxableAmount * bracket.TaxRate
	// 			t.logger.Debug().Msgf("Tax amount after applying rate %.2f: %.2f", bracket.TaxRate, taxAmount)
	// 		}
	// 		break
	// 	} else {
	// 		// Log when taxed income exceeds the bracket upper bound
	// 		t.logger.Debug().Msgf("Taxed income (%.2f) exceeds the bracket upper bound (%.2f)", taxedIncome, bracket.UpperBound)

	// 		taxableAmount := bracket.UpperBound - adjustLowerBound(bracket.LowerBound-1)
	// 		t.logger.Debug().Msgf("Taxable amount: %.2f", taxableAmount)
	// 		taxAmount += taxableAmount * bracket.TaxRate
	// 		t.logger.Debug().Msgf("Tax amount after applying rate %.2f: %.2f", bracket.TaxRate, taxAmount)
	// 	}
	// }

	// for _, bracket := range brackets {

	// 	t.logger.Debug().Msgf("Processing bracket: %+v", bracket)

	// 	// Calculate taxable amount based on bracket boundaries
	// 	taxableAmount := min(taxedIncome, bracket.UpperBound) - bracket.LowerBound + 1

	// 	// If taxable amount is positive, calculate tax amount and append to tax levels
	// 	if taxableAmount > 0 {
	// 		taxAmount = taxableAmount * bracket.TaxRate
	// 	}

	// 	// If taxed income is within the bracket, break the loop
	// 	if taxedIncome <= bracket.UpperBound {
	// 		break
	// 	}
	// }

	taxAmount , err = t.calculateWithTaxTable(taxedIncome)
	if err != nil {
		return 0, apperrs.NewInternalServerError(constant.MSG_BU_GERNERAL_ERROR)
	}


	taxAmount = deductWht(taxAmount,wht)
	return taxAmount, nil
}

func (t *TaxService) calculateWithTaxTable(taxedIncome float64)(float64,error){
	
	brackets, err := getTaxTable()
	if err != nil {
		return 0, apperrs.NewInternalServerError(constant.MSG_BU_GERNERAL_ERROR)
	}

	var taxAmount float64

	for _, bracket := range brackets {
		t.logger.Debug().Msgf("Processing bracket: %+v", bracket)

		// Calculate taxable amount based on bracket boundaries
		taxableAmount := min(taxedIncome, bracket.UpperBound) - bracket.LowerBound + 1

		// If taxable amount is positive, calculate tax amount and append to tax levels
		if taxableAmount > 0 {
			taxAmount = taxableAmount * bracket.TaxRate
		}

		// If taxed income is within the bracket, break the loop
		if taxedIncome <= bracket.UpperBound {
			break
		}
	}
	return taxAmount,nil
}


// func CalculateTaxLevels(brackets []TaxBracket, taxedIncome float64) TaxLevels {
// 	var taxLevels TaxLevels

// 	// Iterate over each tax bracket
// 	for _, bracket := range brackets {
// 		var taxAmount float64

// 		// Calculate taxable amount based on bracket boundaries
// 		taxableAmount := min(taxedIncome, bracket.UpperBound) - bracket.LowerBound + 1

// 		// If taxable amount is positive, calculate tax amount and append to tax levels
// 		if taxableAmount > 0 {
// 			taxAmount = taxableAmount * bracket.TaxRate
// 		}

// 		// Append tax level to tax levels
// 		taxLevels = append(taxLevels, struct {
// 			Level string  `json:"level"`
// 			Tax   float64 `json:"tax"`
// 		}{
// 			Level: fmt.Sprintf("%.0f-%.0f", bracket.LowerBound, bracket.UpperBound),
// 			Tax:   taxAmount,
// 		})

// 		// If taxed income is within the bracket, break the loop
// 		if taxedIncome <= bracket.UpperBound {
// 			break
// 		}
// 	}

// 	return taxLevels
// }




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

func getTaxTable() ([]taxrepo.TaxBracket, error) {
	brackets := []taxrepo.TaxBracket{
		{LowerBound: 0, UpperBound: 150000, TaxRate: 0.00}, // Adjust the tax rate as needed
		{LowerBound: 150001, UpperBound: 500000, TaxRate: 0.10},
		{LowerBound: 500001, UpperBound: 1000000, TaxRate: 0.15},
		{LowerBound: 1000001, UpperBound: 2000000, TaxRate: 0.20},
		{LowerBound: 2000001, UpperBound: 0, TaxRate: 0.35},
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