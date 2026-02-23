package node

// Dynamic represents nodes that contain dynamic content requiring re-evaluation on each render.
// This interface is used by the JIT compiler to identify nodes that cannot be pre-rendered.
type Dynamic interface {
	IsDynamic() bool
}
