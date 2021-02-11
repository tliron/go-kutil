package format

import (
	"bytes"
	templatepkg "text/template"
)

func DecodeYAMLTemplate(code string, data interface{}) (interface{}, error) {
	if template, err := templatepkg.New("").Parse(code); err == nil {
		var buffer bytes.Buffer
		if err := template.Execute(&buffer, data); err == nil {
			return DecodeYAML(buffer.String())
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}
