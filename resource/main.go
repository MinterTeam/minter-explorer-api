package resource

import (
	"reflect"
)

type ItemInterface interface{}

type ResourceItemInterface interface {
	Transform(model ItemInterface) ResourceItemInterface
}

func TransformCollection(collection interface{}, resource ResourceItemInterface) []ResourceItemInterface {
	val := reflect.ValueOf(collection)

	models := make([]ItemInterface, val.Len())
	for i := 0; i < val.Len(); i++ {
		models[i] = val.Index(i).Interface()
	}

	result := make([]ResourceItemInterface, len(models))
	for i := range models {
		result[i] = resource.Transform(models[i])
	}

	return result
}
