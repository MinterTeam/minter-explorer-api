package validators

import (
	"gopkg.in/go-playground/validator.v8"
	"reflect"
	"regexp"
)

func MinterAddress(
	v *validator.Validate, topStruct reflect.Value, currentStructOrField reflect.Value,
	field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string,
) bool {
	if fieldType.String() == "[]string" {
		data, _ := field.Interface().([]string)
		for _, address := range data {
			if !isValidMinterAddress(address) {
				return false
			}
		}

		return true
	}

	return isValidMinterAddress(field.String())
}

func MinterTxHash(
	v *validator.Validate, topStruct reflect.Value, currentStructOrField reflect.Value,
	field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string,
) bool {
	return isValidMinterHash(field.String())
}

func MinterPublicKey(
	v *validator.Validate, topStruct reflect.Value, currentStructOrField reflect.Value,
	field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string,
) bool {
	return isValidMinterPublicKey(field.String())
}

func isValidMinterAddress(address string) bool {
	return regexp.MustCompile("^Mx([A-Fa-f0-9]{40})$").MatchString(address)
}

func isValidMinterHash(hash string) bool {
	return regexp.MustCompile("^Mt([A-Fa-f0-9]{64})$").MatchString(hash)
}

func isValidMinterPublicKey(publicKey string) bool {
	return regexp.MustCompile("^Mp([A-Fa-f0-9]{64})$").MatchString(publicKey)
}
