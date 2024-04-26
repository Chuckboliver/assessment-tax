package tax

import (
	"github.com/chuckboliver/assessment-tax/common"
)

type CalculationRequest struct {
	TotalIncome float64     `json:"totalIncome"`
	Wht         float64     `json:"wht"`
	Allowances  []Allowance `json:"allowances"`
}

type Allowance struct {
	AllowanceType string  `json:"allowanceType"`
	Amount        float64 `json:"amount"`
}

type CalculationResult struct {
	Tax common.Float64 `json:"tax"`
}

type Calculator interface {
	Calculate(param CalculationRequest) CalculationResult
}

type CalculatorImpl struct{}

func NewCalculator() CalculatorImpl {
	return CalculatorImpl{}
}

func (c *CalculatorImpl) Calculate(param CalculationRequest) CalculationResult {
	income := param.TotalIncome - 60000

	tax := 0.0

	if income > 2000000 {
		tax += (income - 2000000) * 0.35
	}

	if income > 1000000 {
		tax += (income - 1000000) * 0.2
	}

	if income > 500000 {
		tax += (income - 500000) * 0.15
	}

	if income > 150000 {
		tax += (income - 150000) * 0.1
	}

	return CalculationResult{
		Tax: common.Float64(tax),
	}
}
