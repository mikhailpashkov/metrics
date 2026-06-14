package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

const HashHeaderKey = "HashSHA256"

func Hash(body []byte, key []byte) string {
	h := hmac.New(sha256.New, key)
	h.Write(body)
	return hex.EncodeToString(h.Sum(nil))
}

func HashEqual(a string, b string) bool {
	return hmac.Equal([]byte(a), []byte(b))
}
