package ard

import (
	"github.com/tliron/yamlkeys"
)

// Ensure data adheres to the ARD map type
// (JSON decoding uses map[string]interface{} instead of map[interface{}]interface{})
func Normalize(value Value) (Value, bool) {
	switch value_ := value.(type) {
	case StringMap:
		return StringMapToMap(value_), true

	case Map:
		map_ := make(Map)
		changed := false
		for key, element := range value_ {
			if value, changed_ := Normalize(element); changed_ {
				map_[key] = value
				changed = true
			}
		}
		if changed {
			return map_, true
		}

	case List:
		list := make(List, len(value_))
		changed := false
		for index, element := range value_ {
			if value, changed_ := Normalize(element); changed_ {
				list[index] = value
				changed = true
			}
		}
		if changed {
			return list, true
		}
	}

	return value, false
}

func MapsToStringMaps(value Value) (Value, bool) {
	switch value_ := value.(type) {
	case Map:
		return MapToStringMap(value_), true

	case StringMap:
		stringMap := make(StringMap)
		changed := false
		for key, element := range value_ {
			if value, changed_ := MapsToStringMaps(element); changed_ {
				stringMap[key] = value
				changed = true
			}
		}
		if changed {
			return stringMap, true
		}

	case List:
		list := make(List, len(value_))
		changed := false
		for index, element := range value_ {
			if value, changed_ := MapsToStringMaps(element); changed_ {
				list[index] = value
				changed = true
			}
		}
		if changed {
			return list, true
		}
	}

	return value, false
}

// Ensure data adheres to map[string]interface{}
// (JSON encoding does not support map[interface{}]interface{})
func EnsureStringMaps(stringMap StringMap) StringMap {
	stringMap_, _ := MapsToStringMaps(stringMap)
	return stringMap_.(StringMap)
}

// Recursive
func StringMapToMap(stringMap StringMap) Map {
	map_ := make(Map)
	for key, value := range stringMap {
		map_[key], _ = Normalize(value)
	}
	return map_
}

// Recursive
func MapToStringMap(map_ Map) StringMap {
	stringMap := make(StringMap)
	for key, value := range map_ {
		stringMap[yamlkeys.KeyString(key)], _ = MapsToStringMaps(value)
	}
	return stringMap
}
