// Package node defines the interfaces and helpers that compose a
// fluent render tree. Every renderable item satisfies [Node]: HTML
// elements, text, conditionals, function wrappers, and memoised nodes.
//
// Tree-level helpers create nodes without importing an HTML element
// package: [Func], [Funcs], [Condition], [When], [Unless], and [Empty]
// (a node that renders nothing). Memoisation nodes live in fluent-jit,
// alongside the diff engine that consumes them.
//
// All node types that produce dynamic content support reactive tracking
// via the [Dynamic] interface. Call .Dynamic(key) on elements, function
// components, and conditionals to assign a tracking key. The diff
// engine in fluent-jit (and Tether's reactive UI) uses these keys to
// produce targeted patches when content changes between renders.
package node

import (
	"bytes"
	"io"
)

// Node represents any renderable item in an HTML document tree.
// This is the base contract that all renderable types satisfy - text, elements,
// function components, and conditionals alike.
type Node interface {
	// Render writes the HTML to w. Fire-and-forget for HTTP handlers.
	//
	// Write errors are deliberately discarded rather than returned. A write
	// failure during rendering is almost always a disconnected client or a
	// closed stream - a condition the render tree cannot act on. The real
	// error path lives with the writer's owner (the HTTP handler, the buffer
	// consumer), not inside rendering. Choosing Render over WriteTo is
	// choosing to drop that error on purpose; use WriteTo when it matters.
	Render(w io.Writer)

	// WriteTo renders the HTML and writes it to w, returning the byte
	// count and any write error. It satisfies [io.WriterTo], so a node
	// drops straight into io.Copy and anything else that accepts one.
	// This is the one place a render's write error surfaces.
	WriteTo(w io.Writer) (int64, error)

	// RenderBytes returns the HTML as a byte slice. Use it where no
	// writer is involved: caching, snapshots, tests.
	RenderBytes() []byte

	// RenderBuilder writes HTML into a shared buffer to avoid allocations
	// when composing a tree of nodes. Parent nodes call this on their children.
	RenderBuilder(*bytes.Buffer)

	// Nodes returns the children that will be rendered. For conditionals
	// this is the active branch only; for function components this is the
	// evaluated output. Tree walkers rely on this matching Render output.
	Nodes() []Node
}
