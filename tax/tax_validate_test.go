package tax

import (
	"testing"
)

func TestValidateTaxRequest(t *testing.T) {
	testCases := []struct {
		name      string
		taxRequest *TaxRequest
		expectErr bool
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
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateTaxRequest(tc.taxRequest)
			if tc.expectErr && err == nil {
				t.Errorf("expected error but got nil")
			}
			if !tc.expectErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
