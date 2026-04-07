package html

// Byte constants for HTML rendering.
var (
	TagOpen  = []byte("<html")
	TagClose = []byte("</html>")

	AttrLang     = []byte(" lang=\"")
	AttrXmlns    = []byte(" xmlns=\"")
	AttrManifest = []byte(" manifest=\"")
)
