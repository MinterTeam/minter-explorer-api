package helpers

import (
	"log"
)

func CheckErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func CheckErrBool(ok bool) {
	if !ok {
		log.Panic(ok)
	}
}
