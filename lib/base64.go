package lib

import (
	"encoding/base64"
	"encoding/json"
	"strings"
)

// Encode exported
func B64Encode(s interface{}) string {

	b, _ := json.Marshal(s)
	return strings.TrimRight(base64.StdEncoding.EncodeToString(b), "=")
}

// Decode exported
func B64Decode(src string) (string, error) {

	if l := len(src) % 4; l > 0 {
		src += strings.Repeat("=", 4-l)
	}
	decoded, err := base64.URLEncoding.DecodeString(src)
	if err != nil {
		e := new(B64DecodeError).SetArgs(err.Error())
		return "", &e
	}
	return string(decoded), nil
}
