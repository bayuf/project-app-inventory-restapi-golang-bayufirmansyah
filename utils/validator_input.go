package utils

import (
	"fmt"
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/shopspring/decimal"
)

var (
	validate *validator.Validate
	once     sync.Once
)

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func getValidator() *validator.Validate {
	once.Do(func() {
		validate = validator.New()

		// REGISTER CUSTOM DECIMAL VALIDATOR
		validate.RegisterValidation("decimal_gt_zero", func(fl validator.FieldLevel) bool {
			d, ok := fl.Field().Interface().(decimal.Decimal)
			if !ok {
				return false
			}
			return d.GreaterThan(decimal.Zero)
		})
	})

	return validate
}
func ValidateInput(data any) (string, error) {
	validate := getValidator()

	err := validate.Struct(data)
	if err == nil {
		return "", nil
	}

	if _, ok := err.(*validator.InvalidValidationError); ok {
		return "", err
	}

	var errors []string
	for _, e := range err.(validator.ValidationErrors) {
		var message string
		switch e.Tag() {
		case "email":
			message = "Please input correct email format"
		case "decimal_gt_zero":
			message = fmt.Sprintf("%s must be greater than 0", e.Field())
		default:
			message = fmt.Sprintf("%s must %s", e.Field(), e.Tag())
		}
		errors = append(errors, message)
	}

	return fmt.Sprint(errors), err
}

func ValidateErrors(data any) ([]FieldError, error) {
	validate := getValidator()

	err := validate.Struct(data)
	if err == nil {
		return nil, nil
	}

	var errors []FieldError

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, err := range validationErrors {
			var message string

			switch err.Tag() {
			case "required":
				message = fmt.Sprintf("%s is required", err.Field())
			case "email":
				message = "Please enter a valid email format"
			case "decimal_gt_zero":
				message = fmt.Sprintf("%s must be greater than 0", err.Field())
			case "gte":
				message = fmt.Sprintf("%s must be a non-negative number", err.Field())
			case "min":
				message = fmt.Sprintf("%s must be at least %s characters long", err.Field(), err.Param())
			case "eqfield":
				message = fmt.Sprintf("%s must match %s", err.Field(), err.Param())
			default:
				message = fmt.Sprintf("%s is invalid", err.Field())
			}

			errors = append(errors, FieldError{
				Field:   err.Field(),
				Message: message,
			})
		}
		return errors, err
	}

	return nil, err
}
