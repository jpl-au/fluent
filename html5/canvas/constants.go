package canvas

// Byte constants for HTML rendering.
var (
	TagOpen  = []byte("<canvas")
	TagClose = []byte("</canvas>")

	AttrHeight = []byte(" height=\"")
	AttrWidth  = []byte(" width=\"")
)
