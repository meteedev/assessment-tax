package tax

import (
	"testing"

	taxrepo "github.com/meteedev/assessment-tax/repository"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestTaxService_CalculationTax(t *testing.T) {
	// Mock TaxRequest
	mockIncomeDetail := &TaxRequest{
		TotalIncome: 500000,
		WHT:         0,
		Allowances:  []Allowance{{AllowanceType: "", Amount: 0}, {AllowanceType: "", Amount: 0}},
	}

	logger := zerolog.Nop()

	// Create an instance of TaxService
	taxService := NewTaxService(&logger)

	// Call CalculationTax method
	taxResponse, err := taxService.CalculationTax(mockIncomeDetail)

	// Check for errors
	assert.NoError(t, err, "Error occurred while calculating tax")

	// Check the tax response
	expectedTaxAmount := 29000.0 // This value is based on the provided tax brackets and allowances
	assert.Equal(t, expectedTaxAmount, taxResponse.Tax, "Tax amount mismatch")
}

func TestCalculateTax(t *testing.T) {
	
	// Create an instance of TaxService
	logger := zerolog.Nop()
	taxService := NewTaxService(&logger)


	
	
	// Mock income and allowances
	income := 500000.0
	allowances := []Allowance{{AllowanceType: "", Amount: 0}, {AllowanceType: "", Amount: 0}}

	// Call calculateTax function
	taxAmount, err := taxService.calculateTax(income, allowances)

	// Check for errors
	assert.NoError(t, err, "Error occurred while calculating tax")

	// Check the calculated tax amount
	expectedTaxAmount := 29000.0 // This value is based on the provided tax brackets and allowances
	assert.Equal(t, expectedTaxAmount, taxAmount, "Tax amount mismatch")
}

func TestDeductPersonalAllowance(t *testing.T) {
	// Mock income and allowances
	income := 100000.0
	allowances := []Allowance{{AllowanceType: "standard", Amount: 5000}, {AllowanceType: "additional", Amount: 3000}}

	// Call deductPersonalAllowance function
	taxedIncome, err := deductPersonalAllowance(income, allowances)

	// Check for errors
	assert.NoError(t, err, "Error occurred while deducting personal allowance")

	// Check the taxed income
	expectedTaxedIncome := 32000.0 // This value is based on the provided income and allowances
	assert.Equal(t, expectedTaxedIncome, taxedIncome, "Taxed income mismatch")
}

func TestGetPersonalAllowance(t *testing.T) {
	// Call getPersonalAllowance function
	personalAllowance, err := getPersonalAllowance()

	// Check for errors
	assert.NoError(t, err, "Error occurred while getting personal allowance")

	// Check the personal allowance
	expectedPersonalAllowance := 60000.0 // This value is predefined
	assert.Equal(t, expectedPersonalAllowance, personalAllowance, "Personal allowance mismatch")
}

func TestGetTaxTable(t *testing.T) {
	// Call getTaxTable function
	brackets, err := getTaxTable()

	// Check for errors
	assert.NoError(t, err, "Error occurred while getting tax table")

	// Check the tax brackets
	expectedBrackets := []taxrepo.TaxBracket{
		{LowerBound: 0, UpperBound: 150000, TaxRate: 0.00},
		{LowerBound: 150001, UpperBound: 500000, TaxRate: 0.10},
		{LowerBound: 500001, UpperBound: 1000000, TaxRate: 0.15},
		{LowerBound: 1000001, UpperBound: 2000000, TaxRate: 0.20},
		{LowerBound: 2000001, UpperBound: 0, TaxRate: 0.35},
	}

	// Compare each bracket individually
	for i, bracket := range brackets {
		assert.Equal(t, expectedBrackets[i], bracket, "Tax bracket mismatch")
	}
}
