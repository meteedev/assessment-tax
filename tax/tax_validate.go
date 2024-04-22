package tax

import (
	"errors"
	"strings"

	"github.com/meteedev/assessment-tax/constant"
)


func ValidateTaxRequest(taxRequest *TaxRequest) error {
	//fmt.Println("taxRequest.WHT ",taxRequest.WHT)
	var errMsgs []string

	validateWht(taxRequest.WHT,taxRequest.TotalIncome, &errMsgs)
	
	if len(errMsgs) > 0 {
		return errors.New(strings.Join(errMsgs, "; "))
	}

	return nil
}



func validateWht(wht float64,totalIncome float64,errMsgs *[]string){
	validateWhtGreaterThanOrEqualZero(wht , errMsgs)
	validateWhtNotGreaterThanTotalIncome(wht,totalIncome,errMsgs)
}

func validateWhtGreaterThanOrEqualZero(wht float64, errMsgs *[]string) {		
	//fmt.Println("wht ",wht)
	if wht < 0 {
		*errMsgs = append(*errMsgs, constant.MSG_BU_INVALID_WHT_LESS_THAN_ZERO)
	}
}

func validateWhtNotGreaterThanTotalIncome(wht float64,totalIncome float64, errMsgs *[]string) {		
	//fmt.Println("wht ",wht)
	if wht > totalIncome {
		*errMsgs = append(*errMsgs, constant.MSG_BU_INVALID_WHT_GREATER_THAN_TOTALINCOME) 
	}
}
