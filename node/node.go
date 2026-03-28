// Package node defines the interfaces and helpers that compose a
// fluent render tree. Every renderable item satisfies [Node]: HTML
// elements, text, conditionals, function wrappers, and memoised nodes.
//
// Tree-level helpers create nodes without importing an HTML element
// package: [Func], [Funcs], [Condition], [When], [Unless], and
// [Memoise]. Elements gain reactive tracking via the [Dynamic]
// interface and its .Dynamic(key) method.
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
	Render(w ...io.Writer) []byte

	// RenderBuilder writes HTML into a shared buffer to avoid allocations
	// when composing a tree of nodes. Parent nodes call this on their children.
	RenderBuilder(*bytes.Buffer)

	// Nodes returns the children that will be rendered. For conditionals
	// this is the active branch only; for function components this is the
	// evaluated output. Tree walkers rely on this matching Render output.
	Nodes() []Node
}
