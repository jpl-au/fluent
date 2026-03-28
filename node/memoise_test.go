package node

import (
	"bytes"
	"testing"
)

func TestMemoisedRendersOutput(t *testing.T) {
	got := string(Memoise("v1", func() Node {
		return stub("hello")
	}).Render())
	if got != "hello" {
		t.Errorf("Memoise should render closure output, got %q", got)
	}
}

func TestMemoisedNilClosure(t *testing.T) {
	got := string(Memoise("v1", nil).Render())
	if got != "" {
		t.Errorf("Memoise with nil closure should render nothing, got %q", got)
	}
}

func TestMemoisedNilReturn(t *testing.T) {
	got := string(Memoise("v1", func() Node { return nil }).Render())
	if got != "" {
		t.Errorf("Memoise returning nil should render nothing, got %q", got)
	}
}

func TestMemoisedRenderBuilder(t *testing.T) {
	m := Memoise("v1", func() Node { return stub("hello") })
	var buf bytes.Buffer
	m.RenderBuilder(&buf)
	if buf.String() != "hello" {
		t.Errorf("RenderBuilder should write %q, got %q", "hello", buf.String())
	}
}

func TestMemoisedRenderToWriter(t *testing.T) {
	m := Memoise("v1", func() Node { return stub("hello") })
	var buf bytes.Buffer
	result := m.Render(&buf)
	if result != nil {
		t.Error("Render(writer) should return nil")
	}
	if buf.String() != "hello" {
		t.Errorf("Render(writer) should write %q, got %q", "hello", buf.String())
	}
}

func TestMemoisedNodesReturnsOutput(t *testing.T) {
	m := Memoise("v1", func() Node { return stub("hello") })
	nodes := m.Nodes()
	if len(nodes) != 1 {
		t.Fatalf("Nodes() should return 1 node, got %d", len(nodes))
	}
	if string(nodes[0].Render()) != "hello" {
		t.Errorf("Nodes() should return closure output, got %q", string(nodes[0].Render()))
	}
}

func TestMemoisedNodesNilClosure(t *testing.T) {
	m := Memoise("v1", nil)
	if len(m.Nodes()) != 0 {
		t.Error("Nodes() with nil closure should be empty")
	}
}

func TestMemoisedNodesNilReturn(t *testing.T) {
	m := Memoise("v1", func() Node { return nil })
	if len(m.Nodes()) != 0 {
		t.Error("Nodes() with nil return should be empty")
	}
}

func TestMemoisedKey(t *testing.T) {
	m := Memoise(42, func() Node { return stub("hello") })
	if m.MemoiseKey() != 42 {
		t.Errorf("MemoiseKey() should return 42, got %v", m.MemoiseKey())
	}
}

func TestMemoisedRenderCallsClosure(t *testing.T) {
	m := Memoise("v1", func() Node { return stub("from closure") })
	n := m.MemoiseRender()
	if n == nil {
		t.Fatal("MemoiseRender() should return a node")
	}
	if string(n.Render()) != "from closure" {
		t.Errorf("MemoiseRender() should return closure output, got %q", string(n.Render()))
	}
}

func TestMemoisedRenderNilClosure(t *testing.T) {
	m := Memoise("v1", nil)
	if m.MemoiseRender() != nil {
		t.Error("MemoiseRender() with nil closure should return nil")
	}
}

// Verify Memoise satisfies both Node and Memoiser at compile time.
var (
	_ Node     = (*MemoisedNode)(nil)
	_ Memoiser = (*MemoisedNode)(nil)
)
