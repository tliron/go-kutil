package transcribe

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/fxamacker/cbor/v2"
	"github.com/tliron/kutil/ard"
	"github.com/tliron/kutil/util"
)

func Encode(value any, format string, indent string, strict bool) (string, error) {
	switch format {
	case "yaml", "":
		return EncodeYAML(value, indent, strict)

	case "json":
		return EncodeJSON(value, indent)

	case "cjson":
		return EncodeCompatibleJSON(value, indent)

	case "xml":
		return EncodeCompatibleXML(value, indent)

	case "cbor":
		return EncodeCBOR(value)

	case "messagepack":
		return EncodeMessagePack(value)

	case "go":
		return EncodeGo(value, indent)

	default:
		return "", fmt.Errorf("unsupported format: %q", format)
	}
}

func EncodeYAML(value any, indent string, strict bool) (string, error) {
	var writer strings.Builder
	if err := WriteYAML(value, &writer, indent, strict); err == nil {
		return writer.String(), nil
	} else {
		return "", err
	}
}

func EncodeJSON(value any, indent string) (string, error) {
	var writer strings.Builder
	if err := WriteJSON(value, &writer, indent); err == nil {
		s := writer.String()
		if indent == "" {
			// json.Encoder adds a "\n", unlike json.Marshal
			s = strings.TrimRight(s, "\n")
		}
		return s, nil
	} else {
		return "", err
	}
}

func EncodeCompatibleJSON(value any, indent string) (string, error) {
	return EncodeJSON(ard.EnsureCompatibleJSON(value), indent)
}

func EncodeCompatibleXML(value any, indent string) (string, error) {
	var writer strings.Builder
	if err := WriteCompatibleXML(value, &writer, indent); err == nil {
		return writer.String(), nil
	} else {
		return "", err
	}
}

// Encodes to Base64
func EncodeCBOR(value any) (string, error) {
	if bytes, err := cbor.Marshal(value); err == nil {
		return util.ToBase64(bytes), nil
	} else {
		return "", err
	}
}

// Encodes to Base64
func EncodeMessagePack(value any) (string, error) {
	var buffer bytes.Buffer
	encoder := util.NewMessagePackEncoder(&buffer)
	if err := encoder.Encode(value); err == nil {
		return util.ToBase64(buffer.Bytes()), nil
	} else {
		return "", err
	}
}

func EncodeGo(value any, indent string) (string, error) {
	return NewUtterConfig(indent).Sdump(value), nil
}
