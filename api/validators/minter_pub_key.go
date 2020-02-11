package validators

import (
	"gopkg.in/go-playground/validator.v9"
	"regexp"
)

func MinterPublicKey(fl validator.FieldLevel) bool {
	return isValidMinterPublicKey(fl.Field().Interface().(string))
}

func isValidMinterPublicKey(publicKey string) bool {
	return regexp.MustCompile("^Mp([A-Fa-f0-9]{64})$").MatchString(publicKey)
}
