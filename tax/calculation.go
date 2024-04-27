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
	AllowanceType AllowanceType `json:"allowanceType"`
	Amount        float64       `json:"amount"`
}

type AllowanceType string

const (
	AllowanceDonation AllowanceType = "donation"
)

type CalculationResult struct {
	Tax       common.Float64 `json:"tax"`
	TaxLevels []TaxLevel     `json:"taxLevel"`
}

type TaxLevel struct {
	Level string         `json:"level"`
	Tax   common.Float64 `json:"tax"`
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

	for _, v := range param.Allowances {
		switch v.AllowanceType {
		case AllowanceDonation:
			income -= min(v.Amount, 100000)
		}
	}

	taxLevels := createEmptyTaxLevels()

	tax := 0.0

	if income > 2000000 {
		currentLevelTax := (income - 2000000) * 0.35
		taxLevels[4].Tax = common.Float64(currentLevelTax)
		tax += currentLevelTax
	}

	if income > 1000000 {
		currentLevelTax := (income - 1000000) * 0.2
		taxLevels[3].Tax = common.Float64(currentLevelTax)
		tax += currentLevelTax
	}

	if income > 500000 {
		currentLevelTax := (income - 500000) * 0.15
		taxLevels[2].Tax = common.Float64(currentLevelTax)
		tax += currentLevelTax
	}

	if income > 150000 {
		currentLevelTax := (income - 150000) * 0.1
		taxLevels[1].Tax = common.Float64(currentLevelTax)
		tax += currentLevelTax
	}

	tax -= param.Wht

	return CalculationResult{
		Tax:       common.Float64(tax),
		TaxLevels: taxLevels,
	}
}

func createEmptyTaxLevels() []TaxLevel {
	return []TaxLevel{
		{
			Level: "0-150,000",
			Tax:   0,
		},
		{
			Level: "150,001-500,000",
			Tax:   0,
		},
		{
			Level: "500,001-1,000,000",
			Tax:   0,
		},
		{
			Level: "1,000,001-2,000,000",
			Tax:   0,
		},
		{
			Level: "2,000,001 ขึ้นไป",
			Tax:   0,
		},
	}
}
