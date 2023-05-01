package validation

import (
	"github.com/go-playground/validator/v10"
	"github.com/iwpnd/detectr/errors"
)

// ValidateStruct ...
func ValidateStruct(s interface{}) []*errors.ErrRequestError {
	var errs []*errors.ErrRequestError
	validate := validator.New()
	err := validate.Struct(s)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element errors.ErrRequestError
			element.Source = err.Field()
			element.Title = "Invalid Attribute"
			element.Detail = err.Tag()
			element.Status = 400
			errs = append(errs, &element)
		}
	}
	return errs
}
