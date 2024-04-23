package service

type Allowance struct {
	AllowanceType string  `json:"allowanceType"`
	Amount        float64 `json:"amount"`
}

type TaxRequest struct {
	TotalIncome float64     `json:"totalIncome"`
	WHT         float64     `json:"wht"`
	Allowances  []Allowance `json:"allowances"`
}


type TaxResponse struct {
	Tax float64 		  	`json:"tax"`
	TaxBracket []TaxBracket `json:"taxLevel"`
}


type TaxBracket struct {
	Level 		string		`json:"level"`	
	LowerBound 	float64 	`json:"-"`
	UpperBound 	float64 	`json:"-"`
	TaxRate    	float64 	`json:"-"`
	Tax 		float64 	`json:"tax"`
}


type UpdateDeductRequest struct {
	Amount 		float64		`json:"amount"`	
}

type UpdateDeductResponse struct {
	Amount 		float64		`json:"amount"`	
}



type TaxServicePort interface{
	CalculationTax(*TaxRequest)(*TaxResponse,error)
	UpdatePersonalAllowance(*UpdateDeductRequest)(*UpdateDeductResponse,error)
}

