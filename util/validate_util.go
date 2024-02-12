package util

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func Validator(data interface{}) (string, error) {
	validate := validator.New()
	if err := validate.Struct(data); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errorMessage := "Validation Error: "
		for i, e := range validationErrors {
			errorMessage += fmt.Sprintf("%d) Field: %s, Error: %s Value: %s\n", (i + 1), e.Field(), e.Tag(), e.Value())
		}
		return errorMessage, fmt.Errorf(errorMessage)
	}
	return "", nil
}
