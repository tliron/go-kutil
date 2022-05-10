package util

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
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

func Joinq(s []string, sep string) string {
	var builder strings.Builder
	last := len(s) - 1
	for i, s_ := range s {
		builder.WriteString(strconv.Quote(s_))
		if i != last {
			builder.WriteString(sep)
		}
	}
	return builder.String()
}
