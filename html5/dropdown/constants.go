package dropdown

// Byte constants for HTML rendering.
var (
	TagOpen  = []byte("<select")
	TagClose = []byte("</select>")

	AttrName         = []byte(" name=\"")
	AttrAutoComplete = []byte(" autocomplete=\"")
	AttrDisabled     = []byte(" disabled")
	AttrForm         = []byte(" form=\"")
	AttrMultiple     = []byte(" multiple")
	AttrRequired     = []byte(" required")
	AttrSize         = []byte(" size=\"")
)
