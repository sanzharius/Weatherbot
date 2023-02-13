package apperrors

import "fmt"

type AppError struct {
	Message string
	Code    string
}

var (
	MessageUnmarshallingError = AppError{
		Message: "Couldn't unmarshal a response",
		Code:    "UNMARSHAL_ERR",
	}
	ConfigReadErr = AppError{
		Message: "couldn't read config",
		Code:    "CONFIG_READ_ER",
	}
	DataNotFoundErr = AppError{
		Message: "Cannot get a weather forecast",
		Code:    "DATA_NOT_FOUND_ERR",
	}
	APICallingErr = AppError{
		Code: "API_CALLING_ERR",
	}
)

func (appError *AppError) Error() string {
	return appError.Code + ": " + appError.Message
}
func (appError *AppError) AppendMessage(anyErrs ...interface{}) *AppError {
	return &AppError{
		Message: fmt.Sprintf("%v: %v", appError.Message, anyErrs),
		Code:    appError.Code,
	}
}
