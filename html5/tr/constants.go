package tr

// Byte constants for HTML rendering.
var (
	TagOpen  = []byte("<tr")
	TagClose = []byte("</tr>")

	AttrAlign   = []byte(" align=\"")
	AttrBgColor = []byte(" bgcolor=\"")
	AttrChar    = []byte(" char=\"")
	AttrCharOff = []byte(" charoff=\"")
	AttrVAlign  = []byte(" valign=\"")
)
