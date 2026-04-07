package iframe

// Byte constants for HTML rendering.
var (
	TagOpen  = []byte("<iframe")
	TagClose = []byte("</iframe>")

	AttrSrc             = []byte(" src=\"")
	AttrWidth           = []byte(" width=\"")
	AttrHeight          = []byte(" height=\"")
	AttrLoading         = []byte(" loading=\"")
	AttrAllow           = []byte(" allow=\"")
	AttrAllowFullscreen = []byte(" allowfullscreen")
	AttrName            = []byte(" name=\"")
	AttrReferrerPolicy  = []byte(" referrerpolicy=\"")
	AttrSandbox         = []byte(" sandbox=\"")
	AttrSrcDoc          = []byte(" srcdoc=\"")
	AttrCsp             = []byte(" csp=\"")
	AttrCredentialless  = []byte(" credentialless")
)
