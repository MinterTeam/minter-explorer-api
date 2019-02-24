package validators

import (
	"gopkg.in/go-playground/validator.v8"
	"reflect"
	"regexp"
)

func MinterPublicKey(
	v *validator.Validate, topStruct reflect.Value, currentStructOrField reflect.Value,
	field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string,
) bool {
	return isValidMinterPublicKey(field.String())
}

func isValidMinterPublicKey(publicKey string) bool {
	return regexp.MustCompile("^Mp([A-Fa-f0-9]{64})$").MatchString(publicKey)
}
