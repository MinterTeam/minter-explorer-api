package helpers

import (
	"reflect"
)

const DefaultStatisticsScale = "day"  // TODO: move to a better place, may be config?
const DefaultStatisticsDayDelta = -14 // TODO: move to a better place, may be config?

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
