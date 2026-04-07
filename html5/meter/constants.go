package meter

// Byte constants for HTML rendering.
var (
	TagOpen  = []byte("<meter")
	TagClose = []byte("</meter>")

	AttrValue   = []byte(" value=\"")
	AttrMin     = []byte(" min=\"")
	AttrMax     = []byte(" max=\"")
	AttrLow     = []byte(" low=\"")
	AttrHigh    = []byte(" high=\"")
	AttrOptimum = []byte(" optimum=\"")
)
