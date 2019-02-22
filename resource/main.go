package resource

import (
	"reflect"
)

type ItemInterface interface{}

type Interface interface {
	Transform(model ItemInterface, params ...interface{}) Interface
}

func TransformCollection(collection interface{}, resource Interface) []Interface {
	val := reflect.ValueOf(collection)

	models := make([]ItemInterface, val.Len())
	for i := 0; i < val.Len(); i++ {
		if val.Index(i).Kind() == reflect.Ptr {
			models[i] = val.Index(i).Elem().Interface()
		} else {
			models[i] = val.Index(i).Interface()
		}
	}

	result := make([]Interface, len(models))
	for i := range models {
		result[i] = resource.Transform(models[i])
	}

	return result
}
