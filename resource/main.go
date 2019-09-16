package resource

import (
	"reflect"
)

type ItemInterface interface{}
type ParamInterface interface{}
type ParamsInterface []ParamInterface
type Interface interface {
	Transform(model ItemInterface, params ...ParamInterface) Interface
}

func TransformCollection(collection interface{}, resource Interface) []Interface {
	models := makeItemsFromModelsCollection(collection)
	result := make([]Interface, len(models))
	for i := range models {
		result[i] = resource.Transform(models[i])
	}

	return result
}

func TransformCollectionWithCallback(collection interface{}, resource Interface, callbackFunc func(model ParamInterface) ParamsInterface) []Interface {
	models := makeItemsFromModelsCollection(collection)
	result := make([]Interface, len(models))
	for i := range models {
		result[i] = resource.Transform(models[i], callbackFunc(models[i])...)
	}

	return result
}

func makeItemsFromModelsCollection(collection interface{}) []ItemInterface {
	val := reflect.ValueOf(collection)

	models := make([]ItemInterface, val.Len())
	for i := 0; i < val.Len(); i++ {
		if val.Index(i).Kind() == reflect.Ptr {
			models[i] = val.Index(i).Elem().Interface()
		} else {
			models[i] = val.Index(i).Interface()
		}
	}

	return models
}
