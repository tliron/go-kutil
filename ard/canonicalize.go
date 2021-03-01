package ard

func Canonicalize(value Value) (Value, error) {
	// Try CBOR first (faster), then YAML, and finally Compatible JSON
	if value, err := RoundtripCBOR(value); err == nil {
		return value, nil
	} else {
		if value, err := RoundtripYAML(value); err == nil {
			return value, nil
		} else {
			return RoundtripCompatibleJSON(value)
		}
	}
}
