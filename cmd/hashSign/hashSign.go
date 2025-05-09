package hashsign

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
)

type HashWriter struct {
	W          http.ResponseWriter
	Key        string
	HashSHA256 string
}

func NewHashWriter(w http.ResponseWriter, key string) *HashWriter {
	return &HashWriter{
		W:   w,
		Key: key,
	}
}

func (hw *HashWriter) Write(p []byte) (int, error) {
	h := hmac.New(sha256.New, []byte(hw.Key))
	h.Write(p)
	hEnc := base64.StdEncoding.EncodeToString(h.Sum(nil))
	hw.HashSHA256 = hEnc
	hw.W.Header().Set("HashSHA256", hw.HashSHA256)
	return hw.W.Write(p)
}

func (hw *HashWriter) Header() http.Header {
	return hw.W.Header()
}

func (hw *HashWriter) WriteHeader(code int) {
	if code < 300 {
		hw.W.Header().Set("HashSHA256", hw.HashSHA256)
	}
	hw.W.WriteHeader(code)
}
