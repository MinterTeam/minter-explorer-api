package validators

import (
	"gopkg.in/go-playground/validator.v9"
	"time"
)

func Timestamp(fl validator.FieldLevel) bool {
	timestamp := fl.Field().Interface().(string)
	_, err := time.Parse("2006-01-02", timestamp)
	if err == nil {
		return true
	}

	_, err = time.Parse("2006-01-02 15:04:05", timestamp)
	if err == nil {
		return true
	}

	_, err = time.Parse(time.RFC3339, timestamp)
	if err == nil {
		return true
	}

	return false
}
