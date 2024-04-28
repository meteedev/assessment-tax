package service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/meteedev/assessment-tax/constant"
	"github.com/meteedev/assessment-tax/tax/repository"
	"github.com/rs/zerolog"
)

// MockTaxDeductConfigPort is a mock implementation of the repository.TaxDeductConfigPort interface.
type MockTaxDeductConfigPort struct {
    mock.Mock
}

func (m *MockTaxDeductConfigPort) FindById(id string) (*repository.TaxDeductConfig, error) {
    args := m.Called(id)
    return args.Get(0).(*repository.TaxDeductConfig), args.Error(1)
}



func (m *MockTaxDeductConfigPort) UpdateById(id string, amount float64) (int64, error) {
    args := m.Called(id, amount)
	return 1, args.Error(1)
}

func TestCalculationTax_deduct_donation(t *testing.T) {
    logger := &zerolog.Logger{}
    mockRepo := new(MockTaxDeductConfigPort)
    csvPaser := &CSVParserImpl{}
    taxService := NewTaxService(logger, mockRepo,csvPaser)

    incomeDetail := &TaxRequest{
        TotalIncome: 500000,
        Allowances:  []Allowance{{AllowanceType: "donation", Amount: 200000} ,  },
        WHT:         0,
    }

    mockRepo.On("FindById", "personal").Return(&repository.TaxDeductConfig{Amount: 60000}, nil)
  

    taxResponse, err := taxService.CalculationTax(incomeDetail)

    assert.NoError(t, err)
    assert.Equal(t, 19000.0, taxResponse.Tax)
}


func TestCalculationTax_deduct_Kreceipt(t *testing.T) {
    logger := &zerolog.Logger{}
    mockRepo := new(MockTaxDeductConfigPort)
    csvPaser := &CSVParserImpl{}
    taxService := NewTaxService(logger, mockRepo,csvPaser)

    incomeDetail := &TaxRequest{
        TotalIncome: 500000,
        Allowances:  []Allowance{ {AllowanceType: "k-receipt", Amount: 200000} },
        WHT:         0,
    }

    mockRepo.On("FindById", "k-receipt").Return(&repository.TaxDeductConfig{Amount: 50000.0}, nil)
    mockRepo.On("FindById", "personal").Return(&repository.TaxDeductConfig{Amount: 60000.0}, nil)

    taxResponse, err := taxService.CalculationTax(incomeDetail)

    assert.NoError(t, err)
    assert.Equal(t, 24000.0, taxResponse.Tax)
}

func TestCalculateTax(t *testing.T) {
    logger := &zerolog.Logger{}
    mockRepo := new(MockTaxDeductConfigPort)
    csvPaser := &CSVParserImpl{}
    taxService := NewTaxService(logger, mockRepo,csvPaser)



    incomeDetail := &TaxRequest{
        TotalIncome: 500000,
        Allowances:  []Allowance{{AllowanceType: "donation",Amount: 0}},
        WHT:         25000,
    }

    mockRepo.On("FindById", "personal").Return(&repository.TaxDeductConfig{Amount: 60000}, nil)

	taxResponse, err := taxService.CalculationTax(incomeDetail)	

    assert.NoError(t, err)
    assert.Equal(t, 4000.0, taxResponse.Tax)
}



func TestUpdatePersonalAllowance(t *testing.T) {
    logger := &zerolog.Logger{}
    mockRepo := new(MockTaxDeductConfigPort)
    csvPaser := &CSVParserImpl{}
    taxService := NewTaxService(logger, mockRepo,csvPaser)

    updateReq := UpdateDeductRequest{Amount: 60000.0}

    mockRepo.On("UpdateById", constant.DEDUCT_PERSONAL_ID, 60000.0).Return(1, nil)
    mockRepo.On("FindById", constant.DEDUCT_PERSONAL_ID).Return(&repository.TaxDeductConfig{Amount: 60000.0}, nil)

    updDeductResponse, err := taxService.UpdatePersonalAllowance(&updateReq)

    assert.NoError(t, err)
    assert.Equal(t, 60000.0, updDeductResponse.Amount)
}


func TestUpdateKreceiptAllowance(t *testing.T) {
    logger := &zerolog.Logger{}
    mockRepo := new(MockTaxDeductConfigPort)
    csvPaser := &CSVParserImpl{}
    taxService := NewTaxService(logger, mockRepo,csvPaser)

    updateReq := UpdateDeductRequest{Amount: 60000.0}

    mockRepo.On("UpdateById", constant.DEDUCT_K_RECEIPT_ID, 60000.0).Return(1, nil)
    mockRepo.On("FindById", constant.DEDUCT_K_RECEIPT_ID).Return(&repository.TaxDeductConfig{Amount: 60000.0}, nil)

    updDeductResponse, err := taxService.UpdateKreceiptAllowance(&updateReq)

    assert.NoError(t, err)
    assert.Equal(t, 60000.0, updDeductResponse.Amount)
}


// func TestDeductPersonalAllowance(t *testing.T) {
//    	logger := &zerolog.Logger{}
//     mockRepo := new(MockTaxDeductConfigPort)
    
//     mockRepo.On("FindById", constant.DEDUCT_PERSONAL_ID).Return(&repository.TaxDeductConfig{Amount: 60000.0}, nil)
    
//     csvPaser := &CSVParserImpl{}
//     taxService := NewTaxService(logger, mockRepo,csvPaser)


//     income := 500000
//     taxedIncome, err := taxService.deductPersonalAllowance(float64(income))

//     assert.NoError(t, err)
//     assert.Equal(t, 440000.0, taxedIncome)
// }



// func TestDeductWht(t *testing.T) {
//     taxAmount := 3000.0
//     wht := 500.0
//     logger := &zerolog.Logger{}
//     mockRepo := new(MockTaxDeductConfigPort)
    
//     csvPaser := &CSVParserImpl{}
//     taxService := NewTaxService(logger, mockRepo,csvPaser)



//     deductedWht := taxService.DeductWht(taxAmount, wht)

//     assert.Equal(t, 2500.0, deductedWht)
// }






func TestGetPersonalAllowance(t *testing.T) {
    // Mocking the logger
    logger := &zerolog.Logger{}
    mockRepo := new(MockTaxDeductConfigPort)
    csvPaser := &CSVParserImpl{}
    mockRepo.On("FindById", constant.DEDUCT_PERSONAL_ID).Return(&repository.TaxDeductConfig{Amount: 60000.0}, nil)

    // Creating a TaxService instance with the mocked logger and repository
    taxService := TaxService{logger: logger, DeductRepo: mockRepo,csvParser: csvPaser}

    // Mocking the behavior of DeductRepo.FindById
    

    // Calling the method under test
    amount, err := taxService.getPersonalAllowance()

    // Assertions
    assert.NoError(t, err)
    assert.Equal(t, 60000.0, amount)
}



func TestGetKreceiptAllowance(t *testing.T) {
    // Mocking the logger
    logger := zerolog.Logger{}

    // Mocking the repository
    mockRepo := new(MockTaxDeductConfigPort)
    mockRepo.On("FindById", constant.DEDUCT_K_RECEIPT_ID).Return(&repository.TaxDeductConfig{Amount: 70000.0}, nil)

    // Creating a TaxService instance with the mocked logger and repository
    taxService := &TaxService{logger: &logger, DeductRepo: mockRepo}


    // Calling the method under test
    amount, err := taxService.getKreceiptAllowance()

    // Assertions
    assert.NoError(t, err)
    assert.Equal(t, 70000.0, amount)
}



// func TestAdjustMaximumKreceiptAllowanceDeduct(t *testing.T) {
    
//     kreciptAllowanceRequest := 100000.0

//     logger := &zerolog.Logger{}
//     mockRepo := new(MockTaxDeductConfigPort)
//     csvPaser := &CSVParserImpl{}
//     taxService := NewTaxService(logger, mockRepo,csvPaser)

//     mockRepo.On("FindById", constant.DEDUCT_K_RECEIPT_ID).Return(&repository.TaxDeductConfig{Amount: 70000.0}, nil)
//     adjusted,err := taxService.adjustMaximumKreceiptAllowanceDeduct(kreciptAllowanceRequest)
        
//     assert.NoError(t, err)
//     assert.Equal(t, 70000.0, adjusted)
// }

func TestGetTaxUpload(t *testing.T) {
	// Create sample TaxRequest and TaxResponse
	taxRequest := &TaxRequest{
		TotalIncome: 1000,
		WHT:         200,
		Allowances: []Allowance{
			{AllowanceType: "donation", Amount: 50},
			{AllowanceType: "k-receipt", Amount: 100},
		},
	}
	taxResponse := &TaxResponse{
		Tax:       150,
		TaxRefund: 50,
	}

	// Call the function to get TaxUpload
	taxUpload := getTaxUpload(taxRequest, taxResponse)

	// Assert the values in the returned TaxUpload struct
	expectedTaxUpload := TaxUpload{
		TotalIncome: taxRequest.TotalIncome,
		Tax:         taxResponse.Tax,
		TaxRefund:   taxResponse.TaxRefund,
	}
	assert.Equal(t, expectedTaxUpload, taxUpload, "TaxUpload struct does not match expected values")
}

func TestDeductWht(t *testing.T) {
	taxService := &TaxService{}

	testCases := []struct {
		name          string
		taxAmount     float64
		wht           float64
		expectedDiff  float64
	}{
		{
			name:          "PositiveTaxAmount",
			taxAmount:     1000,
			wht:           200,
			expectedDiff:  800, // Tax amount (1000) - WHT (200)
		},
		{
			name:          "ZeroTaxAmount",
			taxAmount:     0,
			wht:           100,
			expectedDiff:  -100, // Negative because WHT is greater than tax amount
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			diff := taxService.deductWht(tc.taxAmount, tc.wht)
			assert.Equal(t, tc.expectedDiff, diff)
		})
	}
}

func TestAdjustMaximumDonationAllowanceDeduct(t *testing.T) {
	taxService := &TaxService{}

	testCases := []struct {
		name              string
		allowance         float64
		expectedAdjusted  float64
	}{
		{
			name:             "AllowanceBelowMax",
			allowance:        50,
			expectedAdjusted: 50, // No adjustment needed
		},
		{
			name:             "AllowanceAboveMax",
			allowance:        200000,
			expectedAdjusted: constant.MAX_ALLOWANCE_DONATION_DEDUCT, // Should be adjusted to the max allowance
		},
		{
			name:             "AllowanceEqualMax",
			allowance:        constant.MAX_ALLOWANCE_DONATION_DEDUCT,
			expectedAdjusted: constant.MAX_ALLOWANCE_DONATION_DEDUCT, // No adjustment needed, already at max
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			adjusted := taxService.adjustMaximumDonationAllowanceDeduct(tc.allowance)
			assert.Equal(t, tc.expectedAdjusted, adjusted)
		})
	}
}

type mockDeductRepo struct{}

func (m *mockDeductRepo) FindById(id string) (*repository.TaxDeductConfig, error) {
	// Check the ID to determine which config to return
	if id == constant.DEDUCT_K_RECEIPT_ID {
		// Mock implementation to return the fixed value for K-receipt config
		return &repository.TaxDeductConfig{Amount: 50000}, nil
	}
	// Return nil for other IDs
	return nil, errors.New("unexpected ID")
}

func (m *mockDeductRepo) UpdateById(id string, amount float64) (int64, error) {
	// Mock implementation for the UpdateById method
	return 1, nil
}

func TestAdjustMaximumKreceiptAllowanceDeduct(t *testing.T) {
	mockDeductRepo := &mockDeductRepo{}

	mockTaxService := &TaxService{
		DeductRepo: mockDeductRepo,
	}

	testCases := []struct {
		name                    string
		allowance               float64
		expectedAdjusted        float64
		expectedErr             error
	}{
		{
			name:             "AllowanceBelowMax",
			allowance:        40000,
			expectedAdjusted: 40000, // No adjustment needed
			expectedErr:      nil,
		},
		{
			name:             "AllowanceAboveMax",
			allowance:        60000,
			expectedAdjusted: 50000, // Should be adjusted to the max allowance (50000)
			expectedErr:      nil,
		},
		{
			name:             "AllowanceEqualMax",
			allowance:        50000,
			expectedAdjusted: 50000, // No adjustment needed, already at max
			expectedErr:      nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			adjusted, err := mockTaxService.adjustMaximumKreceiptAllowanceDeduct(tc.allowance)
			assert.Equal(t, tc.expectedAdjusted, adjusted)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}