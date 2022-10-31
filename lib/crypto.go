package lib

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"strings"
)

func Hash(src string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(src))
	return strings.TrimRight(base64.StdEncoding.EncodeToString(h.Sum(nil)), "=")
}
