package tax

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestCalculateTax(t *testing.T) {
	taxLevels1 := createEmptyTaxLevels()
	taxLevels1[1].Tax = 29000

	taxLevels2 := createEmptyTaxLevels()
	taxLevels2[1].Tax = 29000

	taxLevels3 := createEmptyTaxLevels()
	taxLevels3[1].Tax = 19000

	taxLevels4 := createEmptyTaxLevels()
	taxLevels4[1].Tax = 20000

	testCases := []struct {
		name     string
		param    CalculationRequest
		expected CalculationResult
	}{
		{
			name: "Should calculate tax correctly, given only total income",
			param: CalculationRequest{
				TotalIncome: 500000,
				Wht:         0,
				Allowances: []Allowance{
					{
						AllowanceType: AllowanceDonation,
						Amount:        0,
					},
				},
			},
			expected: CalculationResult{
				Tax:       29000,
				TaxLevels: taxLevels1,
			},
		},
		{
			name: "Should calculate tax correctly, given total income and withholding tax",
			param: CalculationRequest{
				TotalIncome: 500000,
				Wht:         25000,
				Allowances: []Allowance{
					{
						AllowanceType: AllowanceDonation,
						Amount:        0,
					},
				},
			},
			expected: CalculationResult{
				Tax:       4000,
				TaxLevels: taxLevels2,
			},
		},
		{
			name: "Should calculate tax correctly, given total income and donation (over allowance limit of 100000)",
			param: CalculationRequest{
				TotalIncome: 500000,
				Wht:         0,
				Allowances: []Allowance{
					{
						AllowanceType: AllowanceDonation,
						Amount:        200000,
					},
				},
			},
			expected: CalculationResult{
				Tax:       19000,
				TaxLevels: taxLevels3,
			},
		},
		{
			name: "Should calculate tax correctly, given total income and donation (under allowance limit of 100000)",
			param: CalculationRequest{
				TotalIncome: 500000,
				Wht:         0,
				Allowances: []Allowance{
					{
						AllowanceType: AllowanceDonation,
						Amount:        90000,
					},
				},
			},
			expected: CalculationResult{
				Tax:       20000,
				TaxLevels: taxLevels4,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			taxConfigRepo := NewMockTaxConfigRepository(ctrl)

			ctx := context.Background()
			taxConfigRepo.EXPECT().FindByName(ctx, "personal_deduction").Times(1).Return(
				&Config{
					Name:  "personal_deduction",
					Value: 60000,
				},
				nil,
			)

			calculator := NewCalculator(taxConfigRepo)

			result := calculator.Calculate(ctx, tc.param)

			require.Equal(t, tc.expected.Tax, result.Tax)
			require.Equal(t, tc.expected.TaxLevels, result.TaxLevels)
		})
	}
}
