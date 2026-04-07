package meta

// Byte constants for HTML rendering.
var (
	TagOpen = []byte("<meta")

	AttrName      = []byte(" name=\"")
	AttrContent   = []byte(" content=\"")
	AttrCharset   = []byte(" charset=\"")
	AttrHttpEquiv = []byte(" http-equiv=\"")
	AttrScheme    = []byte(" scheme=\"")
	AttrProperty  = []byte(" property=\"")
	AttrMedia     = []byte(" media=\"")
)
