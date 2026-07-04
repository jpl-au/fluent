package node

import (
	"bytes"
	"io"
	"testing"
)

// stubNode is a minimal Node for tests within the node package, avoiding
// import cycles with text or html5 packages.
type stubNode struct {
	content string
}

func stub(content string) *stubNode { return &stubNode{content: content} }
func (s *stubNode) Render(w io.Writer) {
	_, _ = s.WriteTo(w)
}
func (s *stubNode) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write([]byte(s.content))
	return int64(n), err
}
func (s *stubNode) RenderBytes() []byte             { return []byte(s.content) }
func (s *stubNode) RenderBuilder(buf *bytes.Buffer) { buf.WriteString(s.content) }
func (s *stubNode) Nodes() []Node                   { return nil }

// TestWhenTrueRendersChild verifies that When renders its child when the
// condition is true.
func TestWhenTrueRendersChild(t *testing.T) {
	got := string(When(true, stub("visible")).RenderBytes())
	if got != "visible" {
		t.Errorf("When(true) should render child, got %q", got)
	}
}

// TestWhenFalseRendersNothing verifies that When renders nothing when the
// condition is false.
func TestWhenFalseRendersNothing(t *testing.T) {
	got := string(When(false, stub("hidden")).RenderBytes())
	if got != "" {
		t.Errorf("When(false) should render nothing, got %q", got)
	}
}

// TestUnlessFalseRendersChild verifies that Unless renders its child when
// the condition is false.
func TestUnlessFalseRendersChild(t *testing.T) {
	got := string(Unless(false, stub("visible")).RenderBytes())
	if got != "visible" {
		t.Errorf("Unless(false) should render child, got %q", got)
	}
}

// TestUnlessTrueRendersNothing verifies that Unless renders nothing when
// the condition is true.
func TestUnlessTrueRendersNothing(t *testing.T) {
	got := string(Unless(true, stub("hidden")).RenderBytes())
	if got != "" {
		t.Errorf("Unless(true) should render nothing, got %q", got)
	}
}

// TestConditionTrueBranch verifies that Condition renders the true branch
// when the condition is true.
func TestConditionTrueBranch(t *testing.T) {
	got := string(Condition(true).
		True(stub("yes")).
		False(stub("no")).
		RenderBytes())
	if got != "yes" {
		t.Errorf("Condition(true) should render true branch, got %q", got)
	}
}

// TestConditionFalseBranch verifies that Condition renders the false branch
// when the condition is false.
func TestConditionFalseBranch(t *testing.T) {
	got := string(Condition(false).
		True(stub("yes")).
		False(stub("no")).
		RenderBytes())
	if got != "no" {
		t.Errorf("Condition(false) should render false branch, got %q", got)
	}
}

// TestConditionRenderBuilder verifies that RenderBuilder writes the same
// output as Render.
func TestConditionRenderBuilder(t *testing.T) {
	c := Condition(true).True(stub("hello"))
	var buf bytes.Buffer
	c.RenderBuilder(&buf)
	if buf.String() != "hello" {
		t.Errorf("RenderBuilder should write %q, got %q", "hello", buf.String())
	}
}

// TestConditionRenderToWriter verifies that WriteTo and Render write to a
// provided writer.
func TestConditionRenderToWriter(t *testing.T) {
	c := When(true, stub("hello"))
	var buf bytes.Buffer
	n, err := c.WriteTo(&buf)
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
	c.Render(&buf)
	if buf.String() != "hello" {
		t.Errorf("Render should write %q, got %q", "hello", buf.String())
	}
}

// TestConditionNilTrueNode verifies that a nil true branch renders nothing
// rather than panicking.
func TestConditionNilTrueNode(t *testing.T) {
	got := string(Condition(true).True(nil).RenderBytes())
	if got != "" {
		t.Errorf("nil true node should render nothing, got %q", got)
	}
}

// TestConditionNilFalseNode verifies that a nil false branch renders nothing
// rather than panicking.
func TestConditionNilFalseNode(t *testing.T) {
	got := string(Condition(false).False(nil).RenderBytes())
	if got != "" {
		t.Errorf("nil false node should render nothing, got %q", got)
	}
}

// TestConditionNodesReturnsActiveBranch verifies that Nodes returns only
// the branch that will be rendered - the active branch. This ensures tree
// walkers see the same children that Render produces.
func TestConditionNodesReturnsActiveBranch(t *testing.T) {
	trueNode := stub("yes")
	falseNode := stub("no")

	c := Condition(true).True(trueNode).False(falseNode)
	nodes := c.Nodes()
	if len(nodes) != 1 {
		t.Fatalf("Nodes() with condition=true should return 1 node, got %d", len(nodes))
	}
	if string(nodes[0].RenderBytes()) != "yes" {
		t.Errorf("Nodes() with condition=true should return the true branch, got %q", string(nodes[0].RenderBytes()))
	}

	c = Condition(false).True(trueNode).False(falseNode)
	nodes = c.Nodes()
	if len(nodes) != 1 {
		t.Fatalf("Nodes() with condition=false should return 1 node, got %d", len(nodes))
	}
	if string(nodes[0].RenderBytes()) != "no" {
		t.Errorf("Nodes() with condition=false should return the false branch, got %q", string(nodes[0].RenderBytes()))
	}
}

// TestWhenFalseNodesReturnsEmpty verifies that When(false) exposes no
// children - the conditional renders nothing so there are no active nodes.
func TestWhenFalseNodesReturnsEmpty(t *testing.T) {
	c := When(false, stub("hidden"))
	nodes := c.Nodes()
	if len(nodes) != 0 {
		t.Errorf("When(false) Nodes() should return empty, got %d nodes", len(nodes))
	}
}

// TestConditionIsDynamic verifies that conditionals report themselves as
// dynamic since their output depends on a runtime condition.
func TestConditionIsDynamic(t *testing.T) {
	c := When(true, stub("yes"))
	if !c.IsDynamic() {
		t.Error("conditional should be dynamic")
	}
}
