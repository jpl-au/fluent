package svg

// Byte constants for HTML rendering.
var (
	TagOpen  = []byte("<svg")
	TagClose = []byte("</svg>")

	AttrWidth               = []byte(" width=\"")
	AttrHeight              = []byte(" height=\"")
	AttrViewBox             = []byte(" viewBox=\"")
	AttrPreserveAspectRatio = []byte(" preserveAspectRatio=\"")
	AttrXmlns               = []byte(" xmlns=\"")
)
