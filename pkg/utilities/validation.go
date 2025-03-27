package utilities

import "github.com/go-playground/validator/v10"

// Helper to format validation errors
func FormatValidationErrors(validationErrors validator.ValidationErrors) map[string]string {
	errors := make(map[string]string)
	for _, fieldError := range validationErrors {
		errors[fieldError.Field()] = getValidationErrorMessage(fieldError)
	}
	return errors
}

// Generate human-readable error messages
func getValidationErrorMessage(fieldError validator.FieldError) string {
	switch fieldError.Tag() {
	case "required":
		return "This field is required."
	case "email":
		return "Invalid email address."
	case "min":
		return "Value is too short."
	case "max":
		return "Value is too long."
	case "len":
		return "Value must have an exact length."
	default:
		return "Invalid value."
	}
}
