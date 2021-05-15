package validators

import (
	"github.com/go-playground/validator/v10"
	"regexp"
)

func MinterTxHash(fl validator.FieldLevel) bool {
	return isValidMinterHash(fl.Field().Interface().(string))
}

func isValidMinterHash(hash string) bool {
	return regexp.MustCompile("^Mt([A-Fa-f0-9]{64})$").MatchString(hash)
}
