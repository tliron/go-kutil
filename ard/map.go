package ard

import (
	"github.com/tliron/yamlkeys"
)

func MapsToStringMaps(value Value) (Value, bool) {
	switch value_ := value.(type) {
	case Map:
		return MapToStringMap(value_), true

	case StringMap:
		changedStringMap := make(StringMap)
		changed := false
		for key, element := range value_ {
			var changed_ bool
			if element, changed_ = MapsToStringMaps(element); changed_ {
				changed = true
			}
			changedStringMap[key] = element
		}
		if changed {
			return changedStringMap, true
		}

	case List:
		changedList := make(List, len(value_))
		changed := false
		for index, element := range value_ {
			var changed_ bool
			if element, changed_ = MapsToStringMaps(element); changed_ {
				changed = true
			}
			changedList[index] = element
		}
		if changed {
			return changedList, true
		}
	}

	return value, false
}

// Ensure data adheres to map[string]any
// (JSON encoding does not support map[any]any)
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
