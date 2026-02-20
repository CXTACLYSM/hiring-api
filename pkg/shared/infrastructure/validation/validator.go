package validation

import (
	"github.com/CXTACLYSM/hiring-api/pkg/shared/infrastructure/responses"
	"github.com/go-playground/validator/v10"
)

type Validator struct {
	validate *validator.Validate
}

func NewValidator(validator *validator.Validate) *Validator {
	return &Validator{
		validate: validator,
	}
}

func (v *Validator) Struct(s any) error {
	return v.validate.Struct(s)
}

func FormatErrors(errors validator.ValidationErrors) []responses.ValidationError {
	errorsFormatted := make([]responses.ValidationError, 0, len(errors))
	for _, fieldError := range errors {
		errorsFormatted = append(errorsFormatted, responses.ValidationError{
			Field:   fieldError.Field(),
			Message: message(fieldError),
		})
	}

	return errorsFormatted
}

func message(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "field is required"
	case "email":
		return "invalid email format"
	case "min":
		return "must be at least " + fe.Param() + " characters"
	case "max":
		return "must be at most " + fe.Param() + " characters"
	case "eqfield":
		return "must match " + fe.Param()
	default:
		return "invalid value"
	}
}
