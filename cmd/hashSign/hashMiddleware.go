package hashsign

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log/slog"
	"net/http"
)

func isHashValid(hash string, key string, data []byte) bool {
	h := hmac.New(sha256.New, []byte(key))
	h.Write(data)
	expectedHash := h.Sum(nil)
	receivedHash, err := hex.DecodeString(hash)
	if err != nil {
		return false
	}
	return hmac.Equal(expectedHash, receivedHash)
}

func WithHash(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writer := w
		hw := NewHashWriter(w, "secretkey")
		if h := r.Header.Get("HashSHA256"); h != "" {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				slog.Error("Error while reading body", "Error: ", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			defer r.Body.Close()

			if !isHashValid(h, "secretkey", body) {
				slog.Error("Error in hash")
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			r.Body = io.NopCloser(bytes.NewBuffer(body))
			writer = hw
		}
		next.ServeHTTP(writer, r)
	}
}
