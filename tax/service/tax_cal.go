package service



func (t *TaxService) getStep(level string, amount float64) TaxStep {
	return TaxStep{Level: level, TaxAmount: amount}
}

func (t *TaxService) calculateStep(income, lowerBound, rate float64, level string) (float64, TaxStep) {
	var tax float64
	if income <= lowerBound {
		//fmt.Printf("income:%.2f lowerBound:%.2f rate:%.2f tax:%.2f \n",income,lowerBound,rate,0.0)
		return 0, t.getStep(level, 0)
	}


	tax = (income - lowerBound) * rate
	//fmt.Printf("income:%.2f lowerBound:%.2f rate:%.2f tax:%.2f \n",income,lowerBound,rate,tax)

	return tax, t.getStep(level, tax)
}

func (t *TaxService) calculateTaxTable(salary float64) ([]TaxStep, float64) {
	var totalTax float64
	var steps []TaxStep

	tax, step := t.calculateStep(salary, 0, 0, "0-150,000")
	totalTax += tax
	steps = append(steps, step)

	tax, step = t.calculateStep(salary, 150000, 0.1, "150,001-500,000")
	totalTax += tax
	steps = append(steps, step)

	tax, step = t.calculateStep(salary, 500000, 0.15, "500,001-1,000,000")
	totalTax += tax
	steps = append(steps, step)

	tax, step = t.calculateStep(salary, 1000000, 0.2, "1,000,001-2,000,000")
	totalTax += tax
	steps = append(steps, step)

	tax, step = t.calculateStep(salary, 2000000, 0.35, "2,000,001 ขึ้นไป")
	totalTax += tax
	steps = append(steps, step)

	return steps, totalTax
}