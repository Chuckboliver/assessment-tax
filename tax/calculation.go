package tax

import (
	"context"
	"log/slog"

	"github.com/chuckboliver/assessment-tax/common"
)

type Allowance struct {
	AllowanceType AllowanceType `json:"allowanceType"`
	Amount        float64       `json:"amount"`
}

type AllowanceType string

const (
	AllowanceDonation AllowanceType = "donation"
	AllowanceKReceipt AllowanceType = "k-receipt"
)

const (
	defaultPersonalDeduction    = 60000.0
	defaultMaxKReceiptDeduction = 50000.0
)

type CalculationResultWithTaxLevel struct {
	Tax       common.Float64 `json:"tax"`
	TaxRefund common.Float64 `json:"taxRefund"`
	TaxLevels []TaxLevel     `json:"taxLevel"`
}

type BatchCalculationResult struct {
	Taxes []CalculationResult `json:"taxes"`
}

type CalculationResult struct {
	TotalIncome common.Float64 `json:"totalIncome"`
	Tax         common.Float64 `json:"tax"`
	TaxRefund   common.Float64 `json:"taxRefund"`
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

func (c *CalculatorImpl) calculate(personalDeduction float64, maxKReceiptDeduction float64, param calculationRequest) CalculationResultWithTaxLevel {
	income := param.TotalIncome - personalDeduction

	income = c.applyAllowances(income, param.Allowances, maxKReceiptDeduction)

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
	taxRefund := 0.0
	if tax < 0 {
		taxRefund = -tax
		tax = 0
	}

	return CalculationResultWithTaxLevel{
		Tax:       common.Float64(tax),
		TaxRefund: common.Float64(taxRefund),
		TaxLevels: taxLevels,
	}
}

func (c *CalculatorImpl) Calculate(ctx context.Context, param calculationRequest) CalculationResultWithTaxLevel {
	personalDeduction := c.getPersonalDeduction(ctx)
	maxKReceiptDeduction := c.getMaxKReceiptDeduction(ctx)
	return c.calculate(personalDeduction, maxKReceiptDeduction, param)
}

func (c *CalculatorImpl) BatchCalculate(ctx context.Context, params []calculationRequest) BatchCalculationResult {
	personalDeduction := c.getPersonalDeduction(ctx)
	maxKReceiptDeduction := c.getMaxKReceiptDeduction(ctx)

	calculationResults := make([]CalculationResult, 0, len(params))
	for _, v := range params {
		calculationResultWithTaxLevel := c.calculate(personalDeduction, maxKReceiptDeduction, v)

		calculationResult := CalculationResult{
			TotalIncome: common.Float64(v.TotalIncome),
			Tax:         calculationResultWithTaxLevel.Tax,
			TaxRefund:   calculationResultWithTaxLevel.TaxRefund,
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
		slog.Error("failed to get personal deduction", err)
		return defaultPersonalDeduction
	}

	return config.Value
}

func (c *CalculatorImpl) getMaxKReceiptDeduction(ctx context.Context) float64 {
	config, err := c.taxConfigRepository.FindByName(ctx, "kreceipt_deduction")
	if err != nil {
		slog.Error("failed to get max kreceipt deduction", err)
		return defaultMaxKReceiptDeduction
	}

	return config.Value
}

func (c *CalculatorImpl) applyAllowances(income float64, allowances []Allowance, maxKReceiptDeduction float64) float64 {

	for _, v := range allowances {
		switch v.AllowanceType {
		case AllowanceDonation:
			income -= min(v.Amount, 100000)
		case AllowanceKReceipt:
			income -= min(v.Amount, maxKReceiptDeduction)
		}
	}

	return income
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
