package middlewares

import (
	"net/http"
	"strings"

	"github.com/bobgromozeka/metrics/internal/compress/gzip"
)

// Gzippify Adds gzip encode/decode middleware into handlers chain (only for servers that implements std http library handlers, not with ctx).
// Checks Accept-Encoding header before decoding and adds Content-Encoding header after encoding.
func Gzippify(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			resultWriter := w
			acceptList := r.Header.Get("Accept-Encoding")
			supportsGzip := strings.Contains(acceptList, gzip.Name)
			if supportsGzip {
				w.Header().Set("Content-Encoding", gzip.Name)
				gzw := gzip.NewGzipWriter(w)
				resultWriter = gzw
				defer gzw.Close()
			}

			contentEncoding := r.Header.Get("Content-Encoding")
			gotGzip := strings.Contains(contentEncoding, gzip.Name)
			if gotGzip {
				gzr, err := gzip.NewGzipReader(r.Body)
				if err != nil {
					resultWriter.WriteHeader(http.StatusInternalServerError)
					return
				}
				r.Body = gzr
				defer gzr.Close()
			}

			next.ServeHTTP(resultWriter, r)
		},
	)
}
