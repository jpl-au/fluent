package fieldset

// Byte constants for HTML rendering.
var (
	TagOpen  = []byte("<fieldset")
	TagClose = []byte("</fieldset>")

	AttrDisabled = []byte(" disabled")
	AttrForm     = []byte(" form=\"")
	AttrName     = []byte(" name=\"")
)
