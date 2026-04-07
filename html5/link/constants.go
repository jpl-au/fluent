package link

// Byte constants for HTML rendering.
var (
	TagOpen = []byte("<link")

	AttrRel            = []byte(" rel=\"")
	AttrHref           = []byte(" href=\"")
	AttrAs             = []byte(" as=\"")
	AttrCrossOrigin    = []byte(" crossorigin=\"")
	AttrDisabled       = []byte(" disabled")
	AttrHrefLang       = []byte(" hreflang=\"")
	AttrIntegrity      = []byte(" integrity=\"")
	AttrMedia          = []byte(" media=\"")
	AttrReferrerPolicy = []byte(" referrerpolicy=\"")
	AttrSizes          = []byte(" sizes=\"")
	AttrTitle          = []byte(" title=\"")
	AttrType           = []byte(" type=\"")
	AttrFetchPriority  = []byte(" fetchpriority=\"")
	AttrBlocking       = []byte(" blocking=\"")
	AttrImageSrcset    = []byte(" imagesrcset=\"")
	AttrImageSizes     = []byte(" imagesizes=\"")
	AttrColor          = []byte(" color=\"")
)
