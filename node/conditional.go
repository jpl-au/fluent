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
// Usage:
//
//	Condition(user.IsLoggedIn).
//	    True(p.Text("Welcome back!")).
//	    False(p.Text("Please log in"))
type ConditionalBuilder struct {
	condition bool
	trueNode  Node
	falseNode Node
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

// Render generates the HTML representation based on the condition.
// If a writer is provided, the output is written to it and nil is returned.
// If no writer is provided, the output is returned as a byte slice.
func (c *ConditionalBuilder) Render(w ...io.Writer) []byte {
	buf := fluent.NewBuffer()
	c.RenderBuilder(buf)

	if len(w) > 0 && w[0] != nil {
		buf.WriteTo(w[0])
		fluent.PutBuffer(buf)
		return nil
	}
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

// Nodes returns the potential child nodes of the ConditionalBuilder.
func (c *ConditionalBuilder) Nodes() []Node {
	children := []Node{}
	if c.trueNode != nil {
		children = append(children, c.trueNode)
	}
	if c.falseNode != nil {
		children = append(children, c.falseNode)
	}
	return children
}

// SetAttribute is a no-op for ConditionalBuilder as it does not have attributes.
func (c *ConditionalBuilder) SetAttribute(key string, value string) {
	// ConditionalBuilder does not support attributes
}
