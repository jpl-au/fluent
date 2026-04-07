package progress

// Byte constants for HTML rendering.
var (
	TagOpen  = []byte("<progress")
	TagClose = []byte("</progress>")

	AttrValue = []byte(" value=\"")
	AttrMax   = []byte(" max=\"")
)
