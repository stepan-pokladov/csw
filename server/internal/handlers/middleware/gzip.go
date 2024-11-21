package middleware

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
)

func GzipDecompressorMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Encoding") != "gzip" && r.Method != "POST" {
			next.ServeHTTP(w, r)
			return
		}

		gzReader, err := gzip.NewReader(r.Body)
		if err != nil {
			http.Error(w, "Failed to create gzip reader", http.StatusInternalServerError)
			return
		}
		defer gzReader.Close()

		decompressedBody := new(bytes.Buffer)
		_, err = io.Copy(decompressedBody, gzReader)
		if err != nil {
			http.Error(w, "Failed to decompress request body", http.StatusInternalServerError)
			return
		}

		r.Body = io.NopCloser(decompressedBody)
		r.Header.Set("Content-Type", "text/jsonl")

		next.ServeHTTP(w, r)
	})
}
