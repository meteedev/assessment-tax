package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/meteedev/assessment-tax/tax/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) CalculationTax(incomeDetail *service.TaxRequest) (*service.TaxResponse, error) {
	args := m.Called(incomeDetail)
	return args.Get(0).(*service.TaxResponse), args.Error(1)
}

func (m *MockService) UpdatePersonalAllowance(updateDeductRequest *service.UpdateDeductRequest)(*service.UpdateDeductResponse,error){
	args := m.Called(updateDeductRequest)
	return args.Get(0).(*service.UpdateDeductResponse), args.Error(1)
}

func (m *MockService) UpdateKreceiptAllowance(updateDeductRequest *service.UpdateDeductRequest)(*service.UpdateDeductResponse,error){
	args := m.Called(updateDeductRequest)
	return args.Get(0).(*service.UpdateDeductResponse), args.Error(1)
}

func (m *MockService) UploadCalculationTax(taxRequests *[]service.TaxRequest)(*service.TaxUploadResponse,error){
	args := m.Called(taxRequests)
	return args.Get(0).(*service.TaxUploadResponse), args.Error(1)
}

func TestTaxCalculationsHandler(t *testing.T) {
	// Create a new instance of the mock service
	mockService := new(MockService)

	// Create a new instance of the handler with the mock service
	handler := NewTaxHandler(mockService)

	// Create a new echo context for testing
	e := echo.New()
	reqBody := []byte(`{"TotalIncome":100000,"Allowances":[{"Amount":5000},{"Amount":3000}]}`)
	req := httptest.NewRequest(http.MethodPost, "/tax/calculations", bytes.NewBuffer(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Mock the service method
	expectedResponse := &service.TaxResponse{Tax: 15000.0}
	mockService.On("CalculationTax", mock.Anything).Return(expectedResponse, nil)

	// Call the handler method
	err := handler.TaxCalculation(c)
	assert.NoError(t, err)

	// Check the response
	assert.Equal(t, http.StatusOK, rec.Code)

	var response service.TaxResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse.Tax, response.Tax)

	// Assert that the mock service method was called with the correct arguments
	mockService.AssertCalled(t, "CalculationTax", mock.Anything)
}


func TestDeductionsPersonalHandler(t *testing.T) {
	// Create a new instance of the mock service
	mockService := new(MockService)

	// Create a new instance of the handler with the mock service
	handler := NewTaxHandler(mockService)

	// Create a new echo context for testing
	e := echo.New()
	reqBody := []byte(`{"amount":60000 }`)
	req := httptest.NewRequest(http.MethodPost, "/deductions/personal", bytes.NewBuffer(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Mock the service method
	expectedResponse := &service.UpdateDeductResponse{Amount: 60000.0}
	mockService.On("UpdatePersonalAllowance", mock.Anything).Return(expectedResponse, nil)

	// Call the handler method
	err := handler.DeductionsPersonal(c)
	assert.NoError(t, err)

	// Check the response
	assert.Equal(t, http.StatusOK, rec.Code)

	var response service.UpdateDeductResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse.Amount, response.Amount)

	// Assert that the mock service method was called with the correct arguments
	mockService.AssertCalled(t, "UpdatePersonalAllowance", mock.Anything)
}

func TestDeductionsKreceiptHandler(t *testing.T) {
	// Create a new instance of the mock service
	mockService := new(MockService)

	// Create a new instance of the handler with the mock service
	handler := NewTaxHandler(mockService)

	// Create a new echo context for testing
	e := echo.New()
	reqBody := []byte(`{"amount":70000 }`)
	req := httptest.NewRequest(http.MethodPost, "/deductions/k-receipt", bytes.NewBuffer(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Mock the service method
	expectedResponse := &service.UpdateDeductResponse{Amount: 70000.0}
	mockService.On("UpdateKreceiptAllowance", mock.Anything).Return(expectedResponse, nil)

	// Call the handler method
	err := handler.DeductionsKreceipt(c)
	assert.NoError(t, err)

	// Check the response
	assert.Equal(t, http.StatusOK, rec.Code)

	var response service.UpdateDeductResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse.Amount, response.Amount)

	// Assert that the mock service method was called with the correct arguments
	mockService.AssertCalled(t, "UpdateKreceiptAllowance", mock.Anything)
}