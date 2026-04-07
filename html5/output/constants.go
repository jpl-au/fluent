package output

// Byte constants for HTML rendering.
var (
	TagOpen  = []byte("<output")
	TagClose = []byte("</output>")

	AttrFor  = []byte(" for=\"")
	AttrForm = []byte(" form=\"")
	AttrName = []byte(" name=\"")
)
