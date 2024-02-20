package common

type Error interface {
	GetErrorType() ErrorType
	GetErrorMessage() string
}
