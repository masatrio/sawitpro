package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/sawitpro/UserService/common"
	commonErr "github.com/sawitpro/UserService/common/errors"
	mock_service "github.com/sawitpro/UserService/mocks"
	"github.com/sawitpro/UserService/service"
)

func TestGetProfile(t *testing.T) {
	tests := []struct {
		name                 string
		userIDCtxValue       interface{}
		expectedStatus       int
		expectServiceCall    bool
		expectedServiceResp  *service.UserInfoResponse
		expectedServiceError common.Error
	}{
		{
			name:                 "Success",
			userIDCtxValue:       int64(1),
			expectedStatus:       http.StatusOK,
			expectServiceCall:    true,
			expectedServiceResp:  &service.UserInfoResponse{FullName: "maulana aji satrio", PhoneNumber: "+628232482440"},
			expectedServiceError: nil,
		},
		{
			name:                 "ForbiddenAccess",
			userIDCtxValue:       "invalid",
			expectedStatus:       http.StatusForbidden,
			expectServiceCall:    false,
			expectedServiceResp:  nil,
			expectedServiceError: nil,
		},
		{
			name:                 "ServiceInternalError",
			userIDCtxValue:       int64(1),
			expectedStatus:       http.StatusInternalServerError,
			expectServiceCall:    true,
			expectedServiceResp:  nil,
			expectedServiceError: commonErr.NewError("any", commonErr.SystemErrorType),
		},
		{
			name:                 "ServiceBadReqError",
			userIDCtxValue:       int64(1),
			expectedStatus:       http.StatusBadRequest,
			expectServiceCall:    true,
			expectedServiceResp:  nil,
			expectedServiceError: commonErr.NewError("any", commonErr.BadRequestErrorType),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := mock_service.NewMockServiceInterface(ctrl)
			mockServer := &Server{
				Service: mockService,
			}

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/profile", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.Set(common.USER_ID_CTX_KEY, tc.userIDCtxValue)

			if tc.expectServiceCall {
				mockService.EXPECT().GetProfile(gomock.Any(), gomock.Any()).Return(tc.expectedServiceResp, tc.expectedServiceError)
			}

			mockServer.GetProfile(c)
			assert.Equal(t, tc.expectedStatus, rec.Code)
		})
	}
}
