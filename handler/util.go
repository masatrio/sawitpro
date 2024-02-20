package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/sawitpro/UserService/common"
	cmnErr "github.com/sawitpro/UserService/common/errors"
	"github.com/sawitpro/UserService/generated"
)

func getErrorHttpStatusCode(err common.Error) int {
	code := http.StatusInternalServerError

	switch err.GetErrorType() {
	case cmnErr.BadRequestErrorType:
		code = http.StatusBadRequest
	case cmnErr.ConflictErrorType:
		code = http.StatusConflict
	case cmnErr.SystemErrorType:
		code = http.StatusInternalServerError
	}

	return code
}

func handleServiceError(ctx echo.Context, err common.Error) error {
	return ctx.JSON(
		getErrorHttpStatusCode(err),
		&generated.ErrorResponse{ErrorMsg: err.GetErrorMessage()},
	)
}

func handleBadRequestJSON(ctx echo.Context, err error) error {
	return ctx.JSON(http.StatusBadRequest, &generated.ErrorResponse{ErrorMsg: err.Error()})
}

func handleForbiddenAccessJSON(ctx echo.Context, err error) error {
	return ctx.JSON(http.StatusForbidden, &generated.ErrorResponse{ErrorMsg: err.Error()})
}

func handleSuccessJSON(ctx echo.Context, body interface{}) error {
	return ctx.JSON(http.StatusOK, body)
}

func getContextUserID(ctx echo.Context) (int64, error) {
	userID, ok := ctx.Get(common.USER_ID_CTX_KEY).(int64)

	if !ok {
		return 0, errors.New("error user ID Format")
	}

	return userID, nil
}

func getBearerToken(ctx echo.Context) string {
	const authPrefix = "Bearer "
	authHeader := ctx.Request().Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, authPrefix) {
		return ""
	}
	return strings.TrimPrefix(authHeader, authPrefix)
}
