// Package pool provides the two-tier sync.Pool used by every Fluent
// render call. Buffers below [Threshold] are recycled through the
// small pool; larger buffers go through the large pool, and anything
// above [MaxPoolSize] is discarded so the pool cannot retain
// pathological allocations.
//
// Pooling is on by default and can be toggled at runtime with
// [Enable] and [Disable]. Configure sizing through [SetThreshold] and
// [SetMaxPoolSize], and observe what the pool is doing by attaching a
// JSONL writer with [SetDiagnostics]. Most callers reach the pool
// indirectly through [github.com/jpl-au/fluent.NewBuffer] rather than
// importing this package.
package pool

import (
	"bytes"
	"encoding/json"
	"io"
	"sync"
	"sync/atomic"
)

// Pool behaviour - configure via the exported setter functions. Held
// atomically because Get and Put read these on every render while the
// setters may run at any time; plain variables would be a data race.
var (
	poolThreshold atomic.Int64 // Threshold defines the pool (small, large) the builder is returned to.
	maxPoolSize   atomic.Int64 // Maximum size to keep in pools - discard larger buffers.
)

// enabled controls whether sync.Pool optimisations are enabled globally.
// It is safe to toggle at runtime via atomic operations.
var enabled atomic.Bool

// diagnostics receives JSONL entries for every Get and Put when set.
// Protected by diagMu so SetDiagnostics is safe to call at any time.
var (
	diagWriter io.Writer
	diagMu     sync.Mutex
)

func init() {
	enabled.Store(true) // Enable pool by default
	poolThreshold.Store(4 * 1024)
	maxPoolSize.Store(256 * 1024)
}

// Labels for diagnostic output - identify which pool handled the operation.
const (
	small   = "small"
	large   = "large"
	fresh   = "new"
	discard = "discard"
)

// Pool instances for any poolable objects
var (
	smallPool = sync.Pool{
		New: func() any {
			return &bytes.Buffer{}
		},
	}
	largePool = sync.Pool{
		New: func() any {
			return &bytes.Buffer{}
		},
	}
)

// Enable turns on sync.Pool optimisations. Pooling is enabled by
// default; call this only after a prior [Disable].
func Enable() {
	enabled.Store(true)
}

// Disable turns off sync.Pool optimisations. Subsequent [Get] calls
// return fresh buffers and [Put] becomes a no-op until [Enable] is
// called.
func Disable() {
	enabled.Store(false)
}

// Enabled reports whether pool optimisations are currently enabled.
func Enabled() bool {
	return enabled.Load()
}

// Get retrieves a buffer from the pool, sized according to the hint.
// If pooling is disabled, it returns a new buffer.
func Get(hint int) *bytes.Buffer {
	if !Enabled() {
		return bytes.NewBuffer(make([]byte, 0, hint))
	}

	var src *sync.Pool
	var pool string
	if hint < int(poolThreshold.Load()) {
		src = &smallPool
		pool = small
	} else {
		src = &largePool
		pool = large
	}

	if p := src.Get(); p != nil {
		buf := p.(*bytes.Buffer) //nolint:forcetypeassert // Pool only contains *bytes.Buffer
		buf.Reset()
		if hint > 0 {
			buf.Grow(hint)
		}
		emitDiag("get", hint, 0, buf.Cap(), pool)
		return buf
	}

	buf := bytes.NewBuffer(make([]byte, 0, hint))
	emitDiag("get", hint, 0, buf.Cap(), fresh)
	return buf
}

// Put returns a buffer to the pool.
// If pooling is disabled or the buffer is too large, it is discarded.
func Put(buf *bytes.Buffer) {
	if !Enabled() || buf == nil {
		return
	}

	cap := buf.Cap()

	// Oversized buffers are always discarded rather than pooled, so one
	// pathological render cannot pin its memory for the process lifetime
	// (the sync.Pool pathology fixed the same way in fmt: golang/go
	// issue 23199). Rebounding to a cap-sized buffer instead measured
	// strictly worse on every workload - see benchmarks/pool-ab.
	if cap > int(maxPoolSize.Load()) {
		emitDiag("put", 0, buf.Len(), cap, discard)
		return
	}

	// Capture length before reset so diagnostics show how much was used.
	length := buf.Len()
	buf.Reset()

	if cap < int(poolThreshold.Load()) {
		emitDiag("put", 0, length, cap, small)
		smallPool.Put(buf)
	} else {
		emitDiag("put", 0, length, cap, large)
		largePool.Put(buf)
	}
}

// Configuration setters

// SetThreshold sets the size threshold in bytes between the small and
// large pools. Buffers with capacity below the threshold are recycled
// through the small pool; the rest go through the large pool. Safe to
// call at any time, including while renders are in flight.
func SetThreshold(size int) {
	poolThreshold.Store(int64(size))
}

// SetMaxPoolSize configures the maximum buffer size in bytes to keep
// in pools. Buffers larger than size are discarded when returned, so
// one pathological render cannot pin its memory in the pool. Safe to
// call at any time, including while renders are in flight.
func SetMaxPoolSize(size int) {
	maxPoolSize.Store(int64(size))
}

// Configuration getters

// Threshold returns the size threshold in bytes between the small and
// large pools.
func Threshold() int {
	return int(poolThreshold.Load())
}

// MaxPoolSize returns the maximum buffer size in bytes that the pool
// will retain.
func MaxPoolSize() int {
	return int(maxPoolSize.Load())
}

// SetDiagnostics enables pool diagnostics. When w is non-nil, every Get
// and Put writes a JSONL entry to w. Pass nil to disable. Safe to call
// at any time - typically called once at application startup.
func SetDiagnostics(w io.Writer) {
	diagMu.Lock()
	diagWriter = w
	diagMu.Unlock()
}

// diagEntry is a single JSONL record written by the diagnostics writer.
type diagEntry struct {
	Op       string `json:"op"`
	Hint     int    `json:"hint,omitempty"`
	Length   int    `json:"len"`
	Capacity int    `json:"cap"`
	Pool     string `json:"pool"`
}

func emitDiag(op string, hint, length, capacity int, pool string) {
	diagMu.Lock()
	w := diagWriter
	diagMu.Unlock()

	if w == nil {
		return
	}

	entry := diagEntry{
		Op:       op,
		Hint:     hint,
		Length:   length,
		Capacity: capacity,
		Pool:     pool,
	}
	line, _ := json.Marshal(entry)
	line = append(line, '\n')

	diagMu.Lock()
	w.Write(line) //nolint:errcheck // best-effort diagnostics
	diagMu.Unlock()
}
