package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/sawitpro/UserService/common"
	commonErr "github.com/sawitpro/UserService/common/errors"
	mock_service "github.com/sawitpro/UserService/mocks"
	"github.com/sawitpro/UserService/service"
)

func TestUpdateProfile(t *testing.T) {
	tests := []struct {
		name                 string
		userIDCtxValue       interface{}
		requestBody          string
		expectedStatus       int
		expectServiceCall    bool
		expectedServiceResp  *service.UpdateProfileResponse
		expectedServiceError common.Error
	}{
		{
			name:              "Success",
			userIDCtxValue:    int64(1),
			requestBody:       `{"fullName": "maulana aji satrio", "phoneNumber": "+628232482440"}`,
			expectedStatus:    http.StatusOK,
			expectServiceCall: true,
			expectedServiceResp: &service.UpdateProfileResponse{
				FullName:    "maulana aji satrio",
				PhoneNumber: "+628232482440",
			},
			expectedServiceError: nil,
		},
		{
			name:                 "ForbiddenAccess",
			userIDCtxValue:       "invalid",
			requestBody:          `{"fullName": "maulana aji satrio", "phoneNumber": "+628232482440"}`,
			expectedStatus:       http.StatusForbidden,
			expectServiceCall:    false,
			expectedServiceResp:  nil,
			expectedServiceError: nil,
		},
		{
			name:                 "BadRequest JSON",
			userIDCtxValue:       int64(1),
			requestBody:          `Invalid JSON`,
			expectedStatus:       http.StatusBadRequest,
			expectServiceCall:    false,
			expectedServiceResp:  nil,
			expectedServiceError: nil,
		},
		{
			name:                 "BadRequest All Params Nil",
			userIDCtxValue:       int64(1),
			requestBody:          `{}`,
			expectedStatus:       http.StatusBadRequest,
			expectServiceCall:    false,
			expectedServiceResp:  nil,
			expectedServiceError: nil,
		},
		{
			name:                 "ServiceError Conflict",
			userIDCtxValue:       int64(1),
			requestBody:          `{"fullName": "maulana aji satrio", "phoneNumber": "+628232482440"}`,
			expectedStatus:       http.StatusConflict,
			expectServiceCall:    true,
			expectedServiceResp:  nil,
			expectedServiceError: commonErr.NewError("any", commonErr.ConflictErrorType),
		},
		{
			name:                 "ServiceError Bad Request",
			userIDCtxValue:       int64(1),
			requestBody:          `{"fullName": "maulana aji satrio", "phoneNumber": "+628232482440"}`,
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
			req := httptest.NewRequest(http.MethodPatch, "/profile", strings.NewReader(tc.requestBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.Set(common.USER_ID_CTX_KEY, tc.userIDCtxValue)

			if tc.expectServiceCall {
				mockService.EXPECT().UpdateProfile(gomock.Any(), gomock.Any()).Return(tc.expectedServiceResp, tc.expectedServiceError)
			}

			mockServer.UpdateProfile(c)
			assert.Equal(t, tc.expectedStatus, rec.Code)
		})
	}
}
