package validation

import (
	"GatewayService/internal/handler/response"
	"errors"
	"github.com/go-playground/validator/v10"
	"regexp"
	"time"
)

func ValidateOwnerName(fl validator.FieldLevel) bool {
	ownerName := fl.Field().String()
	regexPattern := `^[A-Za-z\s]+,\s?[A-Za-z\s]+$`
	match, _ := regexp.MatchString(regexPattern, ownerName)
	return match
}

func ValidateAddress(fl validator.FieldLevel) bool {
	address := fl.Field().String()
	regexPattern := `^[A-Za-z\s]+,\s?[A-Za-z\s]+,\s?[A-Za-z0-9\s]+$`
	match, _ := regexp.MatchString(regexPattern, address)
	return match
}

func ValidateTimeFormat(fl validator.FieldLevel) bool {
	timeStr := fl.Field().String()
	regexPattern := `^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}$`
	match, _ := regexp.MatchString(regexPattern, timeStr)
	if !match {
		return false
	}

	_, err := time.Parse("2006-01-02 15:04:05", timeStr)
	return err == nil
}

func RegisterCustomValidators(validate *validator.Validate) error {
	err := validate.RegisterValidation("ownerNameFormat", ValidateOwnerName)
	if err != nil {
		return err
	}
	err = validate.RegisterValidation("addressFormat", ValidateAddress)
	if err != nil {
		return err
	}
	err = validate.RegisterValidation("timeFormat", ValidateTimeFormat)
	if err != nil {
		return err
	}
	return nil
}

// FormatValidatorError function builds error response message if validation fails
func FormatValidatorError(errs error) response.JSONResult {
	res := make(map[string]string)
	var e validator.ValidationErrors
	ok := errors.As(errs, &e)

	if !ok {
		return response.BuildJSONResponse("Error", "Invalid argument passed")
	}

	for _, err := range e {
		res[err.Field()] = err.Tag()
	}

	return response.BuildJSONResponse("Error", res)
}
