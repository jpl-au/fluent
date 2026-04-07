package th

// Byte constants for HTML rendering.
var (
	TagOpen  = []byte("<th")
	TagClose = []byte("</th>")

	AttrColSpan = []byte(" colspan=\"")
	AttrRowSpan = []byte(" rowspan=\"")
	AttrHeaders = []byte(" headers=\"")
	AttrScope   = []byte(" scope=\"")
	AttrAbbr    = []byte(" abbr=\"")
)
