package repository

type TaxDeductConfig struct {
    DeductId  string  `json:"deduct_type"`
    Amount      float64 `json:"amount"`
    Description string  `json:"description"`
}

type TaxDeductConfigPort interface {
	FindById(id string) (*TaxDeductConfig,error)
    UpdateById(id string,amount float64) (int64,error)
}