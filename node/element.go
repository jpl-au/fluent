package node

import "bytes"

// Element extends Node for HTML elements that have attributes and an open/close tag structure.
// Not all nodes are elements — text nodes, function components, and conditionals are not.
// This separation allows extensions (htmx, turbo, shoelace) to accept only types that
// genuinely support attributes, and allows JIT compilation to pre-render static wrapper
// tags independently of dynamic content.
type Element interface {
	Node

	SetAttribute(key string, value string)

	// RenderOpen and RenderClose split the element's rendering so that JIT can
	// cache the opening tag separately from the children.
	// For example: RenderOpen writes <div class="container">, RenderClose writes </div>.
	RenderOpen(buf *bytes.Buffer)
	RenderClose(buf *bytes.Buffer)
}
