package node

// Dynamic represents nodes that contain dynamic content requiring re-evaluation on each render.
// This interface is used by the JIT compiler to identify nodes that cannot be pre-rendered.
type Dynamic interface {
	// Dynamic returns true if this node contains dynamic content that must be re-evaluated on each render.
	Dynamic() bool
}