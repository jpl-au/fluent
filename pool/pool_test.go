package pool

import (
	"bytes"
	"testing"
)

func TestPoolGetPut(t *testing.T) {
	Enabled()
	defer Disable()
	// Test Get
	buf := Get(100)
	if buf == nil {
		t.Fatal("Get returned nil")
	}
	if buf.Cap() < 100 {
		t.Errorf("Get returned buffer with capacity %d, expected >= 100", buf.Cap())
	}

	// Test Put
	Put(buf)

	// Test Get from pool (should reuse if possible, but hard to guarantee in test without inspecting internals)
	// We can inspect internals since we are in the same package

	// Clear pools for deterministic testing
	smallPool.New = func() any { return &bytes.Buffer{} }
	largePool.New = func() any { return &bytes.Buffer{} }

	// Put a buffer into small pool
	b1 := bytes.NewBuffer(make([]byte, 0, 100))
	Put(b1)

	// Get it back
	b2 := Get(100)
	// In a real sync.Pool, we can't guarantee we get the same object back, but we can verify behavior
	if b2.Cap() < 100 {
		t.Errorf("Get returned buffer with capacity %d, expected >= 100", b2.Cap())
	}
}

func TestPoolThresholds(t *testing.T) {
	Enabled()
	threshold := PoolThreshold()

	// Small buffer
	smallBuf := Get(threshold - 1)
	if smallBuf.Cap() < threshold-1 {
		t.Errorf("Expected capacity >= %d, got %d", threshold-1, smallBuf.Cap())
	}
	Put(smallBuf)

	// Large buffer
	largeBuf := Get(threshold + 1)
	if largeBuf.Cap() < threshold+1 {
		t.Errorf("Expected capacity >= %d, got %d", threshold+1, largeBuf.Cap())
	}
	Put(largeBuf)
}

func TestPoolDisabled(t *testing.T) {
	Enable()
	defer Disable()

	buf := Get(100)
	if buf == nil {
		t.Fatal("Get returned nil when disabled")
	}
	Put(buf) // Should be no-op or safe
}

func TestDiscardOversized(t *testing.T) {
	Enable()
	max := MaxPoolSize()

	// Oversized buffer
	buf := bytes.NewBuffer(make([]byte, 0, max+1))
	Put(buf)
	// We can't easily verify it was discarded without internal counters, but we can ensure it doesn't panic
}
