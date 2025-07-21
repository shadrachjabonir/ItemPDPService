package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// ValidationError represents a validation error response
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Value   interface{} `json:"value,omitempty"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error  string            `json:"error"`
	Errors []ValidationError `json:"errors,omitempty"`
}

// validator instance
var validate *validator.Validate

func init() {
	validate = validator.New()
}

// ValidateJSON validates the JSON body against a struct
func ValidateJSON(obj interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := c.ShouldBindJSON(obj); err != nil {
			var validationErrors []ValidationError

			if validationErrs, ok := err.(validator.ValidationErrors); ok {
				for _, fieldErr := range validationErrs {
					validationErrors = append(validationErrors, ValidationError{
						Field:   getFieldName(fieldErr.Field()),
						Message: getValidationMessage(fieldErr),
						Value:   fieldErr.Value(),
					})
				}
			}

			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:  "Validation failed",
				Errors: validationErrors,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// ValidateStruct validates a struct using the validator
func ValidateStruct(s interface{}) []ValidationError {
	var validationErrors []ValidationError

	if err := validate.Struct(s); err != nil {
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			for _, fieldErr := range validationErrs {
				validationErrors = append(validationErrors, ValidationError{
					Field:   getFieldName(fieldErr.Field()),
					Message: getValidationMessage(fieldErr),
					Value:   fieldErr.Value(),
				})
			}
		}
	}

	return validationErrors
}

// ValidateAndRespond validates a struct and responds with errors if any
func ValidateAndRespond(c *gin.Context, s interface{}) bool {
	if errors := ValidateStruct(s); len(errors) > 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:  "Validation failed",
			Errors: errors,
		})
		return false
	}
	return true
}

// getFieldName converts field name to snake_case
func getFieldName(field string) string {
	var result strings.Builder
	for i, r := range field {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}

// getValidationMessage returns a human-readable validation message
func getValidationMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Must be a valid email address"
	case "min":
		return "Must be at least " + fe.Param() + " characters long"
	case "max":
		return "Must be at most " + fe.Param() + " characters long"
	case "len":
		return "Must be exactly " + fe.Param() + " characters long"
	case "url":
		return "Must be a valid URL"
	case "gte":
		return "Must be greater than or equal to " + fe.Param()
	case "lte":
		return "Must be less than or equal to " + fe.Param()
	case "gt":
		return "Must be greater than " + fe.Param()
	case "lt":
		return "Must be less than " + fe.Param()
	case "oneof":
		return "Must be one of: " + fe.Param()
	default:
		return "Invalid value"
	}
} 