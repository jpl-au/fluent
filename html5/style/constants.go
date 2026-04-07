package style

// Byte constants for HTML rendering.
var (
	TagOpen  = []byte("<style")
	TagClose = []byte("</style>")

	AttrType     = []byte(" type=\"")
	AttrMedia    = []byte(" media=\"")
	AttrNonce    = []byte(" nonce=\"")
	AttrTitle    = []byte(" title=\"")
	AttrBlocking = []byte(" blocking=\"")
)
