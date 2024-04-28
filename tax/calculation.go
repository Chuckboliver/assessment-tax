package tax

import (
	"context"

	"github.com/chuckboliver/assessment-tax/common"
)

type Allowance struct {
	AllowanceType AllowanceType `json:"allowanceType"`
	Amount        float64       `json:"amount"`
}

type AllowanceType string

const (
	AllowanceDonation AllowanceType = "donation"
)

const (
	defaultPersonalDeduction = 60000.0
)

type CalculationResultWithTaxLevel struct {
	Tax       common.Float64 `json:"tax"`
	TaxLevels []TaxLevel     `json:"taxLevel"`
}

type BatchCalculationResult struct {
	Taxes []CalculationResult `json:"taxes"`
}

type CalculationResult struct {
	TotalIncome common.Float64 `json:"totalIncome"`
	Tax         common.Float64 `json:"tax"`
}

type TaxLevel struct {
	Level string         `json:"level"`
	Tax   common.Float64 `json:"tax"`
}

type TaxConfigRepository interface {
	FindByName(ctx context.Context, name string) (*Config, error)
}

type Calculator interface {
	Calculate(ctx context.Context, param calculationRequest) CalculationResultWithTaxLevel
	BatchCalculate(ctx context.Context, params []calculationRequest) BatchCalculationResult
}

var _ Calculator = (*CalculatorImpl)(nil)

type CalculatorImpl struct {
	taxConfigRepository TaxConfigRepository
}

func NewCalculator(taxConfigRepository TaxConfigRepository) Calculator {
	return &CalculatorImpl{
		taxConfigRepository: taxConfigRepository,
	}
}

func (c *CalculatorImpl) calculate(personalDeduction float64, param calculationRequest) CalculationResultWithTaxLevel {
	income := param.TotalIncome - personalDeduction

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

	return CalculationResultWithTaxLevel{
		Tax:       common.Float64(tax),
		TaxLevels: taxLevels,
	}
}

func (c *CalculatorImpl) Calculate(ctx context.Context, param calculationRequest) CalculationResultWithTaxLevel {
	personalDeduction := c.getPersonalDeduction(ctx)
	return c.calculate(personalDeduction, param)
}

func (c *CalculatorImpl) BatchCalculate(ctx context.Context, params []calculationRequest) BatchCalculationResult {
	personalDeduction := c.getPersonalDeduction(ctx)

	calculationResults := make([]CalculationResult, 0, len(params))
	for _, v := range params {
		calculationResultWithTaxLevel := c.calculate(personalDeduction, v)

		calculationResult := CalculationResult{
			TotalIncome: common.Float64(v.TotalIncome),
			Tax:         calculationResultWithTaxLevel.Tax,
		}

		calculationResults = append(calculationResults, calculationResult)
	}

	return BatchCalculationResult{
		Taxes: calculationResults,
	}
}

func (c *CalculatorImpl) getPersonalDeduction(ctx context.Context) float64 {
	config, err := c.taxConfigRepository.FindByName(ctx, "personal_deduction")
	if err != nil {
		return defaultPersonalDeduction
	}

	return config.Value
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
