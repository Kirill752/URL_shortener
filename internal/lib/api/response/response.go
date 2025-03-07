package response

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

const (
	StatusOk    = "OK"
	StatusError = "Error"
)

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

func OK() Response {
	return Response{
		Status: StatusOk,
	}
}

func Error(msg string) Response {
	return Response{
		Status: StatusError,
		Error:  msg,
	}
}

// Обработка ошибки валидатора
func ValidationError(errs validator.ValidationErrors) Response {
	var errMsg []string
	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsg = append(errMsg, fmt.Sprintf("field %s is required field", err.Field()))
		case "url":
			errMsg = append(errMsg, fmt.Sprintf("field %s is not valid URL", err.Field()))
		default:
			errMsg = append(errMsg, fmt.Sprintf("field %s is not valid", err.Field()))
		}
	}
	return Response{
		Status: StatusError,
		Error:  strings.Join(errMsg, " "),
	}
}
