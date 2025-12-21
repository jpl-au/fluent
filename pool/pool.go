package pool

import (
	"bytes"
	"sync"
	"sync/atomic"
)

// Global pool configuration - internal variables, access via getter/setter methods
var (
	poolThreshold    = 4 * 1024   // Threshold defines the pool (small, large) the builder is returned to
	maxPoolSize      = 256 * 1024 // Maximum size to keep in pools - discard larger buffers
	discardOversized = true       // Whether to discard oversized buffers (true by default)
)

// enabled controls whether sync.Pool optimizations are enabled globally.
// Can be safely toggled at runtime using atomic operations.
var enabled atomic.Bool

func init() {
	enabled.Store(true) // Enable pool by default
}

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

	var pooled *bytes.Buffer
	if hint < poolThreshold {
		if p := smallPool.Get(); p != nil {
			pooled = p.(*bytes.Buffer)
		}
	} else {
		if p := largePool.Get(); p != nil {
			pooled = p.(*bytes.Buffer)
		}
	}

	if pooled != nil {
		pooled.Reset()
		if hint > 0 {
			pooled.Grow(hint)
		}
		return pooled
	}

	return bytes.NewBuffer(make([]byte, 0, hint))
}

// Put returns a buffer to the pool.
// If pooling is disabled or the buffer is too large, it is discarded.
func Put(buf *bytes.Buffer) {
	if !Enabled() || buf == nil {
		return
	}

	cap := buf.Cap()

	// Check if buffer is oversized
	if cap > maxPoolSize {
		if discardOversized {
			// Discard oversized buffers to prevent memory bloat
			return
		}
	}

	buf.Reset()
	// Route to appropriate pool based on capacity
	if cap < poolThreshold {
		smallPool.Put(buf)
	} else {
		largePool.Put(buf)
	}
}

// Configuration setters

// SetPoolThreshold sets the size threshold between small and large pools in bytes
func SetPoolThreshold(size int) {
	poolThreshold = size
}

// SetMaxPoolSize configures the maximum buffer size to keep in pools.
// Buffers larger than this will be discarded if discard is true, otherwise
// they will be resized back to maxSize before being pooled.
func SetMaxPoolSize(size int, discard bool) {
	maxPoolSize = size
	discardOversized = discard
}

// Configuration getters

// PoolThreshold returns the size threshold between small and large pools in bytes
func PoolThreshold() int {
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
