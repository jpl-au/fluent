package node

import (
	"bytes"
	"testing"
)

// TestFuncRendersOutput verifies that Func calls the function at render
// time and produces its output.
func TestFuncRendersOutput(t *testing.T) {
	got := string(Func(func() Node {
		return stub("hello")
	}).RenderBytes())
	if got != "hello" {
		t.Errorf("Func should render function output, got %q", got)
	}
}

// TestFuncNilReturn verifies that Func safely renders nothing when the
// function returns nil.
func TestFuncNilReturn(t *testing.T) {
	got := string(Func(func() Node { return nil }).RenderBytes())
	if got != "" {
		t.Errorf("Func returning nil should render nothing, got %q", got)
	}
}

// TestFuncNilFunction verifies that a Func with a nil function renders
// nothing rather than panicking.
func TestFuncNilFunction(t *testing.T) {
	got := string(Func(nil).RenderBytes())
	if got != "" {
		t.Errorf("Func(nil) should render nothing, got %q", got)
	}
}

// TestFuncRenderBuilder verifies that RenderBuilder writes the same output
// as Render.
func TestFuncRenderBuilder(t *testing.T) {
	f := Func(func() Node { return stub("hello") })
	var buf bytes.Buffer
	f.RenderBuilder(&buf)
	if buf.String() != "hello" {
		t.Errorf("RenderBuilder should write %q, got %q", "hello", buf.String())
	}
}

// TestFuncRenderToWriter verifies that Render writes to a provided writer
// and that WriteTo reports the byte count.
func TestFuncRenderToWriter(t *testing.T) {
	f := Func(func() Node { return stub("hello") })
	var buf bytes.Buffer
	n, err := f.WriteTo(&buf)
	if err != nil {
		t.Fatalf("WriteTo returned error: %v", err)
	}
	if n != int64(buf.Len()) {
		t.Errorf("WriteTo reported %d bytes but wrote %d", n, buf.Len())
	}
	if buf.String() != "hello" {
		t.Errorf("WriteTo should write %q, got %q", "hello", buf.String())
	}

	buf.Reset()
	f.Render(&buf)
	if buf.String() != "hello" {
		t.Errorf("Render should write %q, got %q", "hello", buf.String())
	}
}

// TestFuncNodeReturnsOutput verifies that Nodes evaluates the function
// and returns its output. This ensures tree walkers see the same children
// that Render produces.
func TestFuncNodeReturnsOutput(t *testing.T) {
	f := Func(func() Node { return stub("hello") })
	nodes := f.Nodes()
	if len(nodes) != 1 {
		t.Fatalf("Nodes() should return 1 node, got %d", len(nodes))
	}
	if string(nodes[0].RenderBytes()) != "hello" {
		t.Errorf("Nodes() should return function output, got %q", string(nodes[0].RenderBytes()))
	}
}

// TestFuncNodeNilReturn verifies that Nodes returns empty when the
// function returns nil.
func TestFuncNodeNilReturn(t *testing.T) {
	f := Func(func() Node { return nil })
	nodes := f.Nodes()
	if len(nodes) != 0 {
		t.Errorf("Nodes() with nil return should be empty, got %d nodes", len(nodes))
	}
}

// TestFuncIsDynamic verifies that function components report themselves as
// dynamic since their output depends on a function call.
func TestFuncIsDynamic(t *testing.T) {
	f := Func(func() Node { return nil })
	if !f.IsDynamic() {
		t.Error("function component should be dynamic")
	}
}
