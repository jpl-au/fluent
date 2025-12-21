package node

import (
	"bytes"
	"io"
)

// Node represents any element in an HTML document tree.
// All nodes implement both direct output rendering and efficient buffer rendering
// for optimal performance in different scenarios.
type Node interface {
	// Render generates the complete HTML representation of the node.
	// If a writer is provided, the output is written to it and nil is returned.
	// If no writer is provided, the output is returned as a byte slice.
	Render(w ...io.Writer) []byte

	// RenderBuilder writes the HTML representation directly to a buffer.
	// This method is used for building complex node trees efficiently
	// by allowing nodes to render into a shared buffer.
	RenderBuilder(*bytes.Buffer)

	// Nodes returns a slice of child nodes.
	// For nodes that do not have children, it returns an empty slice.
	Nodes() []Node

	// SetAttribute sets an attribute on the node.
	// This is the primary method for extensions to add attributes safely.
	SetAttribute(key string, value string)
}
