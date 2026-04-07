package object

// Byte constants for HTML rendering.
var (
	TagOpen  = []byte("<object")
	TagClose = []byte("</object>")

	AttrData   = []byte(" data=\"")
	AttrType   = []byte(" type=\"")
	AttrForm   = []byte(" form=\"")
	AttrHeight = []byte(" height=\"")
	AttrName   = []byte(" name=\"")
	AttrUseMap = []byte(" usemap=\"")
	AttrWidth  = []byte(" width=\"")
)
