package form

// Byte constants for HTML rendering.
var (
	TagOpen  = []byte("<form")
	TagClose = []byte("</form>")

	AttrAction         = []byte(" action=\"")
	AttrMethod         = []byte(" method=\"")
	AttrAcceptCharset  = []byte(" accept-charset=\"")
	AttrAutoCapitalize = []byte(" autocapitalize=\"")
	AttrAutoComplete   = []byte(" autocomplete=\"")
	AttrEncType        = []byte(" enctype=\"")
	AttrName           = []byte(" name=\"")
	AttrNoValidate     = []byte(" novalidate")
	AttrRel            = []byte(" rel=\"")
	AttrTarget         = []byte(" target=\"")
)
