package node

import (
	"bytes"
	"io"
	"reflect"

	"github.com/jpl-au/fluent"
)

// ConditionalBuilder provides a fluent API for conditional rendering.
// It allows you to specify different content based on a boolean condition.
//
// Nil nodes are safely ignored - if a nil node is provided to True() or False(),
// it will not be stored and nothing will be rendered for that path.
//
// Chain .Dynamic("key") to enable diff engine tracking. Without a key, the
// differ cannot produce targeted patches for conditional content.
//
// Usage:
//
//	node.Condition(user.IsLoggedIn).
//	    True(p.Text("Welcome back!")).
//	    False(p.Text("Please log in")).
//	    Dynamic("auth-message")
type ConditionalBuilder struct {
	condition bool
	trueNode  Node
	falseNode Node
	dynamic   string
}

// Condition creates a new conditional builder with the given boolean condition.
func Condition(condition bool) *ConditionalBuilder {
	return &ConditionalBuilder{
		condition: condition,
	}
}

// When renders the node only when the condition is true.
// This is a shorthand for Condition(cond).True(node).
//
// Usage:
//
//	When(user.IsAdmin, span.Static("Admin"))
func When(condition bool, node Node) *ConditionalBuilder {
	return Condition(condition).True(node)
}

// Unless renders the node only when the condition is false.
// This is a shorthand for Condition(cond).False(node).
//
// Usage:
//
//	Unless(user.IsLoggedIn, a.New().Href("/login").Text("Sign in"))
func Unless(condition bool, node Node) *ConditionalBuilder {
	return Condition(!condition).True(node)
}

// True sets the node to render when the condition is true.
// If node is nil (explicit or typed nil pointer), it is not stored.
func (c *ConditionalBuilder) True(node Node) *ConditionalBuilder {
	if node != nil && !reflect.ValueOf(node).IsNil() {
		c.trueNode = node
	}
	return c
}

// False sets the node to render when the condition is false.
// If node is nil (explicit or typed nil pointer), it is not stored.
func (c *ConditionalBuilder) False(node Node) *ConditionalBuilder {
	if node != nil && !reflect.ValueOf(node).IsNil() {
		c.falseNode = node
	}
	return c
}

// Render writes the active branch's rendered output to w.
// Write errors are intentionally discarded; see [Node] for rationale.
func (c *ConditionalBuilder) Render(w io.Writer) {
	_, _ = c.WriteTo(w)
}

// WriteTo writes the active branch's rendered output to w, returning
// the byte count and any write error. Satisfies [io.WriterTo].
func (c *ConditionalBuilder) WriteTo(w io.Writer) (int64, error) {
	buf := fluent.NewBuffer()
	c.RenderBuilder(buf)
	n, err := buf.WriteTo(w)
	fluent.PutBuffer(buf)
	return n, err
}

// RenderBytes returns the active branch's rendered output as a byte
// slice.
func (c *ConditionalBuilder) RenderBytes() []byte {
	var buf bytes.Buffer
	c.RenderBuilder(&buf)
	return buf.Bytes()
}

// RenderBuilder writes the HTML representation directly to a buffer.
// Renders the appropriate node based on the condition.
func (c *ConditionalBuilder) RenderBuilder(buf *bytes.Buffer) {
	if c.condition && c.trueNode != nil {
		c.trueNode.RenderBuilder(buf)
	} else if !c.condition && c.falseNode != nil {
		c.falseNode.RenderBuilder(buf)
	}
	// If condition doesn't match or no node is set, render nothing
}

// IsDynamic returns true - conditionals always contain dynamic content
// because their output depends on a runtime condition.
func (c *ConditionalBuilder) IsDynamic() bool {
	return true
}

// Dynamic marks this conditional for reactive tracking by the diff engine.
// The key identifies this node across renders so the diff engine can detect
// changes and send targeted patches.
func (c *ConditionalBuilder) Dynamic(key string) *ConditionalBuilder {
	c.dynamic = key
	return c
}

// DynamicKey returns the developer-assigned key for diff engine tracking.
// Returns an empty string if the conditional has not been marked as dynamic.
func (c *ConditionalBuilder) DynamicKey() string {
	return c.dynamic
}

// Nodes returns only the active branch - the one that Render will actually
// produce output for. Returning both branches would mislead tree walkers
// into believing that inactive content exists in the rendered tree, which
// breaks the Differ's structural change detection for keyed elements.
func (c *ConditionalBuilder) Nodes() []Node {
	if c.condition {
		if c.trueNode != nil {
			return []Node{c.trueNode}
		}
	} else {
		if c.falseNode != nil {
			return []Node{c.falseNode}
		}
	}
	return nil
}
