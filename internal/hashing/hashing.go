package hashing

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
)

type hashingWriter struct {
	w         http.ResponseWriter
	secretKey string
}

func NewHashingWriter(w http.ResponseWriter, secretKey string) *hashingWriter {
	return &hashingWriter{
		w:         w,
		secretKey: secretKey,
	}
}

func (h *hashingWriter) Header() http.Header {
	return h.w.Header()
}

func (h *hashingWriter) Write(data []byte) (int, error) {
	h.w.Header().Set("HashSHA256", GetHash(data, h.secretKey))
	return h.w.Write(data)
}

func (h *hashingWriter) WriteHeader(statusCode int) {
	h.w.WriteHeader(statusCode)
}

func GetHash(data []byte, secretKey string) string {
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(data))
	hashBytes := h.Sum(nil)
	hashHex := hex.EncodeToString(hashBytes)
	return hashHex
}
