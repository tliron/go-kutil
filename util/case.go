package util

import (
	"unicode"
)

func ToCamelCase(name string) string {
	// TODO: should results be cached?
	runes := []rune(name)
	length := len(runes)

	if (length > 0) && unicode.IsUpper(runes[0]) { // sanity check
		if (length > 1) && unicode.IsUpper(runes[1]) {
			// If the second rune is also uppercase we'll keep the name as is
			return name
		}
		runes_ := make([]rune, 1, length-1)
		runes_[0] = unicode.ToLower(runes[0])
		return string(append(runes_, runes[1:]...))
	} else {
		return name
	}
}
