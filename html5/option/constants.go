package option

// Byte constants for HTML rendering.
var (
	TagOpen  = []byte("<option")
	TagClose = []byte("</option>")

	AttrValue    = []byte(" value=\"")
	AttrDisabled = []byte(" disabled")
	AttrLabel    = []byte(" label=\"")
	AttrSelected = []byte(" selected")
)
