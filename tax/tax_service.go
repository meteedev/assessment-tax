package tax

import (
	"fmt"
	"log"

	"github.com/meteedev/assessment-tax/apperrs"
	"github.com/meteedev/assessment-tax/constant"
	taxrepo "github.com/meteedev/assessment-tax/repository"
)

type TaxService struct {
}

func NewTaxService() *TaxService {
	return &TaxService{}
}

func (t *TaxService) CalculationTax(incomeDetail *TaxRequest) (*TaxResponse, error) {
	
	income := incomeDetail.TotalIncome
	allowances := incomeDetail.Allowances


	// Calculate tax
	taxAmount, err := calculateTax(income,allowances)

	if err != nil{
		return nil,err
	}

	// Print the result
	fmt.Printf("Income: %.2f\nTax Amount: %.2f\n", income, taxAmount)

	taxResponse := TaxResponse{
		Tax: taxAmount,
	}

	return &taxResponse, nil

}


func calculateTax(income float64, allowances []Allowance) (float64, error) {
    var taxAmount float64

    brackets, err := getTaxTable()
    if err != nil {
        log.Println(err)
		return 0, apperrs.NewInternalServerError(constant.MSG_BU_GERNERAL_ERROR)
    }

    taxedIncome, err := deductPersonalAllowance(income, allowances)
    if err != nil {
        log.Println(err)
		return 0, apperrs.NewInternalServerError(constant.MSG_BU_GERNERAL_ERROR)
    }

    remainingIncome := taxedIncome

    for _, bracket := range brackets {

        if remainingIncome <= 0 {
            break // No more taxable income
        }

		lower := bracket.LowerBound-1

		if(lower <0){
			lower =0
		}

        taxableAmount := min(remainingIncome, bracket.UpperBound) - lower
        if taxableAmount <= 0 {
            continue // No taxable amount in this bracket
        }

        // Calculate tax for this bracket
        bracketTax := taxableAmount * bracket.TaxRate
        taxAmount += bracketTax

        // Update remaining income for the next iteration
        remainingIncome -= taxableAmount
    }

    return taxAmount, nil
}



func deductPersonalAllowance(income float64 , allowances []Allowance)(float64,error){
	
	totalAllowance := 0.0
	for _, allowance := range allowances {
		totalAllowance += allowance.Amount
	}


	personalAllowance,err := getPersonalAllowance()
	if err != nil{
		log.Println(err)
		return 0, apperrs.NewInternalServerError(constant.MSG_BU_GERNERAL_ERROR)
	}

	taxedIncome := income - totalAllowance - personalAllowance
	
	return taxedIncome,nil
}


func getPersonalAllowance()(float64,error){
	return 60000.0,nil
}

func getTaxTable() ([]taxrepo.TaxBracket,error) {

	brackets := []taxrepo.TaxBracket{
		{LowerBound: 0, UpperBound: 150000, TaxRate: 0.10}, // Adjust the tax rate as needed
		{LowerBound: 150001, UpperBound: 500000, TaxRate: 0.10},
		{LowerBound: 500001, UpperBound: 1000000, TaxRate: 0.15},
		{LowerBound: 1000001, UpperBound: 2000000, TaxRate: 0.20},
		{LowerBound: 2000001, UpperBound: 0, TaxRate: 0.35},
	}

	return brackets,nil;

}