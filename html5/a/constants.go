package a

// Byte constants for HTML rendering.
var (
	TagOpen  = []byte("<a")
	TagClose = []byte("</a>")

	AttrHref           = []byte(" href=\"")
	AttrAttributionSrc = []byte(" attributionsrc=\"")
	AttrDownload       = []byte(" download=\"")
	AttrHrefLang       = []byte(" hreflang=\"")
	AttrPing           = []byte(" ping=\"")
	AttrReferrerPolicy = []byte(" referrerpolicy=\"")
	AttrRel            = []byte(" rel=\"")
	AttrTarget         = []byte(" target=\"")
	AttrType           = []byte(" type=\"")
)
