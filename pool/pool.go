package pool

import (
	"bytes"
	"encoding/json"
	"io"
	"sync"
	"sync/atomic"
)

// Pool behaviour - configure via the exported setter functions.
var (
	poolThreshold    = 4 * 1024   // Threshold defines the pool (small, large) the builder is returned to
	maxPoolSize      = 256 * 1024 // Maximum size to keep in pools - discard larger buffers
	discardOversized = true       // Whether to discard oversized buffers (true by default)
)

// enabled controls whether sync.Pool optimizations are enabled globally.
// Can be safely toggled at runtime using atomic operations.
var enabled atomic.Bool

// diagnostics receives JSONL entries for every Get and Put when set.
// Protected by diagMu so SetDiagnostics is safe to call at any time.
var (
	diagWriter io.Writer
	diagMu     sync.Mutex
)

func init() {
	enabled.Store(true) // Enable pool by default
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

// Enable turns on sync.Pool optimizations
func Enable() {
	enabled.Store(true)
}

// Disable turns off sync.Pool optimizations
func Disable() {
	enabled.Store(false)
}

// Enabled returns whether pool optimizations are currently enabled
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
	if hint < poolThreshold {
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

	if cap > maxPoolSize && discardOversized {
		emitDiag("put", 0, buf.Len(), cap, discard)
		return
	}

	// Capture length before reset so diagnostics show how much was used.
	length := buf.Len()
	buf.Reset()

	if cap < poolThreshold {
		emitDiag("put", 0, length, cap, small)
		smallPool.Put(buf)
	} else {
		emitDiag("put", 0, length, cap, large)
		largePool.Put(buf)
	}
}

// Configuration setters

// SetThreshold sets the size threshold between small and large pools in bytes
func SetThreshold(size int) {
	poolThreshold = size
}

// SetMaxPoolSize configures the maximum buffer size to keep in pools.
// Buffers larger than this will be discarded if drop is true, otherwise
// they will be resized back to maxSize before being pooled.
func SetMaxPoolSize(size int, drop bool) {
	maxPoolSize = size
	discardOversized = drop
}

// Configuration getters

// Threshold returns the size threshold between small and large pools in bytes
func Threshold() int {
	return poolThreshold
}

// MaxPoolSize returns the maximum buffer size to keep in pools in bytes
func MaxPoolSize() int {
	return maxPoolSize
}

// DiscardOversized returns whether oversized buffers should be discarded
func DiscardOversized() bool {
	return discardOversized
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
