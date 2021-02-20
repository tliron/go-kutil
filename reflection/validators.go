package reflection

//
// Validators
//

// ard.TypeValidator signature

// *string
func IsPtrToString(value interface{}) bool {
	_, ok := value.(*string)
	return ok
}

// *int64
func IsPtrToInt64(value interface{}) bool {
	_, ok := value.(*int64)
	return ok
}

// *float64
func IsPtrToFloat64(value interface{}) bool {
	_, ok := value.(*float64)
	return ok
}

// *bool
func IsPtrToBool(value interface{}) bool {
	_, ok := value.(*bool)
	return ok
}

// *[]string
func IsPtrToSliceOfString(value interface{}) bool {
	_, ok := value.(*[]string)
	return ok
}

// *map[string]string
func IsPtrToMapOfStringToString(value interface{}) bool {
	_, ok := value.(*map[string]string)
	return ok
}
