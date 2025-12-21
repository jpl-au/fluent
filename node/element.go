package node

import "bytes"

// Element extends Node with methods for rendering opening and closing tags separately.
// This interface enables JIT compilation to pre-render static wrapper tags while
// preserving dynamic content rendering.
type Element interface {
	Node

	// RenderOpen writes the opening tag and attributes to the buffer.
	// For example, for a div with class="container", this writes: <div class="container">
	RenderOpen(buf *bytes.Buffer)

	// RenderClose writes the closing tag to the buffer.
	// For example, for a div, this writes: </div>
	// For self-closing elements, this is a no-op.
	RenderClose(buf *bytes.Buffer)
}
