package img

// Byte constants for HTML rendering.
var (
	TagOpen = []byte("<img")

	AttrSrc            = []byte(" src=\"")
	AttrAlt            = []byte(" alt=\"")
	AttrWidth          = []byte(" width=\"")
	AttrHeight         = []byte(" height=\"")
	AttrLoading        = []byte(" loading=\"")
	AttrSizes          = []byte(" sizes=\"")
	AttrSrcset         = []byte(" srcset=\"")
	AttrCrossOrigin    = []byte(" crossorigin=\"")
	AttrDecoding       = []byte(" decoding=\"")
	AttrFetchPriority  = []byte(" fetchpriority=\"")
	AttrReferrerPolicy = []byte(" referrerpolicy=\"")
	AttrIsMap          = []byte(" ismap")
	AttrUseMap         = []byte(" usemap=\"")
	AttrAttributionSrc = []byte(" attributionsrc=\"")
	AttrElementTiming  = []byte(" elementtiming=\"")
)
