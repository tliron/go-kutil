package util

import (
	"reflect"
	"unsafe"
)

// See:
// https://go101.org/article/unsafe.html
// https://github.com/golang/go/issues/25484
// https://github.com/golang/go/issues/19367
// https://golang.org/src/strings/builder.go#L45

// This casting *does not* copy data. Note that casting via "string(bytes)" *does* copy data.
func BytesToString(bytes []byte) string {
	return *(*string)(unsafe.Pointer(&bytes))
}

func StringToBytes(string_ string) (bytes []byte) {
	stringHeader := (*reflect.StringHeader)(unsafe.Pointer(&string_))
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&bytes))
	sliceHeader.Data = stringHeader.Data
	sliceHeader.Cap = stringHeader.Len
	sliceHeader.Len = stringHeader.Len
	return
}
