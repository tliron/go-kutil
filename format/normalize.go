package format

func Normalize(data interface{}) (interface{}, error) {
	// TODO: not very efficient
	if code, err := EncodeYAML(data, " ", false); err == nil {
		return DecodeYAML(code)
	} else {
		return nil, err
	}
}
