package reflection

import (
	"reflect"
)

// Compatible with *struct{}
func IsPointerToStruct(type_ reflect.Type) bool {
	return (type_.Kind() == reflect.Pointer) && (type_.Elem().Kind() == reflect.Struct)
}

// Compatible with []*struct{}
func IsSliceOfPointerToStruct(type_ reflect.Type) bool {
	return (type_.Kind() == reflect.Slice) && (type_.Elem().Kind() == reflect.Pointer) && (type_.Elem().Elem().Kind() == reflect.Struct)
}

// Compatible with map[string]*struct{}
func IsMapOfStringToPointerToStruct(type_ reflect.Type) bool {
	return (type_.Kind() == reflect.Map) && (type_.Key().Kind() == reflect.String) && (type_.Elem().Kind() == reflect.Pointer) && (type_.Elem().Elem().Kind() == reflect.Struct)
}

// int64, int32, int16, int8, int
func IsInteger(kind reflect.Kind) bool {
	return (kind == reflect.Int64) || (kind == reflect.Int32) || (kind == reflect.Int16) || (kind == reflect.Int8) || (kind == reflect.Int)
}

// uint64, uint32, uint16, uint8, uint
func IsUInteger(kind reflect.Kind) bool {
	return (kind == reflect.Uint64) || (kind == reflect.Uint32) || (kind == reflect.Uint16) || (kind == reflect.Uint8) || (kind == reflect.Uint)
}

// float64, float32
func IsFloat(kind reflect.Kind) bool {
	return (kind == reflect.Float64) || (kind == reflect.Float32)
}
