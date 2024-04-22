package tax

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
	Tax float64 `json:"tax"`
}


type Service interface{
	CalculationTax(*TaxRequest)(*TaxResponse,error)
}