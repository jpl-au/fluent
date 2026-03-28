package node

import (
	"bytes"
	"testing"
)

func TestMemoRendersOutput(t *testing.T) {
	got := string(Memoise("v1", func() Node {
		return stub("hello")
	}).Render())
	if got != "hello" {
		t.Errorf("Memo should render closure output, got %q", got)
	}
}

func TestMemoNilClosure(t *testing.T) {
	got := string(Memoise("v1", nil).Render())
	if got != "" {
		t.Errorf("Memo with nil closure should render nothing, got %q", got)
	}
}

func TestMemoNilReturn(t *testing.T) {
	got := string(Memoise("v1", func() Node { return nil }).Render())
	if got != "" {
		t.Errorf("Memo returning nil should render nothing, got %q", got)
	}
}

func TestMemoRenderBuilder(t *testing.T) {
	m := Memoise("v1", func() Node { return stub("hello") })
	var buf bytes.Buffer
	m.RenderBuilder(&buf)
	if buf.String() != "hello" {
		t.Errorf("RenderBuilder should write %q, got %q", "hello", buf.String())
	}
}

func TestMemoRenderToWriter(t *testing.T) {
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

func TestMemoKey(t *testing.T) {
	m := Memoise(42, func() Node { return stub("hello") })
	if m.MemoKey() != 42 {
		t.Errorf("MemoKey() should return 42, got %v", m.MemoKey())
	}
}

func TestMemoRenderCallsClosure(t *testing.T) {
	m := Memoise("v1", func() Node { return stub("from closure") })
	n := m.MemoRender()
	if n == nil {
		t.Fatal("MemoRender() should return a node")
	}
	if string(n.Render()) != "from closure" {
		t.Errorf("MemoRender() should return closure output, got %q", string(n.Render()))
	}
}

func TestMemoRenderNilClosure(t *testing.T) {
	m := Memoise("v1", nil)
	if m.MemoRender() != nil {
		t.Error("MemoRender() with nil closure should return nil")
	}
}

// Verify Memo satisfies both Node and Memoiser at compile time.
var (
	_ Node     = (*MemoisedNode)(nil)
	_ Memoiser = (*MemoisedNode)(nil)
)
