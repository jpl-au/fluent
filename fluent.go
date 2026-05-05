// Package fluent is the root of the Fluent HTML rendering library. It
// exposes the buffer pool used throughout the render pipeline so that
// callers and extensions can borrow and return buffers without taking
// a direct dependency on the [pool] package.
//
// The HTML elements live under [github.com/jpl-au/fluent/html5], the
// core node interfaces under [github.com/jpl-au/fluent/node], and text
// nodes under [github.com/jpl-au/fluent/text]. Optional dot-import
// constructors live under [github.com/jpl-au/fluent/dot].
//
// Most callers never touch this package directly. Calling Render(w)
// on any node uses the pool transparently; reach for [NewBuffer] and
// [PutBuffer] only when composing renders manually with RenderBuilder.
package fluent

import (
	"bytes"
	"io"

	"github.com/jpl-au/fluent/pool"
)

// NewBuffer returns a pooled bytes.Buffer sized to the given hint.
// If no hint is provided, the pool's default sizing applies.
// Always return the buffer with PutBuffer when finished to avoid leaking
// pooled memory.
//
//	buf := fluent.NewBuffer(256)
//	defer fluent.PutBuffer(buf)
//	node.RenderBuilder(buf)
func NewBuffer(hint ...int) *bytes.Buffer {
	h := 0
	if len(hint) > 0 {
		h = hint[0]
	}
	return pool.Get(h)
}

// PutBuffer returns a buffer to the pool for reuse. Passing nil is safe.
// Do not use the buffer after calling PutBuffer - the pool may hand it
// to another caller at any time.
func PutBuffer(buf *bytes.Buffer) {
	pool.Put(buf)
}

// SetPoolDiagnostics controls JSONL diagnostic output for the buffer pool.
// When w is non-nil, every NewBuffer and PutBuffer call writes a single
// JSON line to w. Pass nil to disable. Safe to call at any time.
//
// Each line contains: op ("get"/"put"), hint, len, cap, and pool
// ("small", "large", "new", or "discard").
//
//	f, _ := os.Create("pool.jsonl")
//	defer f.Close()
//	fluent.SetPoolDiagnostics(f)
func SetPoolDiagnostics(w io.Writer) {
	pool.SetDiagnostics(w)
}
