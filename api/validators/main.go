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

func isValidMinterAddress(address string) bool {
	return regexp.MustCompile("^Mx([A-Fa-f0-9]{40})$").MatchString(address)
}
