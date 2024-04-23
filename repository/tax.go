package taxrepo

type TaxBracket struct {
	Level int
	LowerBound float64
	UpperBound float64
	TaxRate    float64
}