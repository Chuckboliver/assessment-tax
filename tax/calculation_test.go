package tax

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCalculateTax(t *testing.T) {
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
						AllowanceType: "donation",
						Amount:        0,
					},
				},
			},
			expected: CalculationResult{
				Tax: 29000,
			},
		},
		{
			name: "Should calculate tax correctly, given total income and withholding tax",
			param: CalculationRequest{
				TotalIncome: 500000,
				Wht:         25000,
				Allowances: []Allowance{
					{
						AllowanceType: "donation",
						Amount:        0,
					},
				},
			},
			expected: CalculationResult{
				Tax: 4000,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			calculator := NewCalculator()
			result := calculator.Calculate(tc.param)

			require.Equal(t, tc.expected.Tax, result.Tax)
		})
	}
}
