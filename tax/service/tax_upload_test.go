package service

import (
	"strings"
	"testing"
	"io"	

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/rs/zerolog"
	"github.com/meteedev/assessment-tax/tax/repository"
)

func TestParseTaxRequestRecord(t *testing.T) {
	tests := []struct {
		record   []string
		expected *TaxRequest
		err      error
	}{
		// Test case: Valid input
		{
			record: []string{"1000", "200", "50"},
			expected: &TaxRequest{
				TotalIncome: 1000,
				WHT:         200,
				Allowances: []Allowance{
					{AllowanceType: "donation", Amount: 50},
				},
			},
			err: nil,
		},
	}

	for _, test := range tests {
		result, err := parseTaxRequestRecord(test.record)
		assert.Equal(t, test.err, err, "Error mismatch")
		assert.Equal(t, test.expected, result, "Result mismatch")
	}
}

func TestParseCSVToTaxRequest(t *testing.T) {
	testCases := []struct {
		name          string
		csvData       string
		expectedError error
	}{
		{
			name:    "Valid CSV",
			csvData: "1000,50,25\n2000,75,30\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			csvParser := CSVParserImpl{}
			csvData := strings.NewReader(tc.csvData)
			
			taxRequests, err := csvParser.ParseCSVToTaxRequest(csvData)
			
			assert.NoError(t, err)
			assert.NotNil(t, taxRequests)
			assert.NotEmpty(t, *taxRequests)
			
		})
	}
}


type MockCSVParser struct {
	mock.Mock
}

// ParseCSVToTaxRequest mocks parsing CSV to tax request.
func (m *MockCSVParser) ParseCSVToTaxRequest(file io.Reader) (*[]TaxRequest, error) {
	args := m.Called(file)
	return args.Get(0).(*[]TaxRequest), args.Error(1)
}





func TestUploadCalculationTax(t *testing.T) {
	// Mock dependencies
	logger := &zerolog.Logger{}
	mockRepo := new(MockTaxDeductConfigPort)
	mockCSVParser := new(MockCSVParser)
	mockTaxService := NewTaxService(logger, mockRepo, mockCSVParser)

	// Test cases
	testCases := []struct {
		name          string
		csvData       io.Reader
		mockReturn    *[]TaxRequest
		mockErr       error
		expectedError error
	}{
		{
			name:    "Valid CSV",
			csvData: strings.NewReader("1000,50,25\n2000,75,30\n"),
			mockReturn: &[]TaxRequest{
				{TotalIncome: 500000, WHT: 0, Allowances: []Allowance{{AllowanceType: "donation", Amount: 0}}},
			},
			mockErr:       nil,
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCSVParser.On("ParseCSVToTaxRequest", tc.csvData).Return(tc.mockReturn, tc.mockErr).Once()
			mockRepo.On("FindById", "personal").Return(&repository.TaxDeductConfig{Amount: 60000}, nil)
			response, err := mockTaxService.UploadCalculationTax(tc.csvData)

			if tc.expectedError != nil {
				assert.EqualError(t, err, tc.expectedError.Error())
				assert.Nil(t, response)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, response)
			}

			mockCSVParser.AssertExpectations(t)
		})
	}
}