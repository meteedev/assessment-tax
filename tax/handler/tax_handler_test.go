package handler

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
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

func (m *MockService) UploadCalculationTax(file io.Reader)(*service.TaxUploadResponse,error){
	args := m.Called(file)
	return args.Get(0).(*service.TaxUploadResponse), args.Error(1)
}

func TestTaxCalculationsHandler(t *testing.T) {
	// Create a new instance of the mock service
	mockService := new(MockService)

	// Create a new instance of the handler with the mock service
	handler := NewTaxHandler(mockService)

	// Create a new echo context for testing
	e := echo.New()
	reqBody := []byte(`{"totalIncome":100000,"wht":0.0,"allowances":[{"allowanceType":"donations","amount":0.0}]}`)
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

func TestTaxHandler_TaxUploadCalculation(t *testing.T) {
	// Create a new instance of the Echo framework
	e := echo.New()

	// Create a mock service
	mockService := new(MockService)

	// Create a new instance of the TaxHandler with the mock service
	taxHandler := TaxHandler{
		service: mockService,
	}

	// Create a mock HTTP request with a test CSV file
	fileContents := "TotalIncome,WHT,Donation\n1000,100,50\n2000,200,100\n"
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("taxFile", "test.csv")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := io.Copy(part, strings.NewReader(fileContents)); err != nil {
		t.Fatal(err)
	}
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Mock the UploadCalculationTax method of the service
	mockResponse := &service.TaxUploadResponse{}
	mockService.On("UploadCalculationTax", mock.Anything).Return(mockResponse, nil).Run(func(args mock.Arguments) {
		// Assert that the file passed to the service is the same as the one received by the handler
		file := args.Get(0).(io.Reader)
		fileBytes, err := io.ReadAll(file)
		assert.NoError(t, err)
		assert.Equal(t, fileContents, string(fileBytes))
	})

	// Call the handler function
	err = taxHandler.TaxUploadCalculation(c)

	// Assert that there was no error
	assert.NoError(t, err)

	// Assert the status code is OK
	assert.Equal(t, http.StatusOK, rec.Code)


	// Assert that the UploadCalculationTax method was called with the correct argument
	mockService.AssertCalled(t, "UploadCalculationTax", mock.Anything)
}