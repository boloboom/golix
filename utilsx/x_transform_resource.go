package utilsx

import (
	"log"
	"reflect"
)

type AbstractResourceInterface interface {
	transform(model interface{}) interface{}
}

type ResourceTransformer struct {
	Resource AbstractResourceInterface
}

type ResourceTransformerInterface interface {
	Make(interface{}) interface{}
	Collection(interface{}) []interface{}
}

// NewResourceTransformer creates a new ResourceTransformerInterface.
//
// It takes a parameter `resource` of type `AbstractResourceInterface` and
// returns a `ResourceTransformerInterface`.
func NewResourceTransformer(resource AbstractResourceInterface) ResourceTransformerInterface {
	return &ResourceTransformer{
		Resource: resource,
	}
}

// Make calls the transform method of the Resource field of the ResourceTransformer struct.
//
// Parameters:
//   - model: The model to be transformed.
//
// Returns:
//   - interface: The transformed model.
func (trans *ResourceTransformer) Make(model interface{}) interface{} {
	return trans.Resource.transform(model)
}

// Collection transforms a slice of models into a slice of interfaces.
//
// Parameters:
//   - models: The slice of models to be transformed.
//
// Returns:
//   - interfaceSlice: The slice of transformed models.
func (trans *ResourceTransformer) Collection(models interface{}) []interface{} {
	modelSlice, ok := trans.anySliceToInterfaceSlice(models)
	if !ok {
		log.Fatal("models parameter is not a slice")
	}
	for index := range modelSlice {
		modelSlice[index] = trans.Resource.transform(modelSlice[index])
	}
	return modelSlice
}

// anySliceToInterfaceSlice converts a slice of any type to a slice of interface{} type.
//
// It takes in a parameter `anySlice` of type `interface{}` which represents the slice
// of any type to be converted.
//
// The function returns two values: `interfaceSlice` of type `[]interface{}` which
// represents the converted slice, and `ok` of type `bool` which indicates whether the
// conversion was successful or not.
func (trans *ResourceTransformer) anySliceToInterfaceSlice(anySlice interface{}) (interfaceSlice []interface{}, ok bool) {
	slice, ok := trans.reflectSliceValue(anySlice)
	if !ok {
		ok = false
		return
	}
	sliceLength := slice.Len()
	interfaceSlice = make([]interface{}, sliceLength)
	for i := 0; i < sliceLength; i++ {
		interfaceSlice[i] = slice.Index(i).Interface()
	}
	ok = true
	return
}

// reflectSliceValue returns the reflect.Value of the input slice and a boolean
// indicating if the input is a slice.
//
// Parameters:
//   - i: the input interface{} to be reflected.
//
// Returns:
//   - val: the reflect.Value of the input slice.
//   - ok: a boolean indicating if the input is a slice.
func (trans *ResourceTransformer) reflectSliceValue(i interface{}) (val reflect.Value, ok bool) {
	if val = reflect.ValueOf(i); val.Kind() == reflect.Slice {
		ok = true
	}
	return
}
