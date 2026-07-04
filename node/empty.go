package node

import (
	"bytes"
	"io"
)

// emptyNode renders nothing. It is stateless, so a single shared value
// (see [empty]) serves every call to [Empty].
type emptyNode struct{}

// empty is the shared instance returned by [Empty]; it carries no state.
var empty Node = emptyNode{}

// Empty returns a node that renders nothing. It is the explicit, always-safe way
// to say "render nothing" - from a [Func] or [Funcs] body, a conditional branch,
// a children slice, or anywhere a [Node] is expected.
//
// You can also return an untyped nil, which fluent's render loops skip. Empty is
// preferred because it is safer and clearer: it is a real node with a no-op
// render, so it cannot trigger a nil dereference the way a typed nil pointer can
// (returning, say, a nil *div.Element), and it reads as a deliberate choice
// rather than a forgotten return.
//
// Usage:
//
//	node.Func(func() node.Node {
//	    if user.LoggedIn {
//	        return greeting(user)
//	    }
//	    return node.Empty()
//	})
func Empty() Node {
	return empty
}

// Render writes nothing.
func (emptyNode) Render(io.Writer) {}

// WriteTo writes nothing and reports zero bytes written.
func (emptyNode) WriteTo(io.Writer) (int64, error) {
	return 0, nil
}

// RenderBytes returns nil, which stringifies to "".
func (emptyNode) RenderBytes() []byte {
	return nil
}

// RenderBuilder writes nothing to the buffer.
func (emptyNode) RenderBuilder(*bytes.Buffer) {}

// Nodes returns no children.
func (emptyNode) Nodes() []Node {
	return nil
}
