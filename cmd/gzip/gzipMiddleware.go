package gzip_local

import (
	"log/slog"
	"net/http"
	"strings"
)

func WithGzip(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writer := w
		isGzipSupported := strings.Contains(r.Header.Get("Accept-Encoding"), "gzip")
		isGzipEncoded := strings.Contains(r.Header.Get("Content-Encoding"), "gizp")
		if isGzipSupported {
			cw := NewCompressWriter(w)
			writer = cw
			defer cw.Close()
		}
		if isGzipEncoded {
			cr, err := NewCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				slog.Error("Unable to unzip body")
				return
			}
			r.Body = cr
			defer cr.Close()
		}
		next.ServeHTTP(writer, r)
	}
}
