package admin

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func TestPostUpdatePersonalDeduction(t *testing.T) {
	testCases := []struct {
		name               string
		arg                personalDeductionRequest
		expectedStatusCode int
	}{
		{
			name: "Should response with 200 status code, given valid request",
			arg: personalDeductionRequest{
				Amount: 20000.0,
			},
			expectedStatusCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			adminService := NewMockAdminService(ctrl)
			adminController := NewAdminController(adminService)

			adminService.EXPECT().UpdatePersonalDeduction(gomock.Any(), 20000.0).Times(1).Return(nil)

			e := echo.New()
			adminController.RouteConfig(e)

			jsonData, err := json.Marshal(tc.arg)
			require.NoError(t, err)

			url := "/admin/deductions/personal"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(jsonData))
			require.NoError(t, err)

			request.Header.Set("Content-Type", "application/json")

			recorder := httptest.NewRecorder()
			e.ServeHTTP(recorder, request)

			require.Equal(t, tc.expectedStatusCode, recorder.Code)
		})
	}
}
