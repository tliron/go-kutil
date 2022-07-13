package ard

// TODO: not very efficient

func Canonicalize(value Value) (Value, error) {
	// Try CBOR first (fastest), then Compatible JSON, and finally YAML
	if value, err := RoundtripCBOR(value); err == nil {
		return value, nil
	} else if value, err := RoundtripCompatibleJSON(value); err == nil {
		return value, nil
	} else {
		return RoundtripYAML(value)
	}
}
