package tax

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func TestPostCalculateTax(t *testing.T) {
	testCases := []struct {
		name   string
		param  CalculationRequest
		expect CalculationResult
	}{
		{
			name: "calculate tax correctly, given only total income",
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
			expect: CalculationResult{
				Tax: 29000,
			},
		},
		{
			name: "calculate tax correctly, given total income and withholding tax",
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
			expect: CalculationResult{
				Tax: 4000,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			e := echo.New()
			recorder := httptest.NewRecorder()

			taxCalculator := NewMockCalculator(ctrl)
			taxCalculator.EXPECT().Calculate(tc.param).Times(1).Return(tc.expect)

			taxController := NewTaxController(taxCalculator)
			taxController.RouteConfig(e)

			data, err := json.Marshal(tc.param)
			require.NoError(t, err)

			url := "/tax/calculations"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			request.Header.Set("Content-Type", "application/json")
			require.NoError(t, err)

			e.ServeHTTP(recorder, request)

			require.Equal(t, http.StatusOK, recorder.Code)

			response := recorder.Body
			responseBytes, err := io.ReadAll(response)
			require.NoError(t, err)

			var gotCalculationResult CalculationResult
			err = json.Unmarshal(responseBytes, &gotCalculationResult)
			require.NoError(t, err)

			require.Equal(t, tc.expect, gotCalculationResult)
		})
	}
}
