package format

import (
	"encoding/xml"
	"reflect"

	"github.com/tliron/yamlkeys"
)

func ToXMLWritable(data interface{}) interface{} {
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
			slice[index] = ToXMLWritable(v)
		}
		return slice
	case reflect.Map:
		// Convert to slice of XMLMapEntry
		slice := make([]XMLMapEntry, value.Len())
		for index, key := range value.MapKeys() {
			k := yamlkeys.KeyData(key.Interface())
			v := value.MapIndex(key).Interface()
			slice[index] = XMLMapEntry{ToXMLWritable(k), ToXMLWritable(v)}
		}
		return XMLMap{slice}
	}

	return data
}

//
// XMLMap
//

type XMLMap struct {
	Entries []XMLMapEntry
}

var XMLMapStartElement = xml.StartElement{Name: xml.Name{Local: "map"}}

// xml.Marshaler interface
func (self XMLMap) MarshalXML(encoder *xml.Encoder, start xml.StartElement) error {
	if err := encoder.EncodeToken(XMLMapStartElement); err == nil {
		if err := encoder.Encode(self.Entries); err == nil {
			return encoder.EncodeToken(XMLMapStartElement.End())
		} else {
			return err
		}
	} else {
		return err
	}
}

//
// XMLMapEntry
//

type XMLMapEntry struct {
	Key   interface{}
	Value interface{}
}

var XMLMapEntryStart = xml.StartElement{Name: xml.Name{Local: "entry"}}
var XMLKeyStart = xml.StartElement{Name: xml.Name{Local: "key"}}
var XMLValueStart = xml.StartElement{Name: xml.Name{Local: "value"}}

// xml.Marshaler interface
func (self XMLMapEntry) MarshalXML(encoder *xml.Encoder, start xml.StartElement) error {
	if err := encoder.EncodeToken(XMLMapEntryStart); err == nil {
		if err := encoder.EncodeElement(self.Key, XMLKeyStart); err == nil {
			if err := encoder.EncodeToken(XMLValueStart); err == nil {
				if err := encoder.Encode(self.Value); err == nil {
					if err := encoder.EncodeToken(XMLValueStart.End()); err == nil {
						return encoder.EncodeToken(XMLMapEntryStart.End())
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
