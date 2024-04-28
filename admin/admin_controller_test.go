package admin

import (
	"bytes"
	"encoding/json"
	"io"
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
		checkResponse      func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Should response with 200 status code, given valid request",
			body: `
				{
					"amount": 100000.0
				}
			`,
			adminServiceStub: func(adminService *MockAdminService) {
				adminService.EXPECT().UpdatePersonalDeduction(gomock.Any(), 100000.0).Times(1).Return(100000.0, nil)
			},
			expectedStatusCode: http.StatusOK,
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				responseBody, err := io.ReadAll(recorder.Body)
				require.NoError(t, err)

				var actualResponse updatePersonalDeductionResponse
				err = json.Unmarshal(responseBody, &actualResponse)
				require.NoError(t, err)

				require.Equal(t, common.Float64(100000.0), actualResponse.PersonalDeduction)
			},
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
			checkResponse:      func(t *testing.T, recorder *httptest.ResponseRecorder) {},
		},
		{
			name: "Should response with 400 status code, given amount greater than 100000",
			body: `
				{
					"amount": 100000.1
				}
			`,
			adminServiceStub: func(adminService *MockAdminService) {
				adminService.EXPECT().UpdatePersonalDeduction(gomock.Any(), gomock.Any()).Times(0)
			},
			expectedStatusCode: http.StatusBadRequest,
			checkResponse:      func(t *testing.T, recorder *httptest.ResponseRecorder) {},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			appConfig := common.AppConfig{
				AdminUsername: "admin",
				AdminPassword: "P@ssw0rd",
			}
			adminService := NewMockAdminService(ctrl)
			adminController := NewAdminController(adminService, appConfig)

			tc.adminServiceStub(adminService)

			e := common.NewConfiguredEcho()

			adminController.RouteConfig(e)

			url := "/admin/deductions/personal"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader([]byte(tc.body)))
			require.NoError(t, err)

			request.SetBasicAuth("admin", "P@ssw0rd")
			request.Header.Set("Content-Type", "application/json")

			recorder := httptest.NewRecorder()
			e.ServeHTTP(recorder, request)

			require.Equal(t, tc.expectedStatusCode, recorder.Code)

			tc.checkResponse(t, recorder)
		})
	}
}

func TestPostUpdateKReceiptDeduction(t *testing.T) {
	testCases := []struct {
		name               string
		body               string
		adminServiceStub   func(adminService *MockAdminService)
		expectedStatusCode int
		checkResponse      func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Should response with 200 status code, given valid request",
			body: `
				{
					"amount": 100000.0
				}
			`,
			adminServiceStub: func(adminService *MockAdminService) {
				adminService.EXPECT().UpdateKReceiptDeduction(gomock.Any(), 100000.0).Times(1).Return(100000.0, nil)
			},
			expectedStatusCode: http.StatusOK,
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				responseBody, err := io.ReadAll(recorder.Body)
				require.NoError(t, err)

				var actualResponse updateKReceiptDeductionResponse
				err = json.Unmarshal(responseBody, &actualResponse)
				require.NoError(t, err)

				require.Equal(t, common.Float64(100000.0), actualResponse.KReceipt)
			},
		},
		{
			name: "Should response with 400 status code, given invalid request",
			body: `
				{
					"invalid": null
				}
			`,
			adminServiceStub: func(adminService *MockAdminService) {
				adminService.EXPECT().UpdateKReceiptDeduction(gomock.Any(), gomock.Any()).Times(0)
			},
			expectedStatusCode: http.StatusBadRequest,
			checkResponse:      func(t *testing.T, recorder *httptest.ResponseRecorder) {},
		},
		{
			name: "Should response with 400 status code, given k-receipt greater than 100000",
			body: `
				{
					"amount": 100000.1
				}
			`,
			adminServiceStub: func(adminService *MockAdminService) {
				adminService.EXPECT().UpdateKReceiptDeduction(gomock.Any(), gomock.Any()).Times(0)
			},
			expectedStatusCode: http.StatusBadRequest,
			checkResponse:      func(t *testing.T, recorder *httptest.ResponseRecorder) {},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			appConfig := common.AppConfig{
				AdminUsername: "admin",
				AdminPassword: "P@ssw0rd",
			}
			adminService := NewMockAdminService(ctrl)
			adminController := NewAdminController(adminService, appConfig)

			tc.adminServiceStub(adminService)

			e := common.NewConfiguredEcho()

			adminController.RouteConfig(e)

			url := "/admin/deductions/k-receipt"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader([]byte(tc.body)))
			require.NoError(t, err)

			request.SetBasicAuth("admin", "P@ssw0rd")
			request.Header.Set("Content-Type", "application/json")

			recorder := httptest.NewRecorder()
			e.ServeHTTP(recorder, request)

			require.Equal(t, tc.expectedStatusCode, recorder.Code)

			tc.checkResponse(t, recorder)
		})
	}
}
