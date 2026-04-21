// Package node defines the interfaces and helpers that compose a
// fluent render tree. Every renderable item satisfies [Node]: HTML
// elements, text, conditionals, function wrappers, and memoised nodes.
//
// Tree-level helpers create nodes without importing an HTML element
// package: [Func], [Funcs], [Condition], [When], [Unless], and
// [Memoise].
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
	// Render returns the HTML as a byte slice, or writes it to the provided writer.
	// Use this for top-level rendering where you need the final output.
	//
	// Write errors are deliberately discarded rather than returned. A write
	// failure during rendering is almost always a disconnected client or a
	// closed stream - a condition the render tree cannot act on. The real
	// error path lives with the writer's owner (the HTTP handler, the buffer
	// consumer), not inside rendering. Returning an error here would force
	// every node in the tree to handle a failure that isn't its to handle.
	// Implementations silently drop write errors on purpose.
	Render(w ...io.Writer) []byte

	// RenderBuilder writes HTML into a shared buffer to avoid allocations
	// when composing a tree of nodes. Parent nodes call this on their children.
	RenderBuilder(*bytes.Buffer)

	// Nodes returns the children that will be rendered. For conditionals
	// this is the active branch only; for function components this is the
	// evaluated output. Tree walkers rely on this matching Render output.
	Nodes() []Node
}
