package service

import (
	"errors"
	"fmt"
	"testing"

	"github.com/meteedev/assessment-tax/constant"
	"github.com/stretchr/testify/assert"
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
	expectedErrMsg := constant.MSG_BU_VALIDATE_CSV_DIGIT_ONLY
	assert.Equal([]string{expectedErrMsg}, errMsgs, "Expected error message for string containing non-digit characters")
}


func TestValidateUploadTaxCsvRecord(t *testing.T) {
	tests := []struct {
		record     []string
		expectedErr error
	}{
		{record: []string{"100", "200", "300"}, expectedErr: nil},         // Valid record
		{record: []string{"100", "200", "not a number"}, expectedErr: errors.New(constant.MSG_BU_VALIDATE_CSV_DIGIT_ONLY)}, 
		{record: []string{"100", "200"}, expectedErr: errors.New(constant.MSG_BU_INVALID_CSV_RECORD_COLUMN_NUMBERS)}, 
		
	}

	for _, test := range tests {
		err := ValidateUploadTaxCsvRecord(test.record)
		assert.Equal(t, test.expectedErr, err, "For record %v, expected error: %v, but got: %v", test.record, test.expectedErr, err)
	}
}


func TestValidatePersonalAllowanceMinimum(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		amount     float64
		expected   []string
	}{
		{
			name:       "Amount less than minimum",
			amount:     100,
			expected:   []string{fmt.Sprintf("Personal deductibles start at %.2f baht", constant.MIN_ALLOWANCE_PERSONAL)},
		},
		{
			name:       "Amount equal to minimum",
			amount:     constant.MIN_ALLOWANCE_PERSONAL,
			expected:   []string{},
		},
		{
			name:       "Amount greater than minimum",
			amount:     20000,
			expected:    []string{},
		},
	}

	for _, tc := range tests {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			errMsgs := make([]string, 0)
			validatePersonalAllowanceMinimum(tc.amount, &errMsgs)
			assert.Equal(t, tc.expected, errMsgs)
		})
	}
}


func TestValidatePersonalAllowanceMaximum(t *testing.T) {
	
	tests := []struct {
		name       string
		amount     float64
		expected   []string
	}{
		{
			name:       "Amount exceeds maximum",
			amount:     100001,
			expected:   []string{fmt.Sprintf("Maximum Personal deductibles %.2f baht",  constant.MAX_ALLOWANCE_PERSONAL)},
		},
	
	}

	for _, tc := range tests {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			errMsgs := make([]string, 0)
			validatePersonalAllowanceMaximum(tc.amount, &errMsgs)
			assert.Equal(t, tc.expected, errMsgs)
		})
	}
}