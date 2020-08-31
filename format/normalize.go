package format

func Normalize(data interface{}) (interface{}, error) {
	if code, err := EncodeYAML(data, " ", false); err == nil {
		return DecodeYAML(code)
	} else {
		return nil, err
	}
}
