package repository

import (
	"github.com/go-playground/validator/v10"
)

type ValidatorRepositoryInterface interface {
	validate(data interface{}) []ErrorResponse
	Validator(data interface{}) []ErrorResponse
}

type ErrorResponse struct {
	Error       bool        `json:"error"`
	FailedField string      `json:"failedField"`
	Tag         string      `json:"tag"`
	Value       interface{} `json:"value"`
}

type GlobalErrorHandlerResp struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type xValidator struct {
	validator *validator.Validate
}

func NewXValidator() *xValidator {
	return &xValidator{
		validator: validator.New(),
	}
}

func (v *xValidator) validate(data interface{}) []ErrorResponse {
	validationErrors := []ErrorResponse{}

	errs := v.validator.Struct(data)
	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			var elem ErrorResponse

			elem.FailedField = err.Field()
			elem.Tag = err.Tag()
			elem.Value = err.Value()
			elem.Error = true

			validationErrors = append(validationErrors, elem)
		}
	}

	return validationErrors
}

func (v *xValidator) Validator(data interface{}) []ErrorResponse {
	return v.validate(data)
}
