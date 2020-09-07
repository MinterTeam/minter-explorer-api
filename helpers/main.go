package helpers

import (
	"reflect"
	"strconv"
	"strings"
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

func GetSymbolAndVersionFromStr(symbol string) (string, uint64) {
	items := strings.Split(symbol, "-")
	symbol, version := items[0], uint64(0)

	if len(items) == 2 {
		version, _ = strconv.ParseUint(items[1], 10, 64)
	}

	return symbol, version
}
