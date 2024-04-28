package tax

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chuckboliver/assessment-tax/common"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestPostCalculateTax(t *testing.T) {
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
		body     string
		expected CalculationResult
	}{
		{
			name: "calculate tax correctly, given only total income",
			body: `
				{
					"totalIncome": 500000,
					"wht": 0,
					"allowances": [
						{
							"allowanceType": "donation",
							"amount": 0
						}
					]
				}
			`,
			expected: CalculationResult{
				Tax:       29000,
				TaxLevels: taxLevels1,
			},
		},
		{
			name: "calculate tax correctly, given total income and withholding tax",
			body: `
				{
					"totalIncome": 500000,
					"wht": 25000,
					"allowances": [
						{
							"allowanceType": "donation",
							"amount": 0
						}
					]
				}
			`,
			expected: CalculationResult{
				Tax:       4000,
				TaxLevels: taxLevels2,
			},
		},
		{
			name: "Should calculate tax correctly, given total income and donation (over allowance limit of 100000)",
			body: `
				{
					"totalIncome": 500000,
					"wht": 0,
					"allowances": [
						{
							"allowanceType": "donation",
							"amount": 200000
						}
					]
				}
			`,
			expected: CalculationResult{
				Tax:       19000,
				TaxLevels: taxLevels3,
			},
		},
		{
			name: "Should calculate tax correctly, given total income and donation (under allowance limit of 100000)",
			body: `
				{
					"totalIncome": 500000,
					"wht": 0,
					"allowances": [
						{
							"allowanceType": "donation",
							"amount": 90000
						}
					]
				}
			`,
			expected: CalculationResult{
				Tax:       20000,
				TaxLevels: taxLevels4,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			taxCalculator := NewMockCalculator(ctrl)

			e := common.NewConfiguredEcho()
			taxController := NewTaxController(taxCalculator)
			taxController.RouteConfig(e)

			var expectedInputOfCalculate calculationRequest
			err := json.Unmarshal([]byte(tc.body), &expectedInputOfCalculate)
			require.NoError(t, err)

			ctx := context.Background()
			taxCalculator.EXPECT().Calculate(ctx, expectedInputOfCalculate).Times(1).Return(tc.expected)

			url := "/tax/calculations"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader([]byte(tc.body)))
			require.NoError(t, err)

			request.Header.Set("Content-Type", "application/json")

			recorder := httptest.NewRecorder()
			e.ServeHTTP(recorder, request)

			require.Equal(t, http.StatusOK, recorder.Code)

			response := recorder.Body
			responseBytes, err := io.ReadAll(response)
			require.NoError(t, err)

			var gotCalculationResult CalculationResult
			err = json.Unmarshal(responseBytes, &gotCalculationResult)
			require.NoError(t, err)

			require.Equal(t, tc.expected.Tax, gotCalculationResult.Tax)
			require.Equal(t, tc.expected.TaxLevels, gotCalculationResult.TaxLevels)
		})
	}
}
