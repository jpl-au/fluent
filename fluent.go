package fluent

import (
	"bytes"

	"github.com/jpl-au/fluent/pool"
)

// NewBuffer creates a new bytes.Buffer instance with optional pooling.
func NewBuffer(hint ...int) *bytes.Buffer {
	h := 0
	if len(hint) > 0 {
		h = hint[0]
	}
	return pool.Get(h)
}

// PutBuffer returns a bytes.Buffer to the pool for reuse
func PutBuffer(buf *bytes.Buffer) {
	pool.Put(buf)
}
