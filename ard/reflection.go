package ard

import (
	"fmt"
	"reflect"
	"time"

	"github.com/tliron/kutil/reflection"
)

//
// Reflector
//

type Reflector struct {
	IgnoreMissingStructFields bool
	NilMeansZero              bool
}

func NewReflector() *Reflector {
	return &Reflector{}
}

func (self *Reflector) ToComposite(value Value, compositeValuePtr any) error {
	compositeValuePtr_ := reflect.ValueOf(compositeValuePtr)
	if compositeValuePtr_.Kind() == reflect.Pointer {
		return self.ToCompositeReflect(value, compositeValuePtr_)
	} else {
		return fmt.Errorf("not a pointer: %T", compositeValuePtr)
	}
}

func (self *Reflector) ToCompositeReflect(value Value, compositeValue reflect.Value) error {
	compositeType := compositeValue.Type()

	// Dereference pointers
	if compositeType.Kind() == reflect.Pointer {
		if value == nil {
			return nil
		}

		compositeType = compositeType.Elem()
		if compositeValue.IsNil() {
			// Zero value
			compositeValue.Set(reflect.New(compositeType))
		}
		compositeValue = compositeValue.Elem()
	}

	switch value_ := value.(type) {
	case nil:
		if self.NilMeansZero {
			compositeValue.Set(reflect.Zero(compositeType))
		} else {
			kind := compositeValue.Kind()
			if (kind != reflect.Map) && (kind != reflect.Slice) {
				return fmt.Errorf("not a pointer, map, or slice: %s", compositeType.String())
			}
		}

	case string:
		if compositeValue.Kind() == reflect.String {
			compositeValue.SetString(value_)
		} else {
			return fmt.Errorf("not a string: %s", compositeType.String())
		}

	case bool:
		if compositeValue.Kind() == reflect.Bool {
			compositeValue.SetBool(value_)
		} else {
			return fmt.Errorf("not a bool: %s", compositeType.String())
		}

	case int64, int32, int16, int8, int, uint64, uint32, uint16, uint8, uint:
		if reflection.IsInteger(compositeValue.Kind()) {
			compositeValue.SetInt(ToInt64(value_))
		} else if reflection.IsUInteger(compositeValue.Kind()) {
			compositeValue.SetUint(ToUInt64(value_))
		} else {
			return fmt.Errorf("not an integer: %s", compositeType.String())
		}

	case float64, float32:
		if reflection.IsFloat(compositeValue.Kind()) {
			compositeValue.SetFloat(ToFloat64(value_))
		} else {
			return fmt.Errorf("not a float: %s", compositeType.String())
		}

	case time.Time: // as-is values
		if compositeType == reflect.TypeOf(value_) {
			compositeValue.Set(reflect.ValueOf(value_))
		} else {
			return fmt.Errorf("not a %T: %s", value, compositeType.String())
		}

	case List:
		if compositeValue.Kind() == reflect.Slice {
			elemType := compositeType.Elem()
			length := len(value_)
			list := reflect.MakeSlice(reflect.SliceOf(elemType), length, length)
			for index, elem := range value_ {
				if err := self.ToCompositeReflect(elem, list.Index(index)); err != nil {
					return fmt.Errorf("slice element %d %s", index, err.Error())
				}
			}
			compositeValue.Set(list)
		} else {
			return fmt.Errorf("not a slice: %s", compositeType.String())
		}

	case Map:
		switch compositeValue.Kind() {
		case reflect.Map:
			if compositeValue.IsNil() {
				compositeValue.Set(reflect.MakeMap(compositeType))
			}

			keyType := compositeType.Key()
			valueType := compositeType.Elem()
			for k, v := range value_ {
				k_ := reflect.New(keyType)
				if err := self.ToCompositeReflect(k, k_); err == nil {
					v_ := reflect.New(valueType)
					if err := self.ToCompositeReflect(v, v_); err == nil {
						compositeValue.SetMapIndex(k_.Elem(), v_.Elem())
					} else {
						return fmt.Errorf("map value %s", err)
					}
				} else {
					return fmt.Errorf("map key %s", err)
				}
			}

		case reflect.Struct:
			fieldNames := NewFieldNames(compositeType)
			for k, v := range value_ {
				if k_, ok := k.(string); ok {
					if err := self.setStructField(compositeValue, k_, v, fieldNames); err != nil {
						return err
					}
				} else {
					return fmt.Errorf("key not a string: %T", k)
				}
			}

		default:
			return fmt.Errorf("not a map or struct: %s", compositeType.String())
		}

	case StringMap:
		switch compositeValue.Kind() {
		case reflect.Map:
			if compositeValue.IsNil() {
				compositeValue.Set(reflect.MakeMap(compositeType))
			}

			keyType := compositeType.Key()
			valueType := compositeType.Elem()
			for k, v := range value_ {
				k_ := reflect.New(keyType)
				if err := self.ToCompositeReflect(k, k_); err == nil {
					v_ := reflect.New(valueType)
					if err := self.ToCompositeReflect(v, v_); err == nil {
						compositeValue.SetMapIndex(k_.Elem(), v_.Elem())
					} else {
						return fmt.Errorf("map value %s", err)
					}
				} else {
					return fmt.Errorf("map key %s", err)
				}
			}

		case reflect.Struct:
			fieldNames := NewFieldNames(compositeType)
			for k, v := range value_ {
				if err := self.setStructField(compositeValue, k, v, fieldNames); err != nil {
					return err
				}
			}

		default:
			return fmt.Errorf("not a map or struct: %s", compositeType.String())
		}

	default:
		return fmt.Errorf("unsupported type: %s", compositeType.String())
	}

	return nil
}

func (self *Reflector) FromComposite(compositeValue any) (Value, error) {
	return self.FromCompositeReflect(reflect.ValueOf(compositeValue))
}

var time_ time.Time
var timeType = reflect.TypeOf(time_)

func (self *Reflector) FromCompositeReflect(compositeValue reflect.Value) (Value, error) {
	compositeType := compositeValue.Type()

	// Dereference pointers
	if compositeType.Kind() == reflect.Pointer {
		if compositeValue.IsNil() {
			return nil, nil
		}

		compositeType = compositeType.Elem()
		compositeValue = compositeValue.Elem()
	}

	if compositeType == timeType {
		return compositeValue.Interface(), nil
	}

	switch compositeType.Kind() {
	case reflect.String, reflect.Bool, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8, reflect.Int, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8, reflect.Uint, reflect.Float64, reflect.Float32:
		return compositeValue.Interface(), nil

	case reflect.Slice:
		length := compositeValue.Len()
		list := make(List, length)
		for index := 0; index < length; index++ {
			var err error
			if list[index], err = self.FromCompositeReflect(compositeValue.Index(index)); err != nil {
				return nil, fmt.Errorf("list element %d %s", index, err.Error())
			}
		}
		return list, nil

	case reflect.Map:
		map_ := make(Map)
		keys := compositeValue.MapKeys()
		for _, key := range keys {
			if key_, err := self.FromCompositeReflect(key); err == nil {
				elem := compositeValue.MapIndex(key)
				if elem_, err := self.FromCompositeReflect(elem); err == nil {
					map_[key_] = elem_
				} else {
					return nil, fmt.Errorf("map element %q %s", key_, err.Error())
				}
			} else {
				return nil, fmt.Errorf("map key %q %s", key_, err.Error())
			}
		}
		return map_, nil

	case reflect.Struct:
		map_ := make(Map)
		fieldNames := NewFieldNames(compositeType)
		for name, fieldName := range fieldNames {
			elem := compositeValue.FieldByName(fieldName)
			if elem_, err := self.FromCompositeReflect(elem); err == nil {
				map_[name] = elem_
			} else {
				return nil, fmt.Errorf("struct field %q %s", fieldName, err.Error())
			}
		}
		return map_, nil
	}

	return nil, fmt.Errorf("unsupported type: %s", compositeType.String())
}

func (self *Reflector) setStructField(structValue reflect.Value, fieldName string, value Value, fieldNames FieldNames) error {
	field := fieldNames.GetField(structValue, fieldName)

	if !field.IsValid() {
		if self.IgnoreMissingStructFields {
			return nil
		} else {
			return fmt.Errorf("no %q field", fieldName)
		}
	}

	if !field.CanSet() {
		return fmt.Errorf("field %q cannot be set", fieldName)
	}

	if err := self.ToCompositeReflect(value, field); err == nil {
		return nil
	} else {
		return fmt.Errorf("field %q %s", fieldName, err.Error())
	}
}

//
// FieldNames
//

type FieldNames map[string]string // ARD name to struct field name

func NewFieldNames(type_ reflect.Type) FieldNames {
	self := make(FieldNames)
	tags := reflection.GetFieldTagsForType(type_, "ard")

	// Tagged fields
	for fieldName, tag := range tags {
		self[tag] = fieldName
	}

	// Untagged fields
	for _, field := range reflection.GetStructFields(type_) {
		fieldName := field.Name
		if _, ok := tags[fieldName]; !ok {
			self[fieldName] = fieldName
		}
	}

	return self
}

func (self FieldNames) GetField(structValue reflect.Value, name string) reflect.Value {
	return structValue.FieldByName(self[name])
}
