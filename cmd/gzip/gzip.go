package gzip_local

import (
	"compress/gzip"
	"io"
	"net/http"
)

type CompressWriter struct {
	W  http.ResponseWriter
	Zw gzip.Writer
}

func NewCompressWriter(w http.ResponseWriter) *CompressWriter {
	return &CompressWriter{
		W:  w,
		Zw: *gzip.NewWriter(w),
	}
}

func (c *CompressWriter) Header() http.Header {
	return c.W.Header()
}

func (c *CompressWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		c.W.Header().Set("Content-Encoding", "gzip")
	}
	c.W.WriteHeader(statusCode)
}

func (c *CompressWriter) Write(p []byte) (int, error) {
	return c.Zw.Write(p)
}

func (c *CompressWriter) Close() error {
	return c.Zw.Close()
}

type CompressReader struct {
	R  io.ReadCloser
	Zr gzip.Reader
}

func NewCompressReader(r io.ReadCloser) (*CompressReader, error) {
	reader, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	return &CompressReader{
		R:  r,
		Zr: *reader,
	}, nil
}

func (cr *CompressReader) Read(p []byte) (int, error) {
	return cr.Zr.Read(p)
}

func (cr *CompressReader) Close() error {
	err := cr.R.Close()
	if err != nil {
		return err
	}
	return nil
}
