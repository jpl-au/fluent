package node

import (
	"bytes"
	"io"

	"github.com/jpl-au/fluent"
)

// FunctionsComponent enables dynamic content generation of multiple nodes.
// The function is called during rendering to generate the actual node content.
// This is useful for generating lists of items, e.g. from a loop.
//
// Nil nodes are safely ignored - any nil nodes in the returned slice will be skipped.
//
// Usage:
//
//	FuncNodes(func() []Node {
//	    nodes := []Node{}
//	    for _, item := range items {
//	        nodes = append(nodes, li.Text(item.Name))
//	    }
//	    return nodes
//	})
type FunctionsComponent struct {
	fn func() []Node
}

// FuncNodes creates a new function component that will call the provided function
// during rendering to generate a slice of nodes.
func FuncNodes(fn func() []Node) *FunctionsComponent {
	return &FunctionsComponent{
		fn: fn,
	}
}

// Render generates the HTML representation by calling the function.
// If a writer is provided, the output is written to it and nil is returned.
// If no writer is provided, the output is returned as a byte slice.
func (f *FunctionsComponent) Render(w ...io.Writer) []byte {
	buf := fluent.NewBuffer()
	f.RenderBuilder(buf)

	if len(w) > 0 && w[0] != nil {
		_, _ = buf.WriteTo(w[0])
		fluent.PutBuffer(buf)
		return nil
	}
	return buf.Bytes()
}

// RenderBuilder writes the HTML representation directly to a buffer.
// Calls the function to get the actual nodes and renders them.
// Nil nodes are safely ignored.
func (f *FunctionsComponent) RenderBuilder(buf *bytes.Buffer) {
	if f.fn != nil {
		nodes := f.fn()
		for _, node := range nodes {
			if node != nil {
				node.RenderBuilder(buf)
			}
		}
	}
}

// Nodes returns an empty slice as FunctionsComponent nodes do not have static children.
func (f *FunctionsComponent) Nodes() []Node {
	return []Node{}
}

// SetAttribute is a no-op for FunctionsComponent as it does not have attributes.
func (f *FunctionsComponent) SetAttribute(_ string, _ string) {
	// FunctionsComponent does not support attributes
}
