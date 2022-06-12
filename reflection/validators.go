package reflection

//
// Validators
//

// ard.TypeValidator signature

// *string
func IsPointerToString(value any) bool {
	_, ok := value.(*string)
	return ok
}

// *int64
func IsPointerToInt64(value any) bool {
	_, ok := value.(*int64)
	return ok
}

// *float64
func IsPointerToFloat64(value any) bool {
	_, ok := value.(*float64)
	return ok
}

// *bool
func IsPointerToBool(value any) bool {
	_, ok := value.(*bool)
	return ok
}

// *[]string
func IsPointerToSliceOfString(value any) bool {
	_, ok := value.(*[]string)
	return ok
}

// *map[string]string
func IsPointerToMapOfStringToString(value any) bool {
	_, ok := value.(*map[string]string)
	return ok
}
