package node

import (
	"bytes"
	"io"

	"github.com/jpl-au/fluent"
)

// Memoiser is satisfied by nodes that carry a cache key for the diff
// engine's memoisation layer. When the key matches the previous
// render at the same tree position, the subtree is skipped entirely
// - the closure never runs and no diff is computed for that region.
//
// Regular Diff does not check for this interface. Only the Memoiser
// wrapper in fluent-jit inspects it. When Diff encounters a memo
// node, it calls Render/RenderBuilder/Nodes like any other node
// (the closure executes unconditionally). This makes memo nodes
// transparent to code that does not use memoisation.
type Memoiser interface {
	MemoKey() any
	MemoRender() Node
}

// Memo creates a lazy node with a cache key. The closure produces
// the subtree on demand. When used with a plain Differ, the closure
// is called unconditionally on every render (memo nodes are
// transparent). When used with fluent-jit's Memoiser wrapper, the
// closure is skipped if the key matches the previous render at the
// same tree position.
//
// The key is compared with ==. Use a version counter, a comparable
// struct, or any value where equality means "the subtree has not
// changed". Slices, maps, and functions are not comparable and will
// panic.
//
//	node.Memo(s.ItemsVersion, func() node.Node {
//	    return renderTable(s.Items)
//	})
func Memo(key any, fn func() Node) *MemoNode {
	return &MemoNode{key: key, fn: fn}
}

// MemoNode wraps a lazy closure with a cache key. It satisfies Node
// so it slots into any position in a render tree, and Memoiser so
// the diff engine's memoisation layer can inspect the key.
type MemoNode struct {
	key any
	fn  func() Node
}

// MemoKey returns the cache key for this node. The memoisation layer
// compares this with the previous render's key at the same position.
func (m *MemoNode) MemoKey() any { return m.key }

// MemoRender calls the closure and returns the resulting subtree.
// The memoisation layer calls this only when the key does not match
// the previous render (cache miss).
func (m *MemoNode) MemoRender() Node {
	if m.fn == nil {
		return nil
	}
	return m.fn()
}

// Render calls the closure unconditionally and renders the result.
// This is the path taken by a plain Differ (which does not check
// for Memoiser). The closure always executes - no key checking.
func (m *MemoNode) Render(w ...io.Writer) []byte {
	buf := fluent.NewBuffer()
	m.RenderBuilder(buf)

	if len(w) > 0 && w[0] != nil {
		_, _ = buf.WriteTo(w[0])
		fluent.PutBuffer(buf)
		return nil
	}
	return buf.Bytes()
}

// RenderBuilder calls the closure and writes the result into the
// buffer. Nil closures and nil returns render nothing.
func (m *MemoNode) RenderBuilder(buf *bytes.Buffer) {
	if m.fn == nil {
		return
	}
	if n := m.fn(); n != nil {
		n.RenderBuilder(buf)
	}
}

// Nodes calls the closure and returns its output so tree walkers
// see the same children that Render produces.
func (m *MemoNode) Nodes() []Node {
	if m.fn == nil {
		return nil
	}
	if n := m.fn(); n != nil {
		return []Node{n}
	}
	return nil
}
