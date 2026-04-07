package dialog

// Byte constants for HTML rendering.
var (
	TagOpen  = []byte("<dialog")
	TagClose = []byte("</dialog>")

	AttrOpen     = []byte(" open")
	AttrClosedBy = []byte(" closedby=\"")
)
