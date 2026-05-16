package corevalidator

import "github.com/go-playground/validator/v10"

var validate = validator.New()

type FieldError struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
}

func ValidateStruct(s interface{}) ([]FieldError, error) {
	err := validate.Struct(s)
	if err == nil {
		return nil, nil
	}

	validationErrors := err.(validator.ValidationErrors)

	var errs []FieldError
	for _, e := range validationErrors {
		errs = append(errs, FieldError{
			Field: e.Field(),
			Tag:   e.Tag(),
		})
	}

	return errs, err
}
