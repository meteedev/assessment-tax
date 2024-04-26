package service

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/meteedev/assessment-tax/constant"
)

func TestValidateTotalIncomeGreaterThanOrEqualZero(t *testing.T) {
    // Test case 1: amount is greater than zero
    amount := 100.0
    errMsgs := []string{}
    validateTotalIncomeGreaterThanOrEqualZero(amount, &errMsgs)
    assert.Empty(t, errMsgs, "Expected no error messages for amount greater than zero")

    // Test case 2: amount is equal to zero
    amount = 0
    errMsgs = []string{}
    validateTotalIncomeGreaterThanOrEqualZero(amount, &errMsgs)
    assert.NotEmpty(t, errMsgs, "Expected error message for amount equal to zero")
    assert.Equal(t, errMsgs[0], constant.MSG_BU_INVALID_TOTAL_INCOME_LESS_THAN_OR_EQUAL_ZERO, "Incorrect error message for amount equal to zero")

    // Test case 3: amount is less than zero
    amount = -50.0
    errMsgs = []string{}
    validateTotalIncomeGreaterThanOrEqualZero(amount, &errMsgs)
    assert.NotEmpty(t, errMsgs, "Expected error message for amount less than zero")
    assert.Equal(t, errMsgs[0], constant.MSG_BU_INVALID_TOTAL_INCOME_LESS_THAN_OR_EQUAL_ZERO, "Incorrect error message for amount less than zero")
}

func TestValidateTaxRequest(t *testing.T) {
	testCases := []struct {
		name         string
		taxRequest   *TaxRequest
		expectErr    bool
		expectedMsgs []string
	}{
		{
			name: "ValidTaxRequest",
			taxRequest: &TaxRequest{
				WHT:         100.0,
				TotalIncome: 1000.0,
			},
			expectErr: false,
		},
		{
			name: "NegativeWHT",
			taxRequest: &TaxRequest{
				WHT:         -100.0,
				TotalIncome: 1000.0,
			},
			expectErr: true,
		},
		{
			name: "WHTGreaterThanTotalIncome",
			taxRequest: &TaxRequest{
				WHT:         1100.0,
				TotalIncome: 1000.0,
			},
			expectErr: true,
		},

		{
			name: "InvalidAllowanceType",
			taxRequest: &TaxRequest{
				TotalIncome: 1000.0,
				WHT:         100.0,
				Allowances: []Allowance{
					{AllowanceType: "invalid", Amount: 200.0},
				},
			},
			expectErr:    true,
			expectedMsgs: []string{"donation", "k-receipt"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateTaxRequest(tc.taxRequest)

			if tc.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if len(tc.expectedMsgs) != 0 {
				for _, expectedMsg := range tc.expectedMsgs {
					assert.Contains(t, err.Error(), expectedMsg)
				}
			}
		})
	}
}

func TestOnlyDigits(t *testing.T) {
	// Initialize testify's assert package
	assert := assert.New(t)

	// Test case: string contains only digits
	var errMsgs []string
	onlyDigits("12345", &errMsgs)
	assert.Empty(errMsgs, "Expected no error messages for string containing only digits")

	// Test case: string contains non-digit characters
	errMsgs = nil
	onlyDigits("12a45", &errMsgs)
	expectedErrMsg := constant.MSG_BU_VALIDATE_DIGIT_ONLY
	assert.Equal([]string{expectedErrMsg}, errMsgs, "Expected error message for string containing non-digit characters")
}