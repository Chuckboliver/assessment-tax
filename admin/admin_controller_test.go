package admin

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chuckboliver/assessment-tax/common"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestPostUpdatePersonalDeduction(t *testing.T) {
	testCases := []struct {
		name               string
		body               string
		adminServiceStub   func(adminService *MockAdminService)
		expectedStatusCode int
	}{
		{
			name: "Should response with 200 status code, given valid request",
			body: `
				{
					"amount": 20000.0
				}
			`,
			adminServiceStub: func(adminService *MockAdminService) {
				adminService.EXPECT().UpdatePersonalDeduction(gomock.Any(), 20000.0).Times(1).Return(nil)
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "Should response with 400 status code, given invalid request",
			body: `
				{
					"invalid": null
				}
			`,
			adminServiceStub: func(adminService *MockAdminService) {
				adminService.EXPECT().UpdatePersonalDeduction(gomock.Any(), gomock.Any()).Times(0)
			},
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			adminService := NewMockAdminService(ctrl)
			adminController := NewAdminController(adminService)

			tc.adminServiceStub(adminService)

			e := common.NewConfiguredEcho()

			adminController.RouteConfig(e)

			url := "/admin/deductions/personal"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader([]byte(tc.body)))
			require.NoError(t, err)

			request.Header.Set("Content-Type", "application/json")

			recorder := httptest.NewRecorder()
			e.ServeHTTP(recorder, request)

			require.Equal(t, tc.expectedStatusCode, recorder.Code)
		})
	}
}
