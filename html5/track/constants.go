package track

// Byte constants for HTML rendering.
var (
	TagOpen = []byte("<track")

	AttrSrc     = []byte(" src=\"")
	AttrKind    = []byte(" kind=\"")
	AttrLabel   = []byte(" label=\"")
	AttrSrclang = []byte(" srclang=\"")
	AttrDefault = []byte(" default")
)
