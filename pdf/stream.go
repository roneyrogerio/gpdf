package pdf

import (
	"bytes"
	"compress/zlib"
	"fmt"
)

// CompressFlate compresses data using the zlib/deflate algorithm.
// PDF's FlateDecode filter expects zlib-wrapped deflate (RFC 1950),
// not raw deflate (RFC 1951).
func CompressFlate(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	w, err := zlib.NewWriterLevel(&buf, zlib.DefaultCompression)
	if err != nil {
		return nil, fmt.Errorf("pdf: failed to create zlib writer: %w", err)
	}
	if _, err := w.Write(data); err != nil {
		return nil, fmt.Errorf("pdf: failed to write zlib data: %w", err)
	}
	if err := w.Close(); err != nil {
		return nil, fmt.Errorf("pdf: failed to close zlib writer: %w", err)
	}
	return buf.Bytes(), nil
}
