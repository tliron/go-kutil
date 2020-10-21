package util

import (
	"encoding/base64"
)

func ToBase64(bytes []byte) string {
	return base64.StdEncoding.EncodeToString(bytes)

	/*var builder strings.Builder
	encoder := base64.NewEncoder(base64.StdEncoding, &builder)
	if _, err := encoder.Write(bytes); err == nil {
		if err := encoder.Close(); err == nil {
			return builder.String(), nil
		} else {
			return "", err
		}
	} else {
		return "", err
	}*/
}

func FromBase64(b64 string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(b64)

	/*reader := strings.NewReader(b64)
	decoder := base64.NewDecoder(base64.StdEncoding, reader)
	return ioutil.ReadAll(decoder)*/
}
