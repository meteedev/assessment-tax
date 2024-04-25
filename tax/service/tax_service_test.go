package service

import (
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"

    "github.com/rs/zerolog"
    "github.com/meteedev/assessment-tax/constant"
    "github.com/meteedev/assessment-tax/tax/repository"
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

func TestCalculationTax(t *testing.T) {
    logger := zerolog.Logger{}
    mockRepo := new(MockTaxDeductConfigPort)
    taxService := NewTaxService(&logger, mockRepo)

    incomeDetail := &TaxRequest{
        TotalIncome: 500000,
        Allowances:  []Allowance{{AllowanceType: "donation", Amount: 200000}},
        WHT:         0,
    }

    mockRepo.On("FindById", "personal").Return(&repository.TaxDeductConfig{Amount: 60000}, nil)

    taxResponse, err := taxService.CalculationTax(incomeDetail)

    assert.NoError(t, err)
    assert.Equal(t, 19000.0, taxResponse.Tax)
}

func TestCalculateTax(t *testing.T) {
    logger := zerolog.Logger{}
    mockRepo := new(MockTaxDeductConfigPort)
    taxService := NewTaxService(&logger, mockRepo)



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
    logger := zerolog.Logger{}
    mockRepo := new(MockTaxDeductConfigPort)
    taxService := NewTaxService(&logger, mockRepo)

    updateReq := UpdateDeductRequest{Amount: 60000.0}

    mockRepo.On("UpdateById", constant.DEDUCT_PERSONAL_ID, 60000.0).Return(1, nil)
    mockRepo.On("FindById", "personal").Return(&repository.TaxDeductConfig{Amount: 60000.0}, nil)

    updDeductResponse, err := taxService.UpdatePersonalAllowance(&updateReq)

    assert.NoError(t, err)
    assert.Equal(t, 60000.0, updDeductResponse.Amount)
}



func TestDeductPersonalAllowance(t *testing.T) {
   	logger := zerolog.Logger{}
    mockRepo := new(MockTaxDeductConfigPort)
    
    mockRepo.On("FindById", constant.DEDUCT_PERSONAL_ID).Return(&repository.TaxDeductConfig{Amount: 60000.0}, nil)
    
    taxService := TaxService{&logger,mockRepo}



    income := 500000
    taxedIncome, err := taxService.deductPersonalAllowance(float64(income))

    assert.NoError(t, err)
    assert.Equal(t, 440000.0, taxedIncome)
}



func TestDeductWht(t *testing.T) {
    taxAmount := 3000.0
    wht := 500.0
    logger := zerolog.Logger{}
    mockRepo := new(MockTaxDeductConfigPort)
    taxService := TaxService{&logger,mockRepo}

    deductedWht := taxService.deductWht(taxAmount, wht)

    assert.Equal(t, 2500.0, deductedWht)
}



func TestAdjustMaximumDonationAllowanceDeduct(t *testing.T) {
    
    allowance := 2000.0

    logger := zerolog.Logger{}
    mockRepo := new(MockTaxDeductConfigPort)
    taxService := TaxService{&logger,mockRepo}

    adjusted := taxService.adjustMaximumDonationAllowanceDeduct(allowance)
    assert.Equal(t, 2000.0, adjusted)

}


func TestGetPersonalAllowance(t *testing.T) {
    // Mocking the logger
    logger := zerolog.Logger{}

    // Mocking the repository
    mockRepo := new(MockTaxDeductConfigPort)
    mockRepo.On("FindById", constant.DEDUCT_PERSONAL_ID).Return(&repository.TaxDeductConfig{Amount: 60000.0}, nil)

    // Creating a TaxService instance with the mocked logger and repository
    taxService := &TaxService{logger: &logger, DeductRepo: mockRepo}

    // Mocking the behavior of DeductRepo.FindById
    

    // Calling the method under test
    amount, err := taxService.getPersonalAllowance()

    // Assertions
    assert.NoError(t, err)
    assert.Equal(t, 60000.0, amount)
}
