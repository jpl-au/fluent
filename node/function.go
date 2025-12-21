package node

import (
	"bytes"
	"io"

	"github.com/jpl-au/fluent"
)

// FunctionComponent enables dynamic content generation through function calls.
// The function is called during rendering to generate the actual node content.
// This is useful for complex conditional logic, loops, and data transformations.
//
// Nil nodes are safely ignored - if the function returns nil, nothing will be rendered.
//
// Usage:
//
//	Func(func() Node {
//	    if user.IsLoggedIn {
//	        return div.Text("Welcome back!")
//	    }
//	    return div.Text("Please log in")
//	})
type FunctionComponent struct {
	fn func() Node
}

// Func creates a new function component that will call the provided function
// during rendering to generate the actual node content.
func Func(fn func() Node) *FunctionComponent {
	return &FunctionComponent{
		fn: fn,
	}
}

// Render generates the HTML representation by calling the function.
// If a writer is provided, the output is written to it and nil is returned.
// If no writer is provided, the output is returned as a byte slice.
func (f *FunctionComponent) Render(w ...io.Writer) []byte {
	buf := fluent.NewBuffer()
	f.RenderBuilder(buf)

	if len(w) > 0 && w[0] != nil {
		buf.WriteTo(w[0])
		fluent.PutBuffer(buf)
		return nil
	}
	return buf.Bytes()
}

// RenderBuilder writes the HTML representation directly to a buffer.
// Calls the function to get the actual node and renders it.
// Nil nodes are safely ignored.
func (f *FunctionComponent) RenderBuilder(buf *bytes.Buffer) {
	if f.fn != nil {
		node := f.fn()
		if node != nil {
			node.RenderBuilder(buf)
		}
	}
}

// Nodes returns an empty slice as FunctionComponent nodes do not have static children.
func (f *FunctionComponent) Nodes() []Node {
	return []Node{}
}

// SetAttribute is a no-op for FunctionComponent as it does not have attributes.
func (f *FunctionComponent) SetAttribute(key string, value string) {
	// FunctionComponent does not support attributes
}
