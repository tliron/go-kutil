package reflection

import (
	"reflect"
	"runtime"
)

// See: https://stackoverflow.com/a/7053871/849021
func GetFunctionName(fn interface{}) string {
	runtimeFn := runtime.FuncForPC(reflect.ValueOf(fn).Pointer())
	if runtimeFn == nil {
		return "<unknown function>"
	}
	return runtimeFn.Name()
}

func IsNil(value reflect.Value) bool {
	// https://golang.org/pkg/reflect/#Value.IsNil
	switch value.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Slice, reflect.Ptr:
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

	case reflect.Ptr:
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
