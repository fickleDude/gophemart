package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w *gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

type gzipReader struct {
	io.ReadCloser
	Reader *gzip.Reader
}

func (r *gzipReader) Close() error {
	return r.Reader.Close()
}

func (r gzipReader) Read(p []byte) (n int, err error) {
	return r.Reader.Read(p)
}

func Gzip(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//reader
		encoded := strings.Contains(r.Header.Get("Content-Encoding"), "gzip")
		if encoded {
			gzip, err := gzip.NewReader(r.Body)
			if err != nil {
				io.WriteString(w, err.Error())
				return
			}
			defer gzip.Close()

			r.Body = &gzipReader{ReadCloser: r.Body, Reader: gzip}
			r.Header.Set("Content-Type", "application/json")
		}

		//writer
		rw := w
		acceptEncoding := strings.Contains(r.Header.Get("Accept-Encoding"), "gzip")
		if acceptEncoding {
			gzip := gzip.NewWriter(w)
			defer gzip.Close()
			w.Header().Set("Content-Encoding", "gzip")
			rw = &gzipWriter{ResponseWriter: w, Writer: gzip}
		}
		next.ServeHTTP(rw, r)
	})
}
