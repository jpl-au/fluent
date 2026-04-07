package optgroup

// Byte constants for HTML rendering.
var (
	TagOpen  = []byte("<optgroup")
	TagClose = []byte("</optgroup>")

	AttrLabel    = []byte(" label=\"")
	AttrDisabled = []byte(" disabled")
)
