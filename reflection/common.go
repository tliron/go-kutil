package reflection

import (
	"reflect"
	"runtime"
)

// See: https://stackoverflow.com/a/7053871/849021
func GetFunctionName(fn any) string {
	if function := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()); function != nil {
		return function.Name()
	} else {
		return "<unknown function>"
	}
}

func IsNil(value reflect.Value) bool {
	// https://golang.org/pkg/reflect/#Value.IsNil
	switch value.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Slice, reflect.Pointer:
		return value.IsNil()

	default:
		return false
	}
}

// See: https://stackoverflow.com/questions/23555241/golang-reflection-how-to-get-zero-value-of-a-field-type
func IsZero(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Slice:
		return value.IsNil()

	case reflect.Pointer:
		return value.IsNil() || IsZero(value.Elem())

	case reflect.Array:
		length := value.Len()
		for i := 0; i < length; i++ {
			if !IsZero(value.Index(i)) {
				return false
			}
		}
		return true

	case reflect.Struct:
		numField := value.NumField()
		for i := 0; i < numField; i++ {
			if !IsZero(value.Field(i)) {
				return false
			}
		}
		return true

	default:
		zero := reflect.Zero(value.Type()).Interface()
		return value.Interface() == zero
	}
}

func IsEmpty(value any) bool {
	// From JSON documentation:
	// defined as false, 0, a nil pointer, a nil interface value, and any empty array, slice, map, or string

	switch value_ := value.(type) {
	case bool:
		return value_ == false
	case int, int64, int32, int16, int8, uint, uint64, uint32, uint16, uint8, float64, float32:
		return value_ == 0
	case string:
		return len(value_) == 0
	}

	value_ := reflect.ValueOf(value)
	switch value_.Kind() {
	case reflect.Map, reflect.Slice, reflect.Array:
		if value_.Len() == 0 {
			return true
		}
	}

	return value == nil
}
