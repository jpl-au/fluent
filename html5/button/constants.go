package button

// Byte constants for HTML rendering.
var (
	TagOpen  = []byte("<button")
	TagClose = []byte("</button>")

	AttrCommand             = []byte(" command=\"")
	AttrCommandfor          = []byte(" commandfor=\"")
	AttrDisabled            = []byte(" disabled")
	AttrForm                = []byte(" form=\"")
	AttrFormaction          = []byte(" formaction=\"")
	AttrFormenctype         = []byte(" formenctype=\"")
	AttrFormmethod          = []byte(" formmethod=\"")
	AttrFormnovalidate      = []byte(" formnovalidate")
	AttrFormtarget          = []byte(" formtarget=\"")
	AttrName                = []byte(" name=\"")
	AttrPopOverTarget       = []byte(" popovertarget=\"")
	AttrPopovertargetaction = []byte(" popovertargetaction=\"")
	AttrType                = []byte(" type=\"")
	AttrValue               = []byte(" value=\"")
)
