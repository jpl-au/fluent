package node

import (
	"bytes"
	"testing"
)

func TestEmptyRendersNothing(t *testing.T) {
	if got := Empty().RenderBytes(); len(got) != 0 {
		t.Errorf("Render() = %q, want empty", got)
	}

	var buf bytes.Buffer
	Empty().RenderBuilder(&buf)
	if buf.Len() != 0 {
		t.Errorf("RenderBuilder wrote %q, want nothing", buf.String())
	}

	if n := Empty().Nodes(); n != nil {
		t.Errorf("Nodes() = %v, want nil", n)
	}
}

// TestEmptyIsNonNil guards the whole point of Empty: it is a real, non-nil node,
// so it is safe wherever a typed nil pointer would dereference.
func TestEmptyIsNonNil(t *testing.T) {
	if Empty() == nil {
		t.Fatal("Empty() returned a nil Node")
	}
}

// TestEmptyInConditionalBranches covers the recommended pattern from the
// Empty docs: returning Empty() from a conditional branch. The old
// True/False constructors called reflect.Value.IsNil on every node,
// which panics for struct values, so this exact usage used to crash.
func TestEmptyInConditionalBranches(t *testing.T) {
	if got := When(true, Empty()).RenderBytes(); len(got) != 0 {
		t.Errorf("When(true, Empty()) rendered %q, want nothing", got)
	}
	if got := Condition(false).False(Empty()).RenderBytes(); len(got) != 0 {
		t.Errorf("Condition(false).False(Empty()) rendered %q, want nothing", got)
	}
}
