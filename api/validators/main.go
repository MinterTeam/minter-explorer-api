package validators

import (
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"gopkg.in/go-playground/validator.v8"
	"reflect"
	"regexp"
)

func MinterAddress(
    v *validator.Validate, topStruct reflect.Value, currentStructOrField reflect.Value,
    field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string,
) bool {
	if fieldType.String() == "[]string" {
		data, err := field.Interface().([]string)
		helpers.CheckErrBool(err)

		for _, address := range data {
			if validateMinterAddress(address) == false {
				return false
			}
		}

		return true
	}

	return validateMinterAddress(field.String())

}

func validateMinterAddress(address string) bool {
	return regexp.MustCompile("^Mx([A-Fa-f0-9]{40})$").MatchString(address)
}