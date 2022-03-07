package ard

import (
	"fmt"
	"time"
)

//
// TypeName
//

type TypeName string

const (
	NoType TypeName = ""

	// Failsafe schema: https://yaml.org/spec/1.2/spec.html#id2802346
	TypeMap    TypeName = "ard.map"
	TypeList   TypeName = "ard.list"
	TypeString TypeName = "ard.string"

	// JSON schema: https://yaml.org/spec/1.2/spec.html#id2803231
	TypeBoolean TypeName = "ard.boolean"
	TypeInteger TypeName = "ard.integer"
	TypeFloat   TypeName = "ard.float"

	// Other schemas: https://yaml.org/spec/1.2/spec.html#id2805770
	TypeNull      TypeName = "ard.null"
	TypeTimestamp TypeName = "ard.timestamp"
)

func GetTypeName(value Value) TypeName {
	switch value.(type) {
	case Map:
		return TypeMap
	case List:
		return TypeList
	case string:
		return TypeString
	case bool:
		return TypeBoolean
	case int64, int32, int16, int8, int, uint64, uint32, uint16, uint8, uint:
		return TypeInteger
	case float64, float32:
		return TypeFloat
	case nil:
		return TypeNull
	case time.Time:
		return TypeTimestamp
	default:
		return TypeName(fmt.Sprintf("%T", value))
	}
}

//
// TypeZeroes
//

var TypeZeroes = map[TypeName]Value{
	TypeMap:       make(Map), // YAML parser returns Map, but JSON parser returns StringMap
	TypeList:      List{},
	TypeString:    "",
	TypeBoolean:   false,
	TypeInteger:   int(0),       // YAML parser returns int
	TypeFloat:     float64(0.0), // YAML and JSON parsers return float64
	TypeNull:      nil,
	TypeTimestamp: time.Time{}, // YAML parser returns time.Time
}

//
// TypeValidator
//

type TypeValidator = func(Value) bool

var TypeValidators = map[TypeName]TypeValidator{
	TypeMap:       IsMap,
	TypeList:      IsList,
	TypeString:    IsString,
	TypeBoolean:   IsBoolean,
	TypeInteger:   IsInteger,
	TypeFloat:     IsFloat,
	TypeNull:      IsNull,
	TypeTimestamp: IsTimestamp,
}

// Map = map[any]any
func IsMap(value Value) bool {
	_, ok := value.(Map)
	return ok
}

// List = []any
func IsList(value Value) bool {
	_, ok := value.(List)
	return ok
}

// string
func IsString(value Value) bool {
	_, ok := value.(string)
	return ok
}

// bool
func IsBoolean(value Value) bool {
	_, ok := value.(bool)
	return ok
}

// int64, int32, int16, int8, int, uint64, uint32, uint16, uint8, uint
func IsInteger(value Value) bool {
	switch value.(type) {
	case int64, int32, int16, int8, int, uint64, uint32, uint16, uint8, uint:
		return true
	}
	return false
}

// float64, float32
func IsFloat(value Value) bool {
	switch value.(type) {
	case float64, float32:
		return true
	}
	return false
}

func IsNull(value Value) bool {
	return value == nil
}

// time.Time
func IsTimestamp(value Value) bool {
	_, ok := value.(time.Time)
	return ok
}
