package logger

import (
	"errors"
	"log/slog"
	"net/http"
	"os"
	"time"
)

type ResData struct {
	statusCode int
	size       int
}

type LoggingResWriter struct {
	http.ResponseWriter
	resData *ResData
}

func (lrw *LoggingResWriter) Write(data []byte) (int, error) {
	size, err := lrw.ResponseWriter.Write(data)
	if err != nil {
		return 0, errors.New("error while writing response")
	}
	lrw.resData.size += size
	return size, nil
}

func (lrw *LoggingResWriter) WriteHeader(statusCode int) {
	lrw.ResponseWriter.WriteHeader(statusCode)
	lrw.resData.statusCode = statusCode
}

func newLrw(rw http.ResponseWriter) *LoggingResWriter {
	resData := &ResData{}
	return &LoggingResWriter{
		rw,
		resData,
	}
}

var Log *slog.Logger

func Init(level string) {
	Log = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(Log)
}

func WithLogger(h http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		uri := r.RequestURI
		method := r.Method
		lrw := newLrw(w)
		h.ServeHTTP(lrw, r)
		ellapsed := time.Since(start)

		slog.Info("request obtained", "uri", uri, "method", method, "ellapsed", ellapsed, "status", lrw.resData.statusCode, "size", lrw.resData.size)
	})
}
