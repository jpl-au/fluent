package node

import (
	"bytes"
	"testing"
)

// TestFuncsRendersOutput verifies that Funcs calls the function at
// render time and produces the output of all returned nodes.
func TestFuncsRendersOutput(t *testing.T) {
	got := string(Funcs(func() []Node {
		return []Node{stub("one"), stub("two")}
	}).Render())
	if got != "onetwo" {
		t.Errorf("Funcs should render all nodes, got %q", got)
	}
}

// TestFuncsEmptySlice verifies that Funcs renders nothing when the
// function returns an empty slice.
func TestFuncsEmptySlice(t *testing.T) {
	got := string(Funcs(func() []Node { return nil }).Render())
	if got != "" {
		t.Errorf("Funcs returning nil should render nothing, got %q", got)
	}
}

// TestFuncsSkipsNilEntries verifies that nil nodes in the returned
// slice are safely skipped rather than causing a panic.
func TestFuncsSkipsNilEntries(t *testing.T) {
	got := string(Funcs(func() []Node {
		return []Node{stub("a"), nil, stub("b")}
	}).Render())
	if got != "ab" {
		t.Errorf("Funcs should skip nil entries, got %q", got)
	}
}

// TestFuncsNilFunction verifies that a Funcs with a nil function
// renders nothing rather than panicking.
func TestFuncsNilFunction(t *testing.T) {
	got := string(Funcs(nil).Render())
	if got != "" {
		t.Errorf("Funcs(nil) should render nothing, got %q", got)
	}
}

// TestFuncsRenderBuilder verifies that RenderBuilder writes the same
// output as Render.
func TestFuncsRenderBuilder(t *testing.T) {
	f := Funcs(func() []Node { return []Node{stub("hello")} })
	var buf bytes.Buffer
	f.RenderBuilder(&buf)
	if buf.String() != "hello" {
		t.Errorf("RenderBuilder should write %q, got %q", "hello", buf.String())
	}
}

// TestFuncsRenderToWriter verifies that Render writes to a provided
// writer and returns nil.
func TestFuncsRenderToWriter(t *testing.T) {
	f := Funcs(func() []Node { return []Node{stub("hello")} })
	var buf bytes.Buffer
	result := f.Render(&buf)
	if result != nil {
		t.Error("Render(writer) should return nil")
	}
	if buf.String() != "hello" {
		t.Errorf("Render(writer) should write %q, got %q", "hello", buf.String())
	}
}

// TestFuncsNodesReturnsOutput verifies that Nodes evaluates the
// function and returns its output. This ensures tree walkers see the same
// children that Render produces.
func TestFuncsNodesReturnsOutput(t *testing.T) {
	f := Funcs(func() []Node { return []Node{stub("a"), stub("b")} })
	nodes := f.Nodes()
	if len(nodes) != 2 {
		t.Fatalf("Nodes() should return 2 nodes, got %d", len(nodes))
	}
	if string(nodes[0].Render()) != "a" {
		t.Errorf("first node should be %q, got %q", "a", string(nodes[0].Render()))
	}
	if string(nodes[1].Render()) != "b" {
		t.Errorf("second node should be %q, got %q", "b", string(nodes[1].Render()))
	}
}

// TestFuncsNodesNilReturn verifies that Nodes returns nil when the
// function returns nil.
func TestFuncsNodesNilReturn(t *testing.T) {
	f := Funcs(func() []Node { return nil })
	nodes := f.Nodes()
	if nodes != nil {
		t.Errorf("Nodes() with nil return should be nil, got %d nodes", len(nodes))
	}
}

// TestFuncsIsDynamic verifies that function node components report
// themselves as dynamic since their output depends on a function call.
func TestFuncsIsDynamic(t *testing.T) {
	f := Funcs(func() []Node { return nil })
	if !f.IsDynamic() {
		t.Error("functions component should be dynamic")
	}
}
