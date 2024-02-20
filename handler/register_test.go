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

func TestRegister(t *testing.T) {
	tests := []struct {
		name                 string
		requestBody          string
		expectedStatus       int
		expectServiceCall    bool
		expectedServiceResp  *service.RegisterResponse
		expectedServiceError common.Error
	}{
		{
			name:                 "Success",
			requestBody:          `{"fullName": "maulana aji satrio", "phoneNumber": "+628232482440", "password": "Maulana1996@"}`,
			expectedStatus:       http.StatusOK,
			expectServiceCall:    true,
			expectedServiceResp:  &service.RegisterResponse{UserID: 123},
			expectedServiceError: nil,
		},
		{
			name:                 "BadRequest",
			requestBody:          "Invalid JSON",
			expectedStatus:       http.StatusBadRequest,
			expectServiceCall:    false,
			expectedServiceResp:  nil,
			expectedServiceError: nil,
		},
		{
			name:                 "ServiceError",
			requestBody:          `{"fullName": "maulana aji satrio", "phoneNumber": "+628232482440", "password": "Maulana1996@"}`,
			expectedStatus:       http.StatusBadRequest,
			expectServiceCall:    true,
			expectedServiceResp:  nil,
			expectedServiceError: commonErr.NewError("any", commonErr.BadRequestErrorType),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := mock_service.NewMockServiceInterface(ctrl)
			mockServer := &Server{
				Service: mockService,
			}

			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(tc.requestBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			if tc.expectServiceCall {
				mockService.EXPECT().Register(gomock.Any(), gomock.Any()).Return(tc.expectedServiceResp, tc.expectedServiceError)
			}

			// Execute
			mockServer.Register(c)

			// Assert
			assert.Equal(t, tc.expectedStatus, rec.Code)
		})
	}
}
