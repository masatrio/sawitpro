package errors

import (
	"fmt"

	"github.com/sawitpro/UserService/common"
)

type Error struct {
	Message string
	Type    common.ErrorType
}

const (
	SystemErrorType     common.ErrorType = "SystemErrorType"
	BadRequestErrorType common.ErrorType = "BadRequestErrorType"
	ConflictErrorType   common.ErrorType = "ConflictedErrorType"
)

const (
	WrongPhonePasswordErrorMessage string = "password or phone number is incorrect."
	phoneAlreadyUsedErrorMessage   string = "phone number %s already used."
	UserDataNotFoundErrorMessage   string = "user data not found."
)

func NewError(message string, errorType common.ErrorType) common.Error {
	return &Error{
		Message: message,
		Type:    errorType,
	}
}

func (e *Error) GetErrorType() common.ErrorType {
	return e.Type
}

func (e *Error) GetErrorMessage() string {
	return e.Message
}

func NewPhoneAlreadyUsedErrorMessage(phoneNumber string) string {
	return fmt.Sprintf(phoneAlreadyUsedErrorMessage, phoneNumber)
}
