package node

// Dynamic represents nodes that contain dynamic content requiring
// re-evaluation on each render. The fluent-jit diff engine uses this
// interface to identify trackable nodes in the render tree.
//
// IsDynamic reports whether the node produces different output across
// renders. DynamicKey returns the developer-assigned key used by the
// diff engine to track changes. Nodes without a key return an empty
// string and are not individually tracked.
//
// HTML elements, [FunctionComponent], [FuncsComponent], and
// [ConditionalBuilder] all satisfy this interface. Call .Dynamic(key)
// on any of them to assign a tracking key.
type Dynamic interface {
	IsDynamic() bool
	DynamicKey() string
}
