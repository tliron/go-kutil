package protobuf

import (
	"fmt"
	"reflect"

	"github.com/tliron/kutil/util"
	"google.golang.org/protobuf/types/known/structpb"
)

// A more flexible version of structpb.NewStruct that accepts structs, automatically
// converts number types, and supports variable slice and map types.
func NewStruct(value any) (*structpb.Struct, error) {
	if map_, err := ToCompatibleStringMap(value); err == nil {
		return structpb.NewStruct(map_)
	} else {
		return nil, err
	}
}

func ToCompatibleStringMap(value any) (map[string]any, error) {
	if value_, err := ToCompatibleValue(value); err == nil {
		if map_, ok := value_.(map[string]any); ok {
			return map_, nil
		} else {
			return nil, fmt.Errorf("not a map: %T", value_)
		}
	} else {
		return nil, err
	}
}

func ToCompatibleValue(value any) (any, error) {
	if value == nil {
		return nil, nil
	}

	switch value_ := value.(type) {
	// Supported types
	case bool, float64, string:
		return value, nil

	// Other numbers
	case int:
		return float64(value_), nil
	case int8:
		return float64(value_), nil
	case int16:
		return float64(value_), nil
	case int32:
		return float64(value_), nil
	case int64:
		return float64(value_), nil
	case uint:
		return float64(value_), nil
	case uint8:
		return float64(value_), nil
	case uint16:
		return float64(value_), nil
	case uint32:
		return float64(value_), nil
	case uint64:
		return float64(value_), nil
	case float32:
		return float64(value_), nil

	case []any:
		length := len(value_)
		list := make([]any, length)
		var err error
		for index, v := range value_ {
			if list[index], err = ToCompatibleValue(v); err != nil {
				return nil, err
			}
		}
		return list, nil

	case map[string]any:
		map_ := make(map[string]any)
		var err error
		for k, v := range value_ {
			if map_[k], err = ToCompatibleValue(v); err != nil {
				return nil, err
			}
		}
		return map_, nil
	}

	type_ := reflect.TypeOf(value)
	switch type_.Kind() {
	case reflect.Pointer:
		value_ := reflect.ValueOf(value).Elem().Interface()
		return ToCompatibleValue(value_)

	case reflect.Slice:
		value_ := reflect.ValueOf(value)
		length := value_.Len()
		list := make([]any, length)
		var err error
		for index := 0; index < length; index++ {
			v := value_.Index(index).Interface()
			if list[index], err = ToCompatibleValue(v); err != nil {
				return nil, err
			}
		}
		return list, nil

	case reflect.Map:
		value_ := reflect.ValueOf(value)
		keys := value_.MapKeys()
		map_ := make(map[string]any)
		var err error
		for _, k := range keys {
			if k.Type().Kind() != reflect.String {
				return nil, fmt.Errorf("unsupported map key type: %T", k.Interface())
			}
			k_ := k.Interface().(string)
			v := value_.MapIndex(k)
			if map_[k_], err = ToCompatibleValue(v); err != nil {
				return nil, err
			}
		}
		return map_, nil

	case reflect.Struct:
		value_ := reflect.ValueOf(value)
		length := type_.NumField()
		map_ := make(map[string]any)
		var err error
		for index := 0; index < length; index++ {
			field := type_.Field(index)
			v := value_.Field(index).Interface()
			if map_[util.ToCamelCase(field.Name)], err = ToCompatibleValue(v); err != nil {
				return nil, err
			}
		}
		return map_, nil
	}

	return nil, fmt.Errorf("unsupported type: %T", value)
}

func UnpackStringMap(value any, ptrToStruct any) error {
	type_ := reflect.TypeOf(ptrToStruct)
	if (type_.Kind() != reflect.Pointer) && (type_.Elem().Kind() != reflect.Struct) {
		return fmt.Errorf("not a pointer to a struct: %T", ptrToStruct)
	}
	return UnpackReflectValue(value, reflect.ValueOf(ptrToStruct).Elem())
}

func UnpackReflectValue(value any, field reflect.Value) error {
	type_ := field.Type()
	switch type_.Kind() {
	case reflect.String:
		if value_, ok := value.(string); ok {
			field.SetString(value_)
		} else {
			return fmt.Errorf("not a string: %T", value)
		}

	case reflect.Float64, reflect.Float32:
		if value_, ok := value.(float64); ok {
			field.SetFloat(value_)
		} else {
			return fmt.Errorf("not a float64: %T", value)
		}

	case reflect.Int64, reflect.Int32, reflect.Int8, reflect.Int:
		if value_, ok := value.(float64); ok {
			field.SetInt(int64(value_))
		} else {
			return fmt.Errorf("not a float64: %T", value)
		}

	case reflect.Slice:
		if list, ok := value.([]any); ok {
			length := len(list)
			slice := reflect.MakeSlice(type_, length, length)
			for index, v := range list {
				if err := UnpackReflectValue(v, slice.Index(index)); err != nil {
					return err
				}
			}
			field.Set(slice)
		} else {
			return fmt.Errorf("not a list: %T", value)
		}

	case reflect.Map:
		if map_, ok := value.(map[string]any); ok {
			map__ := reflect.MakeMap(type_)
			for k, v := range map_ {
				k_ := reflect.New(type_.Key())
				v_ := reflect.New(type_.Elem())
				if err := UnpackReflectValue(k, k_); err != nil {
					return err
				}
				if err := UnpackReflectValue(v, v_); err != nil {
					return err
				}
				map__.SetMapIndex(k_, v_)
			}
			field.Set(map__)
		} else {
			return fmt.Errorf("not a map: %T", value)
		}

	case reflect.Struct:
		if map_, ok := value.(map[string]any); ok {
			length := type_.NumField()
			for index := 0; index < length; index++ {
				field_ := type_.Field(index)
				key := util.ToCamelCase(field_.Name)
				if v, ok := map_[key]; ok {
					if err := UnpackReflectValue(v, field.Field(index)); err != nil {
						return err
					}
				} else {
					return fmt.Errorf("missing key value: %s", key)
				}
			}
		} else {
			return fmt.Errorf("not a map: %T", value)
		}

	default:
		return fmt.Errorf("unsupported field type: %s", field.Type().String())
	}

	return nil
}
