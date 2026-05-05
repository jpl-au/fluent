package node

// Attribute is a single HTML attribute pair held in an element's
// generic attribute slice. Typed attributes live in dedicated struct
// fields; this type backs the catch-all storage that
// [Element.SetAttribute], SetAria, and SetData write into.
type Attribute struct {
	// Key is the attribute name as it appears in the rendered HTML
	// (for example "class", "data-id", or "hx-get"). It is written
	// verbatim - callers are expected to pass a valid attribute name.
	Key string

	// Value is the attribute value as it appears in the rendered HTML.
	// It is written verbatim and is not HTML-escaped, so callers that
	// accept untrusted input must escape before storing.
	Value string
}
