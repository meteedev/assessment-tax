package service

import (
	"testing"

	"github.com/stretchr/testify/assert"

)

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