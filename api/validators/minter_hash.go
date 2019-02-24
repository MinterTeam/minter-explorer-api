package validators

import (
	"gopkg.in/go-playground/validator.v8"
	"reflect"
	"regexp"
)

func MinterTxHash(
	v *validator.Validate, topStruct reflect.Value, currentStructOrField reflect.Value,
	field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string,
) bool {
	return isValidMinterHash(field.String())
}

func isValidMinterHash(hash string) bool {
	return regexp.MustCompile("^Mt([A-Fa-f0-9]{64})$").MatchString(hash)
}
