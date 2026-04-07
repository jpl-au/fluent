package table

// Byte constants for HTML rendering.
var (
	TagOpen  = []byte("<table")
	TagClose = []byte("</table>")

	AttrBorder      = []byte(" border=\"")
	AttrCellPadding = []byte(" cellpadding=\"")
	AttrCellSpacing = []byte(" cellspacing=\"")
	AttrFrame       = []byte(" frame=\"")
	AttrRules       = []byte(" rules=\"")
	AttrSummary     = []byte(" summary=\"")
	AttrWidth       = []byte(" width=\"")
)
