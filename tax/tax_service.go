package tax

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

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
	income := incomeDetail.TotalIncome
	allowances := incomeDetail.Allowances

	// Calculate tax
	taxAmount, err := calculateTax(income, allowances)
	if err != nil {
		t.logger.Error().Err(err).Msgf("Error occurred during tax calculation: %v", err)
		return nil, err
	}

	// Log the calculation details with formatting
	t.logger.Info().Msgf("Income: %.2f, Tax Amount: %.2f", income, taxAmount)

	taxResponse := TaxResponse{
		Tax: taxAmount,
	}

	return &taxResponse, nil
}

func calculateTax(income float64, allowances []Allowance) (float64, error) {
	var taxAmount float64

	brackets, err := getTaxTable()
	if err != nil {
		return 0, apperrs.NewInternalServerError(constant.MSG_BU_GERNERAL_ERROR)
	}

	// Log the income for calculation
	log.Info().Msgf("Calculating tax for income: %.2f", income)

	taxedIncome, err := deductPersonalAllowance(income, allowances)
	if err != nil {
		return 0, apperrs.NewInternalServerError(constant.MSG_BU_GERNERAL_ERROR)
	}

	for _, bracket := range brackets {
		// Log processing of each bracket
		log.Info().Msgf("Processing bracket: %+v", bracket)

		if taxedIncome <= bracket.UpperBound || bracket.UpperBound == 0 {
			// Log when taxed income is within the bracket upper bound
			log.Info().Msgf("Taxed income (%.2f) is within the bracket upper bound (%.2f)", taxedIncome, bracket.UpperBound)

			taxableAmount := taxedIncome - adjustLowerBound(bracket.LowerBound-1)
			if taxableAmount > 0 {
				// Log taxable amount and tax amount after applying rate
				log.Info().Msgf("Taxable amount: %.2f", taxableAmount)
				taxAmount += taxableAmount * bracket.TaxRate
				log.Info().Msgf("Tax amount after applying rate %.2f: %.2f", bracket.TaxRate, taxAmount)
			}
			break
		} else {
			// Log when taxed income exceeds the bracket upper bound
			log.Info().Msgf("Taxed income (%.2f) exceeds the bracket upper bound (%.2f)", taxedIncome, bracket.UpperBound)

			taxableAmount := bracket.UpperBound - adjustLowerBound(bracket.LowerBound-1)
			log.Info().Msgf("Taxable amount: %.2f", taxableAmount)
			taxAmount += taxableAmount * bracket.TaxRate
			log.Info().Msgf("Tax amount after applying rate %.2f: %.2f", bracket.TaxRate, taxAmount)
		}
	}

	return taxAmount, nil
}

func deductPersonalAllowance(income float64, allowances []Allowance) (float64, error) {
	totalAllowance := 0.0
	for _, allowance := range allowances {
		totalAllowance += allowance.Amount
	}

	personalAllowance, err := getPersonalAllowance()
	if err != nil {
		return 0, apperrs.NewInternalServerError(constant.MSG_BU_GERNERAL_ERROR)
	}

	taxedIncome := income - totalAllowance - personalAllowance
	return taxedIncome, nil
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
