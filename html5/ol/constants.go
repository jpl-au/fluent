package ol

// Byte constants for HTML rendering.
var (
	TagOpen  = []byte("<ol")
	TagClose = []byte("</ol>")

	AttrReversed = []byte(" reversed")
	AttrStart    = []byte(" start=\"")
	AttrType     = []byte(" type=\"")
)
