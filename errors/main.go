package errors

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
	"net/http"
	"runtime/debug"
)

type Error struct {
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields,omitempty"`
}

type Response struct {
	Error Error `json:"error"`
}

// Return error response
func SetErrorResponse(statusCode int, text string, c *gin.Context) {
	c.JSON(statusCode, Response{
		Error: Error{
			Message: text,
		},
	})
}

// Returns validation errors
func SetValidationErrorResponse(err error, c *gin.Context) {
	errs := err.(validator.ValidationErrors)

	errorFieldsList := make(map[string]string)
	for _, err := range errs {
		errorFieldsList[err.Field()] = validationErrorToText(err)
	}

	c.JSON(http.StatusUnprocessableEntity, Response{
		Error{
			Message: "Validation failed.",
			Fields:  errorFieldsList,
		},
	})
}

// Get validation error field message
func validationErrorToText(e validator.FieldError) string {
	field := e.Field()
	tag := e.Tag()

	switch tag {
	case "numeric":
		return fmt.Sprintf("The %s must be a number", field)
	case "max":
		return fmt.Sprintf("The length of field %s must be less than %s", field, e.Param())
	case "required":
		return fmt.Sprintf("The field %s is required", field)
	case "oneof":
		return fmt.Sprintf("The field %s can have the next values: %s", field, e.Param())
	default:
		return fmt.Sprintf("%s is not valid", field)
	}
}

func Recovery() {
	if err := recover(); err != nil {
		log.WithField("stacktrace", string(debug.Stack())).Error(err)
	}
}
