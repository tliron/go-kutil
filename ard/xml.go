package ard

import (
	"encoding/xml"
	"fmt"
	"reflect"
	"strconv"

	"github.com/beevik/etree"
	"github.com/tliron/yamlkeys"
)

var CompatibleXMLListTag = "list"
var CompatibleXMLMapTag = "map"
var CompatibleXMLMapEntryTag = "entry"
var CompatibleXMLMapEntryKeyTag = "key"
var CompatibleXMLMapEntryValueTag = "value"

func ToCompatibleXML(data interface{}) interface{} {
	if data == nil {
		return nil
	}

	value := reflect.ValueOf(data)

	switch value.Type().Kind() {
	case reflect.Slice:
		length := value.Len()
		slice := make([]interface{}, length)
		for index := 0; index < length; index++ {
			v := value.Index(index).Interface()
			slice[index] = ToCompatibleXML(v)
		}
		return CompatibleXMLList{slice}

	case reflect.Map:
		// Convert to slice of XMLMapEntry
		slice := make([]CompatibleXMLMapEntry, value.Len())
		for index, key := range value.MapKeys() {
			k := yamlkeys.KeyData(key.Interface())
			v := value.MapIndex(key).Interface()
			slice[index] = CompatibleXMLMapEntry{
				Key:   ToCompatibleXML(k),
				Value: ToCompatibleXML(v),
			}
		}
		return CompatibleXMLMap{slice}
	}

	return data
}

func FromCompatibleXML(element *etree.Element) (interface{}, error) {
	fmt.Println(element.Tag)

	switch element.Tag {
	case CompatibleXMLListTag:
		children := element.ChildElements()
		list := make(List, len(children))
		for index, entry := range children {
			if entry_, err := FromCompatibleXML(entry); err == nil {
				list[index] = entry_
			} else {
				return nil, err
			}
		}
		return list, nil

	case CompatibleXMLMapTag:
		map_ := make(Map)
		for _, entry := range element.ChildElements() {
			if entry_, err := NewCompatibleXMLMapEntry(entry); err == nil {
				//fmt.Printf("%T\n", entry_.Key)
				map_[entry_.Key] = entry_.Value
			} else {
				return nil, err
			}
		}
		return map_, nil

	case "string":
		return element.Text(), nil

	case "int":
		if int_, err := strconv.ParseInt(element.Text(), 10, 64); err == nil {
			return int(int_), nil
		} else {
			return nil, err
		}

	case "int64":
		return strconv.ParseInt(element.Text(), 10, 64)

	case "uint":
		if uint_, err := strconv.ParseUint(element.Text(), 10, 64); err == nil {
			return uint(uint_), nil
		} else {
			return nil, err
		}

	case "uint64":
		return strconv.ParseUint(element.Text(), 10, 64)

	case "float64":
		return strconv.ParseFloat(element.Text(), 64)

	case "bool":
		return strconv.ParseBool(element.Text())

	default:
		return nil, fmt.Errorf("element has unsupported tag: %s", xmlElementToString(element))
	}
}

//
// CompatibleXMLList
//

type CompatibleXMLList struct {
	Entries []interface{}
}

// xml.Marshaler interface
func (self CompatibleXMLList) MarshalXML(encoder *xml.Encoder, start xml.StartElement) error {
	listStart := xml.StartElement{Name: xml.Name{Local: CompatibleXMLListTag}}

	if err := encoder.EncodeToken(listStart); err == nil {
		if err := encoder.Encode(self.Entries); err == nil {
			return encoder.EncodeToken(listStart.End())
		} else {
			return err
		}
	} else {
		return err
	}
}

//
// CompatibleXMLMap
//

type CompatibleXMLMap struct {
	Entries []CompatibleXMLMapEntry
}

// xml.Marshaler interface
func (self CompatibleXMLMap) MarshalXML(encoder *xml.Encoder, start xml.StartElement) error {
	mapStart := xml.StartElement{Name: xml.Name{Local: CompatibleXMLMapTag}}

	if err := encoder.EncodeToken(mapStart); err == nil {
		if err := encoder.Encode(self.Entries); err == nil {
			return encoder.EncodeToken(mapStart.End())
		} else {
			return err
		}
	} else {
		return err
	}
}

//
// CompatibleXMLMapEntry
//

type CompatibleXMLMapEntry struct {
	Key   interface{}
	Value interface{}
}

// xml.Marshaler interface
func (self CompatibleXMLMapEntry) MarshalXML(encoder *xml.Encoder, start xml.StartElement) error {
	mapEntryStart := xml.StartElement{Name: xml.Name{Local: CompatibleXMLMapEntryTag}}
	keyStart := xml.StartElement{Name: xml.Name{Local: CompatibleXMLMapEntryKeyTag}}
	valueStart := xml.StartElement{Name: xml.Name{Local: CompatibleXMLMapEntryValueTag}}

	if err := encoder.EncodeToken(mapEntryStart); err == nil {
		if err := encoder.EncodeToken(keyStart); err == nil {
			if err := encoder.Encode(self.Key); err == nil {
				if err := encoder.EncodeToken(keyStart.End()); err == nil {
					if err := encoder.EncodeToken(valueStart); err == nil {
						if err := encoder.Encode(self.Value); err == nil {
							if err := encoder.EncodeToken(valueStart.End()); err == nil {
								return encoder.EncodeToken(mapEntryStart.End())
							} else {
								return err
							}
						} else {
							return err
						}
					} else {
						return err
					}
				} else {
					return err
				}
			} else {
				return err
			}
		} else {
			return err
		}
	} else {
		return err
	}
}

func NewCompatibleXMLMapEntry(element *etree.Element) (CompatibleXMLMapEntry, error) {
	var self CompatibleXMLMapEntry

	if element.Tag == CompatibleXMLMapEntryTag {
		for _, child := range element.ChildElements() {
			switch child.Tag {
			case CompatibleXMLMapEntryKeyTag:
				if key, err := getXmlElementSingleChild(child); err == nil {
					self.Key = key
				} else {
					return self, err
				}

			case CompatibleXMLMapEntryValueTag:
				if value, err := getXmlElementSingleChild(child); err == nil {
					self.Value = value
				} else {
					return self, err
				}

			default:
				return self, fmt.Errorf("element has unsupported tag: %s", xmlElementToString(element))
			}
		}
	}

	return self, nil
}

// Utilities

func xmlElementToString(element *etree.Element) string {
	document := etree.NewDocument()
	document.SetRoot(element)
	if s, err := document.WriteToString(); err == nil {
		return s
	} else {
		return element.GetPath()
	}
}

func getXmlElementSingleChild(element *etree.Element) (interface{}, error) {
	children := element.ChildElements()
	length := len(children)
	if length == 1 {
		return FromCompatibleXML(children[0])
	} else if length == 0 {
		return nil, nil
	} else {
		return nil, fmt.Errorf("element has more than one child: %s", xmlElementToString(element))
	}
}
