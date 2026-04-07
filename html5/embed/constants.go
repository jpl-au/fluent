package embed

// Byte constants for HTML rendering.
var (
	TagOpen = []byte("<embed")

	AttrSrc    = []byte(" src=\"")
	AttrType   = []byte(" type=\"")
	AttrWidth  = []byte(" width=\"")
	AttrHeight = []byte(" height=\"")
)
