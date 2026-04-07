package td

// Byte constants for HTML rendering.
var (
	TagOpen  = []byte("<td")
	TagClose = []byte("</td>")

	AttrColSpan = []byte(" colspan=\"")
	AttrRowSpan = []byte(" rowspan=\"")
	AttrHeaders = []byte(" headers=\"")
)
