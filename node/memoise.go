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
// Regular Diff does not check for this interface. Only fluent-jit's
// [jit.Memoiser] struct inspects it. When Diff encounters a memoised
// node, it calls Render/RenderBuilder/Nodes like any other node
// (the closure executes unconditionally). This makes memoised nodes
// transparent to code that does not use memoisation.
type Memoiser interface {
	MemoiseKey() any
	MemoiseRender() Node
}

// SharedMemoiser is a [Memoiser] whose rendered bytes may be reused
// across sessions, not just across renders of one session. When the
// diff engine's memoisation layer meets a shared node on a cache miss,
// it looks the key up in a process-global store: the first session to
// render a given key populates it, and every other session with the
// same key is served those bytes instead of running the closure. On a
// busy broadcast - a shared header, a live leaderboard - the render
// runs once for the whole process rather than once per session.
//
// The contract is stricter than plain memoisation: the key must be
// globally unique and must fully determine the rendered bytes. A plain
// [Memoise] key is only compared within one session at one tree
// position, so a bare counter is safe; a shared key is compared across
// every session in the process, so it must be namespaced and derived
// from the content itself (a version or hash), never from per-session
// state. [Shared] enforces none of this - correctness is the caller's.
type SharedMemoiser interface {
	Memoiser
	MemoiseShared() bool
}

// Memoise creates a lazy node with a cache key. The closure produces
// the subtree on demand. When used with a plain Differ, the closure
// is called unconditionally on every render (memoised nodes are
// transparent). When used with fluent-jit's Memoiser wrapper, the
// closure is skipped if the key matches the previous render at the
// same tree position.
//
// The key is compared with ==. Use a version counter, a comparable
// struct, or any value where equality means "the subtree has not
// changed". Slices, maps, and functions are not comparable and will
// panic.
//
//	node.Memoise(s.ItemsVersion, func() node.Node {
//	    return renderTable(s.Items)
//	})
func Memoise(key any, fn func() Node) *MemoisedNode {
	return &MemoisedNode{key: key, fn: fn}
}

// Shared is like [Memoise] but marks the region as reusable across
// sessions: the diff engine caches its rendered bytes in a process-
// global store keyed by key, so the closure runs at most once per
// distinct key for the whole process rather than once per session.
// Use it for regions that render identically for every user - a shared
// header, a navigation bar, a live scoreboard broadcast to a room.
//
// The key MUST be globally unique and MUST fully determine the rendered
// bytes. Namespace it and derive it from the content ("nav:v3",
// "board:" + hash), never from per-session state - two sessions with
// the same key are served the same bytes. See [SharedMemoiser] for the
// full contract.
//
//	node.Shared("leaderboard:"+s.BoardVersion, func() node.Node {
//	    return renderBoard(s.Board)
//	})
func Shared(key any, fn func() Node) *MemoisedNode {
	return &MemoisedNode{key: key, fn: fn, shared: true}
}

// MemoisedNode wraps a lazy closure with a cache key. It satisfies
// [Node] so it slots into any position in a render tree, and
// [Memoiser] so the diff engine's memoisation layer can inspect the
// key.
type MemoisedNode struct {
	key    any
	fn     func() Node
	shared bool
}

// MemoiseKey returns the cache key for this node. The memoisation layer
// compares this with the previous render's key at the same position.
func (m *MemoisedNode) MemoiseKey() any { return m.key }

// MemoiseShared reports whether this node's rendered bytes may be
// reused across sessions via the process-global store. True only for
// nodes created with [Shared]. See [SharedMemoiser].
func (m *MemoisedNode) MemoiseShared() bool { return m.shared }

// MemoiseRender calls the closure and returns the resulting subtree.
// The memoisation layer calls this only when the key does not match
// the previous render (cache miss).
func (m *MemoisedNode) MemoiseRender() Node {
	if m.fn == nil {
		return nil
	}
	return m.fn()
}

// Render calls the closure unconditionally and renders the result.
// This is the path taken by a plain Differ (which does not check
// for Memoiser). The closure always executes - no key checking.
func (m *MemoisedNode) Render(w ...io.Writer) []byte {
	buf := fluent.NewBuffer()
	m.RenderBuilder(buf)

	if len(w) > 0 && w[0] != nil {
		// Write errors are intentionally discarded; see [node.Node] for rationale.
		_, _ = buf.WriteTo(w[0])
		fluent.PutBuffer(buf)
		return nil
	}
	return buf.Bytes()
}

// RenderBuilder calls the closure and writes the result into the
// buffer. Nil closures and nil returns render nothing.
func (m *MemoisedNode) RenderBuilder(buf *bytes.Buffer) {
	if m.fn == nil {
		return
	}
	if n := m.fn(); n != nil {
		n.RenderBuilder(buf)
	}
}

// Nodes calls the closure and returns its output so tree walkers
// see the same children that Render produces.
func (m *MemoisedNode) Nodes() []Node {
	if m.fn == nil {
		return nil
	}
	if n := m.fn(); n != nil {
		return []Node{n}
	}
	return nil
}
