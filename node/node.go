package node

import (
	"bytes"
	"io"
)

// Node represents any renderable item in an HTML document tree.
// This is the base contract that all renderable types satisfy — text, elements,
// function components, and conditionals alike.
type Node interface {
	// Render returns the HTML as a byte slice, or writes it to the provided writer.
	// Use this for top-level rendering where you need the final output.
	Render(w ...io.Writer) []byte

	// RenderBuilder writes HTML into a shared buffer to avoid allocations
	// when composing a tree of nodes. Parent nodes call this on their children.
	RenderBuilder(*bytes.Buffer)

	// Nodes returns the direct children of this node.
	Nodes() []Node
}
