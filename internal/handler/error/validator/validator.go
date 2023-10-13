package validator

import (
	"GatewayService/internal/handler/response"
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

// ProcessValidatorError converts error generated during request data validation in response.JSONResult
func ProcessValidatorError(errs error) response.JSONResult {
	res := make(map[string]string)
	var e validator.ValidationErrors
	ok := errors.As(errs, &e)

	if !ok {
		return response.CreateJSONResult("Error", "Invalid argument passed through request as param or part of param")
	}

	for _, err := range e {
		res[err.Field()] = "not " + err.Tag()
	}

	return response.CreateJSONResult("Error", res)
}

type errorMsgWithExtraData struct {
	Description string   `json:"description"`
	Extra       []string `json:"extra,omitempty"`
}

// ErrorMsg returns generated from given string error response.JSONResult
func ErrorMsg(field string, extra ...string) response.JSONResult {
	msg := errorMsgWithExtraData{
		Description: fmt.Sprintf("Invalid argument passed through request as param or part of param(%s)", field),
		Extra:       extra,
	}
	return response.CreateJSONResult("Error", msg)
}
