package util

import (
	"fmt"
	"reflect"
	"strconv"
	stringspkg "strings"
	"unsafe"
)

// See:
// https://go101.org/article/unsafe.html
// https://github.com/golang/go/issues/25484
// https://github.com/golang/go/issues/19367
// https://golang.org/src/strings/builder.go#L45

// This casting *does not* copy data. Note that casting via "string(value)" *does* copy data.
func BytesToString(bytes []byte) string {
	return *(*string)(unsafe.Pointer(&bytes))
}

// This casting *does not* copy data. Note that casting via "[]byte(value)" *does* copy data.
func StringToBytes(string_ string) (bytes []byte) {
	stringHeader := (*reflect.StringHeader)(unsafe.Pointer(&string_))
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&bytes))
	sliceHeader.Data = stringHeader.Data
	sliceHeader.Cap = stringHeader.Len
	sliceHeader.Len = stringHeader.Len
	return
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
	case error:
		return value_.Error()
	default:
		return fmt.Sprintf("%v", value_)
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
