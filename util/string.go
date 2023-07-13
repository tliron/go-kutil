package util

import (
	"fmt"
	"strconv"
	stringspkg "strings"
	"unsafe"
)

// See:
// https://go101.org/article/unsafe.html
// https://github.com/golang/go/issues/25484
// https://github.com/golang/go/issues/19367
// https://golang.org/src/strings/builder.go

// This conversion *does not* copy data. Note that converting via "(string)([]byte)" *does* copy data.
// Also note that you *should not* change the byte slice after conversion, because Go strings
// are treated as immutable. This would cause a segmentation violation panic.
func BytesToString(bytes []byte) string {
	return unsafe.String(unsafe.SliceData(bytes), len(bytes))

	// return *(*string)(unsafe.Pointer(&bytes))
}

// This conversion *does not* copy data. Note that converting via "([]byte)(string)" *does* copy data.
// Also note that you *should not* change the byte slice after conversion, because Go strings
// are treated as immutable. This would cause a segmentation violation panic.
func StringToBytes(string_ string) (bytes []byte) {
	return unsafe.Slice(unsafe.StringData(string_), len(string_))

	/*
		// StringHeader and SliceHeader have been deprecated in Go 1.21
		stringHeader := (*reflect.StringHeader)(unsafe.Pointer(&string_))
		sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&bytes))
		sliceHeader.Data = stringHeader.Data
		sliceHeader.Cap = stringHeader.Len
		sliceHeader.Len = stringHeader.Len
		return
	*/
}

func ToString(value any) string {
	if value == nil {
		return "nil"
	}
	switch value_ := value.(type) {
	case string:
		return value_
	case fmt.Stringer:
		return value_.String()
	case []byte:
		return BytesToString(value_)
	case error:
		return value_.Error()
	default:
		return fmt.Sprintf("%+v", value_)
	}
}

func ToStrings(values []any) []string {
	length := len(values)
	if length == 0 {
		return nil
	}
	strings := make([]string, length)
	for index, value := range values {
		strings[index] = ToString(value)
	}
	return strings
}

func JoinQuote(strings []string, separator string) string {
	var builder stringspkg.Builder
	ultimateIndex := len(strings) - 1
	for index, value := range strings {
		builder.WriteString(strconv.Quote(value))
		if index != ultimateIndex {
			builder.WriteString(separator)
		}
	}
	return builder.String()
}

func JoinQuoteL(strings []string, separator string, lastSeparator string, coupleSeparator string) string {
	var builder stringspkg.Builder
	if len(strings) == 2 {
		builder.WriteString(strconv.Quote(strings[0]))
		builder.WriteString(coupleSeparator)
		builder.WriteString(strconv.Quote(strings[1]))
	} else {
		ultimateIndex := len(strings) - 1
		penultimateIndex := ultimateIndex - 1
		for index, value := range strings {
			builder.WriteString(strconv.Quote(value))
			if index != ultimateIndex {
				if index == penultimateIndex {
					builder.WriteString(lastSeparator)
				} else {
					builder.WriteString(separator)
				}
			}
		}
	}
	return builder.String()
}
