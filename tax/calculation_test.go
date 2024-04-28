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

	taxLevels5 := createEmptyTaxLevels()
	taxLevels5[1].Tax = 22500

	testCases := []struct {
		name              string
		arg               calculationRequest
		taxConfigRepoStub func(taxConfigRepo *MockTaxConfigRepository)
		expected          CalculationResult
	}{
		{
			name: "Should calculate tax correctly, given only total income",
			arg: calculationRequest{
				TotalIncome: 500000,
				Wht:         0,
				Allowances: []Allowance{
					{
						AllowanceType: AllowanceDonation,
						Amount:        0,
					},
				},
			},
			taxConfigRepoStub: func(taxConfigRepo *MockTaxConfigRepository) {
				taxConfigRepo.EXPECT().
					FindByName(gomock.Any(), "personal_deduction").
					Times(1).
					Return(&Config{
						Name:  "personal_deduction",
						Value: 60000.0,
					}, nil)
			},
			expected: CalculationResult{
				Tax:       29000,
				TaxLevels: taxLevels1,
			},
		},
		{
			name: "Should calculate tax correctly, given total income and withholding tax",
			arg: calculationRequest{
				TotalIncome: 500000,
				Wht:         25000,
				Allowances: []Allowance{
					{
						AllowanceType: AllowanceDonation,
						Amount:        0,
					},
				},
			},
			taxConfigRepoStub: func(taxConfigRepo *MockTaxConfigRepository) {
				taxConfigRepo.EXPECT().
					FindByName(gomock.Any(), "personal_deduction").
					Times(1).
					Return(&Config{
						Name:  "personal_deduction",
						Value: 60000.0,
					}, nil)
			},
			expected: CalculationResult{
				Tax:       4000,
				TaxLevels: taxLevels2,
			},
		},
		{
			name: "Should calculate tax correctly, given total income and donation (over allowance limit of 100000)",
			arg: calculationRequest{
				TotalIncome: 500000,
				Wht:         0,
				Allowances: []Allowance{
					{
						AllowanceType: AllowanceDonation,
						Amount:        200000,
					},
				},
			},
			taxConfigRepoStub: func(taxConfigRepo *MockTaxConfigRepository) {
				taxConfigRepo.EXPECT().
					FindByName(gomock.Any(), "personal_deduction").
					Times(1).
					Return(&Config{
						Name:  "personal_deduction",
						Value: 60000.0,
					}, nil)
			},
			expected: CalculationResult{
				Tax:       19000,
				TaxLevels: taxLevels3,
			},
		},
		{
			name: "Should calculate tax correctly, given total income and donation (under allowance limit of 100000)",
			arg: calculationRequest{
				TotalIncome: 500000,
				Wht:         0,
				Allowances: []Allowance{
					{
						AllowanceType: AllowanceDonation,
						Amount:        90000,
					},
				},
			},
			taxConfigRepoStub: func(taxConfigRepo *MockTaxConfigRepository) {
				taxConfigRepo.EXPECT().
					FindByName(gomock.Any(), "personal_deduction").
					Times(1).
					Return(&Config{
						Name:  "personal_deduction",
						Value: 60000.0,
					}, nil)
			},
			expected: CalculationResult{
				Tax:       20000,
				TaxLevels: taxLevels4,
			},
		},
		{
			name: "Should calculate tax correctly, when personal deduction is configured",
			arg: calculationRequest{
				TotalIncome: 500000,
				Wht:         0,
				Allowances: []Allowance{
					{
						AllowanceType: AllowanceDonation,
						Amount:        90000,
					},
				},
			},
			taxConfigRepoStub: func(taxConfigRepo *MockTaxConfigRepository) {
				taxConfigRepo.EXPECT().
					FindByName(gomock.Any(), "personal_deduction").
					Times(1).
					Return(&Config{
						Name:  "personal_deduction",
						Value: 35000.0,
					}, nil)
			},
			expected: CalculationResult{
				Tax:       22500,
				TaxLevels: taxLevels5,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			taxConfigRepo := NewMockTaxConfigRepository(ctrl)
			calculator := NewCalculator(taxConfigRepo)

			tc.taxConfigRepoStub(taxConfigRepo)

			ctx := context.Background()
			result := calculator.Calculate(ctx, tc.arg)

			require.Equal(t, tc.expected.Tax, result.Tax)
			require.Equal(t, tc.expected.TaxLevels, result.TaxLevels)
		})
	}
}
