package validators

import (
	"github.com/go-playground/validator/v10"
	"regexp"
)

func MinterAddress(fl validator.FieldLevel) bool {
	if data, ok := fl.Field().Interface().([]string); ok {
		for _, address := range data {
			if !isValidMinterAddress(address) {
				return false
			}
		}

		return true
	}

	return isValidMinterAddress(fl.Field().Interface().(string))
}

func isValidMinterAddress(address string) bool {
	return regexp.MustCompile("^Mx([A-Fa-f0-9]{40})$").MatchString(address)
}
