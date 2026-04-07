package template

// Byte constants for HTML rendering.
var (
	TagOpen  = []byte("<template")
	TagClose = []byte("</template>")

	AttrShadowRootMode           = []byte(" shadowrootmode=\"")
	AttrShadowRootClonable       = []byte(" shadowrootclonable")
	AttrShadowRootDelegatesFocus = []byte(" shadowrootdelegatesfocus")
	AttrShadowRootSerializable   = []byte(" shadowrootserializable")
)
