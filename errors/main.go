package errors

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/validator.v8"
	"net/http"
)

type Error struct {
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields,omitempty"`
}

type Response struct {
	Error Error `json:"error"`
}

// Return error response
func SetErrorResponse(statusCode int, errorCode int, text string, c *gin.Context) {
	c.JSON(statusCode, Response{
		Error: Error{
			Code:    errorCode,
			Message: text,
		},
	})
}

// Returns validation errors
func SetValidationErrorResponse(err error, c *gin.Context) {
	errs := err.(validator.ValidationErrors)

	errorFieldsList := make(map[string]string)
	for _, err := range errs {
		errorFieldsList[err.Field] = validationErrorToText(err)
	}

	c.JSON(http.StatusUnprocessableEntity, Response{
		Error{
			Code:    1,
			Message: "Validation failed.",
			Fields:  errorFieldsList,
		},
	})
}

// Get validation error field message
func validationErrorToText(e *validator.FieldError) string {
	field := e.Field
	tag := e.Tag

	switch tag {
	case "numeric":
		return fmt.Sprintf("The %s must be a number", field)
	case "max":
		return fmt.Sprintf("The length of field %s must be less than %s", field, e.Param)
	case "required":
		return fmt.Sprintf("The field %s is required", field)
	default:
		return fmt.Sprintf("%s is not valid", field)
	}
}
