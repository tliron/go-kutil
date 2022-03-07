package reflection

//
// Validators
//

// ard.TypeValidator signature

// *string
func IsPtrToString(value any) bool {
	_, ok := value.(*string)
	return ok
}

// *int64
func IsPtrToInt64(value any) bool {
	_, ok := value.(*int64)
	return ok
}

// *float64
func IsPtrToFloat64(value any) bool {
	_, ok := value.(*float64)
	return ok
}

// *bool
func IsPtrToBool(value any) bool {
	_, ok := value.(*bool)
	return ok
}

// *[]string
func IsPtrToSliceOfString(value any) bool {
	_, ok := value.(*[]string)
	return ok
}

// *map[string]string
func IsPtrToMapOfStringToString(value any) bool {
	_, ok := value.(*map[string]string)
	return ok
}
