package script

// Byte constants for HTML rendering.
var (
	TagOpen  = []byte("<script")
	TagClose = []byte("</script>")

	AttrSrc            = []byte(" src=\"")
	AttrType           = []byte(" type=\"")
	AttrAsync          = []byte(" async")
	AttrCrossOrigin    = []byte(" crossorigin=\"")
	AttrDefer          = []byte(" defer")
	AttrIntegrity      = []byte(" integrity=\"")
	AttrNoModule       = []byte(" nomodule")
	AttrNonce          = []byte(" nonce=\"")
	AttrReferrerPolicy = []byte(" referrerpolicy=\"")
	AttrBlocking       = []byte(" blocking=\"")
	AttrFetchPriority  = []byte(" fetchpriority=\"")
)
