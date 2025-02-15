package hashing

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func GetHash(data []byte, secretKey string) string {
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(data))
	hashBytes := h.Sum(nil)
	hashHex := hex.EncodeToString(hashBytes)
	return hashHex
}
