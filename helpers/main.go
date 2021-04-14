package helpers

import (
	"reflect"
	"strconv"
	"time"
)

type ChannelData struct {
	Value interface{}
	Error error
}

func NewChannelData(value interface{}, err error) ChannelData {
	return ChannelData{
		Value: value,
		Error: err,
	}
}

func InArray(needle interface{}, haystack interface{}) bool {
	switch reflect.TypeOf(haystack).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(haystack)
		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(needle, s.Index(i).Interface()) == true {
				return true
			}
		}
	}

	return false
}

func StartOfTheDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}

func StrToUint64(str string) uint64 {
	val, _ := strconv.ParseUint(str, 10, 64)
	return val
}
