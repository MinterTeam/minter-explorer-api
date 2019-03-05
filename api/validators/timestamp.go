package validators

import (
	"gopkg.in/go-playground/validator.v8"
	"reflect"
	"time"
)

func Timestamp(
	v *validator.Validate, topStruct reflect.Value, currentStructOrField reflect.Value,
	field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string,
) bool {
	timestamp := field.String()
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
