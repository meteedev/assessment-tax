package service

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestTaxService_calculateTaxTable(t *testing.T) {
	taxService := TaxService{}

	tests := []struct {
		name     string
		salary   float64
		expected float64
	}{
		{"No Tax", 100000.0, 0.0},
		{"Tax in first bracket", 440000.0, 29000.0},
		{"Tax in second bracket", 600000.0, 60000.0},
		{"Tax in third bracket", 1200000.0, 250000.0},
		{"Tax in fourth bracket", 2500000.0, 1010000.00},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, totalTax := taxService.calculateTaxTable(test.salary)
			assert.Equal(t, test.expected, totalTax, "For %s: Total tax amount mismatch", test.name)
		})
	}
}

func TestTaxService_calculateStep(t *testing.T) {
	taxService := TaxService{}

	tests := []struct {
		name       string
		income     float64
		lowerBound float64
		rate       float64
		level      string
		expected   float64
	}{
		{"No Tax", 100000.0, 0, 0, "0-150,000", 0},
		{"Tax in range", 200000.0, 0, 0.1, "0-150,000", 20000},
		{"Tax beyond range", 300000.0, 0, 0.1, "0-150,000", 30000},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tax, step := taxService.calculateStep(test.income, test.lowerBound, test.rate, test.level)
			assert.Equal(t, test.expected, tax, "For %s: Tax amount mismatch", test.name)
			assert.Equal(t, test.expected, step.TaxAmount, "For %s: Step tax amount mismatch", test.name)
		})
	}
}