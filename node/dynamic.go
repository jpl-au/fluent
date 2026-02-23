package node

// Dynamic represents nodes that contain dynamic content requiring re-evaluation on each render.
// IsDynamic reports whether the node produces different output across renders.
// DynamicKey returns the developer-assigned key used by the diff engine to track
// changes across renders. Nodes without a key return an empty string.
type Dynamic interface {
	IsDynamic() bool
	DynamicKey() string
}
