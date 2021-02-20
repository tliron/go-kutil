package reflection

import (
	"reflect"
)

// We can't do type assertions here

// Compatible with *struct{}
func IsPtrToStruct(type_ reflect.Type) bool {
	return (type_.Kind() == reflect.Ptr) && (type_.Elem().Kind() == reflect.Struct)
}

// Compatible with []*struct{}
func IsSliceOfPtrToStruct(type_ reflect.Type) bool {
	return (type_.Kind() == reflect.Slice) && (type_.Elem().Kind() == reflect.Ptr) && (type_.Elem().Elem().Kind() == reflect.Struct)
}

// Compatible with map[string]*struct{}
func IsMapOfStringToPtrToStruct(type_ reflect.Type) bool {
	return (type_.Kind() == reflect.Map) && (type_.Key().Kind() == reflect.String) && (type_.Elem().Kind() == reflect.Ptr) && (type_.Elem().Elem().Kind() == reflect.Struct)
}
