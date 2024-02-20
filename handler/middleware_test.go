package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sawitpro/UserService/helper"
	"github.com/stretchr/testify/assert"
)

func TestInitMiddleware(t *testing.T) {

	middlewareFuncs, err := InitMiddleware()
	assert.NoError(t, err)
	assert.NotNil(t, middlewareFuncs)
	assert.Len(t, middlewareFuncs, 3)
}

func TestAuthMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		requestPath    string
		isValidToken   bool
		expectedUserID int64
		expectError    bool
		expectedStatus int
		httpMethod     string
	}{
		{
			name:           "ValidToken",
			requestPath:    "/profile",
			isValidToken:   true,
			expectedUserID: 1,
			expectError:    false,
			expectedStatus: http.StatusOK,
			httpMethod:     http.MethodPost,
		},
		{
			name:           "InvalidToken",
			requestPath:    "/profile",
			isValidToken:   false,
			expectedUserID: 0,
			expectError:    true,
			expectedStatus: http.StatusForbidden,
			httpMethod:     http.MethodGet,
		},
		{
			name:           "wrong prefix",
			requestPath:    "/profile",
			isValidToken:   false,
			expectedUserID: 0,
			expectError:    true,
			expectedStatus: http.StatusForbidden,
			httpMethod:     http.MethodGet,
		},
		{
			name:           "WhitelistedPath 1",
			requestPath:    "/auth/login",
			isValidToken:   false,
			expectedUserID: 0,
			expectError:    false,
			expectedStatus: http.StatusOK,
			httpMethod:     http.MethodPost,
		},
		{
			name:           "WhitelistedPath 1",
			requestPath:    "/auth/register",
			isValidToken:   false,
			expectedUserID: 0,
			expectError:    false,
			expectedStatus: http.StatusOK,
			httpMethod:     http.MethodPost,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(tc.httpMethod, tc.requestPath, nil)
			if tc.isValidToken {
				token, _ := helper.CreateToken(tc.expectedUserID, time.Duration(1)*time.Hour)
				req.Header.Set("Authorization", "Bearer "+token)
			} else {
				req.Header.Set("Authorization", "Bear 123")
			}

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := AuthMiddleware(func(c echo.Context) error {
				userID, ok := c.Get("userID").(int64)
				if ok {
					assert.Equal(t, tc.expectedUserID, userID)
				}

				return nil
			})(c)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			if !tc.expectError {
				assert.Equal(t, tc.expectedStatus, rec.Code)
			}
		})
	}
}
