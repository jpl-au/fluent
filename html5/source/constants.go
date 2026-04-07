package source

// Byte constants for HTML rendering.
var (
	TagOpen = []byte("<source")

	AttrSrc    = []byte(" src=\"")
	AttrType   = []byte(" type=\"")
	AttrSrcset = []byte(" srcset=\"")
	AttrSizes  = []byte(" sizes=\"")
	AttrMedia  = []byte(" media=\"")
	AttrWidth  = []byte(" width=\"")
	AttrHeight = []byte(" height=\"")
)
