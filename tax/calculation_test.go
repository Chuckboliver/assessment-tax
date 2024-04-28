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

	taxLevels6 := createEmptyTaxLevels()
	taxLevels6[1].Tax = 20000

	taxLevels7 := createEmptyTaxLevels()
	taxLevels7[1].Tax = 17000

	taxLevels8 := createEmptyTaxLevels()
	taxLevels8[1].Tax = 15000

	testCases := []struct {
		name              string
		arg               calculationRequest
		taxConfigRepoStub func(taxConfigRepo *MockTaxConfigRepository)
		expected          CalculationResultWithTaxLevel
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
			expected: CalculationResultWithTaxLevel{
				Tax:       29000,
				TaxRefund: 0,
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
			expected: CalculationResultWithTaxLevel{
				Tax:       4000,
				TaxRefund: 0,
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
			expected: CalculationResultWithTaxLevel{
				Tax:       19000,
				TaxRefund: 0,
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
			expected: CalculationResultWithTaxLevel{
				Tax:       20000,
				TaxRefund: 0,
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
			expected: CalculationResultWithTaxLevel{
				Tax:       22500,
				TaxRefund: 0,
				TaxLevels: taxLevels5,
			},
		},
		{
			name: "Should calculate tax refund correctly, when withholding tax is more than calculated tax",
			arg: calculationRequest{
				TotalIncome: 500000,
				Wht:         30000,
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
			expected: CalculationResultWithTaxLevel{
				Tax:       0,
				TaxRefund: 10000,
				TaxLevels: taxLevels6,
			},
		},
		{
			name: "Should calculate tax correctly, given allowance type of k-receipt",
			arg: calculationRequest{
				TotalIncome: 500000,
				Wht:         21000,
				Allowances: []Allowance{
					{
						AllowanceType: AllowanceDonation,
						Amount:        90000,
					},
					{
						AllowanceType: AllowanceKReceipt,
						Amount:        30000,
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
			expected: CalculationResultWithTaxLevel{
				Tax:       0,
				TaxRefund: 4000,
				TaxLevels: taxLevels7,
			},
		},
		{
			name: "Should calculate tax correctly, given allowance type of k-receipt (over allowance limit of 50000)",
			arg: calculationRequest{
				TotalIncome: 500000,
				Wht:         15000,
				Allowances: []Allowance{
					{
						AllowanceType: AllowanceDonation,
						Amount:        90000,
					},
					{
						AllowanceType: AllowanceKReceipt,
						Amount:        60000,
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
			expected: CalculationResultWithTaxLevel{
				Tax:       0,
				TaxRefund: 0,
				TaxLevels: taxLevels8,
			},
		},
		{
			name: "Should calculate tax correctly, given allowance type of k-receipt (equal to allowance limit of 50000)",
			arg: calculationRequest{
				TotalIncome: 500000,
				Wht:         15000,
				Allowances: []Allowance{
					{
						AllowanceType: AllowanceDonation,
						Amount:        90000,
					},
					{
						AllowanceType: AllowanceKReceipt,
						Amount:        50000,
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
			expected: CalculationResultWithTaxLevel{
				Tax:       0,
				TaxRefund: 0,
				TaxLevels: taxLevels8,
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
			require.Equal(t, tc.expected.TaxRefund, result.TaxRefund)
		})
	}
}

func TestBatchCalculate(t *testing.T) {
	testCases := []struct {
		name              string
		arg               []calculationRequest
		taxConfigRepoStub func(taxConfigRepo *MockTaxConfigRepository)
		expected          BatchCalculationResult
	}{
		{
			name: "Should calculate tax for batch input correctly",
			arg: []calculationRequest{
				{
					TotalIncome: 500000,
					Wht:         0,
					Allowances: []Allowance{
						{
							AllowanceType: AllowanceDonation,
							Amount:        0,
						},
					},
				},
				{
					TotalIncome: 600000,
					Wht:         55000,
					Allowances: []Allowance{
						{
							AllowanceType: AllowanceDonation,
							Amount:        20000,
						},
					},
				},
				{
					TotalIncome: 750000,
					Wht:         50000,
					Allowances: []Allowance{
						{
							AllowanceType: AllowanceDonation,
							Amount:        15000,
						},
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
			expected: BatchCalculationResult{
				Taxes: []CalculationResult{
					{
						TotalIncome: 500000,
						Tax:         29000,
						TaxRefund:   0,
					},
					{
						TotalIncome: 600000,
						Tax:         0,
						TaxRefund:   15000,
					},
					{
						TotalIncome: 750000,
						Tax:         28750,
						TaxRefund:   0,
					},
				},
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
			result := calculator.BatchCalculate(ctx, tc.arg)

			require.Equal(t, len(tc.expected.Taxes), len(result.Taxes))

			for i := 0; i < len(tc.expected.Taxes); i++ {
				require.Equal(t, tc.expected.Taxes[i].TotalIncome, result.Taxes[i].TotalIncome)
				require.Equal(t, tc.expected.Taxes[i].Tax, result.Taxes[i].Tax)
				require.Equal(t, tc.expected.Taxes[i].TaxRefund, result.Taxes[i].TaxRefund)
			}
		})
	}
}
