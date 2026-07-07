package node

import "bytes"

// Element extends Node for HTML elements that have attributes and an open/close tag structure.
// Not all nodes are elements - text nodes, function components, and conditionals are not.
// This separation allows extensions (htmx, turbo, shoelace) to accept only types that
// genuinely support attributes, and allows JIT compilation to pre-render static wrapper
// tags independently of dynamic content.
type Element interface {
	Node

	// SetAttribute is the escape hatch for attributes that Fluent's typed
	// API does not cover - framework directives like Alpine.js (x-on:click)
	// or HTMX (hx-get), and any other custom or non-standard attribute.
	//
	// For everything else, prefer the typed methods on each element
	// (.Class, .Href, .Checked, ...). They produce the right escaping,
	// validate enumerated values at compile time, and chain. SetAttribute
	// returns nothing on purpose - it cannot be chained, which keeps the
	// Fluent API the obvious default.
	//
	// On concrete element types, prefer SetAria(name, value) and
	// SetData(name, value) for ARIA and data-* attributes - they prefix
	// the key and return the element for chaining. Those chainable
	// methods cannot appear on this interface (their concrete return
	// type cannot satisfy an interface method), so code holding an
	// Element uses the package-level [SetAttribute], [SetData] and
	// [SetAria] functions instead.
	SetAttribute(key string, value string)

	// RenderOpen writes the opening tag, including the element name and
	// every attribute, into buf. JIT compilation caches this output
	// separately from the children so static wrappers can be pre-rendered
	// while dynamic content re-evaluates each render. For <div class="x">
	// it writes the entire opening fragment.
	RenderOpen(buf *bytes.Buffer)

	// RenderClose writes the closing tag (or the self-closing terminator
	// for void elements) into buf. Paired with RenderOpen so the JIT can
	// cache the open and close fragments independently of the children
	// rendered between them.
	RenderClose(buf *bytes.Buffer)
}
