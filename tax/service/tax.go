package service

import(
	"mime/multipart"
)

type TaxServicePort interface{
	CalculationTax(*TaxRequest)(*TaxResponse,error)
	UploadCalculationTax(file *multipart.FileHeader)(*TaxUploadResponse,error)
	UpdatePersonalAllowance(*UpdateDeductRequest)(*UpdateDeductResponse,error)
	UpdateKreceiptAllowance(*UpdateDeductRequest)(*UpdateDeductResponse,error)
}

type TaxRequest struct {
	TotalIncome float64     `json:"totalIncome"`
	WHT         float64     `json:"wht"`
	Allowances  []Allowance `json:"allowances"`
}

type Allowance struct {
	AllowanceType string  `json:"allowanceType"`
	Amount        float64 `json:"amount"`
}




type TaxResponse struct {
	Tax 		float64 	`json:"tax"`
	TaxRefund	float64		`json:"taxRefund"`
	TaxStep 	[]TaxStep 	`json:"taxLevel"`
}


type TaxBracket struct {
	Level 		string		`json:"level"`	
	LowerBound 	float64 	`json:"-"`
	UpperBound 	float64 	`json:"-"`
	TaxRate    	float64 	`json:"-"`
	Tax 		float64 	`json:"tax"`
}

type TaxStep struct {
	Level     string	`json:"level"`
	TaxAmount float64  	`json:"tax"`
}


type UpdateDeductRequest struct {
	Amount 		float64		`json:"amount"`	
}

type UpdateDeductResponse struct {
	Amount 		float64		`json:"amount"`	
}


type TaxUpload struct {
    TotalIncome float64 `json:"totalIncome"`
    Tax         float64 `json:"tax"`
	TaxRefund	float64 `json:"taxRefund"`
}

type TaxUploadResponse struct {
    Taxes []TaxUpload `json:"taxes"`
}



