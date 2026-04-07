package textarea

// Byte constants for HTML rendering.
var (
	TagOpen  = []byte("<textarea")
	TagClose = []byte("</textarea>")

	AttrName         = []byte(" name=\"")
	AttrRows         = []byte(" rows=\"")
	AttrCols         = []byte(" cols=\"")
	AttrAutoComplete = []byte(" autocomplete=\"")
	AttrDisabled     = []byte(" disabled")
	AttrForm         = []byte(" form=\"")
	AttrMaxLength    = []byte(" maxlength=\"")
	AttrMinLength    = []byte(" minlength=\"")
	AttrPlaceholder  = []byte(" placeholder=\"")
	AttrReadOnly     = []byte(" readonly")
	AttrRequired     = []byte(" required")
	AttrSpellCheck   = []byte(" spellcheck=\"")
	AttrWrap         = []byte(" wrap=\"")
	AttrDirName      = []byte(" dirname=\"")
)
