package node

// The package-level setters below serve code that works through the
// [Element] interface - jit and third-party extensions. The
// chainable SetData and SetAria methods on concrete element types
// cannot appear on the interface (their concrete return type cannot
// satisfy an interface method), which used to force interface-typed
// code to hand-build prefixed keys against the raw escape hatch.
// These functions restore the semantic API - the io.WriteString
// pattern - and route through [Element.SetAttribute], the single
// implementation point every attribute write shares.

// SetAttribute sets a custom attribute on an element reached through
// the [Element] interface. It is the package-level twin of
// [Element.SetAttribute], provided for symmetry with [SetData] and
// [SetAria].
//
// The key is written to the rendered output verbatim. It is code, not
// data: pass a fixed, developer-controlled key. Never build the key from
// user input - a key containing a space, quote, "=", "/" or ">" changes
// the markup structure. The value is escaped; the key is not.
func SetAttribute(e Element, key, value string) {
	e.SetAttribute(key, value)
}

// SetAttributeRaw sets a custom attribute on an element reached through
// the [Element] interface without escaping the value. It is the
// package-level twin of [Element.SetAttributeRaw]; prefer [SetAttribute]
// unless the value is already trusted. The key is written verbatim, as
// on [SetAttribute]; here the value is too.
func SetAttributeRaw(e Element, key, value string) {
	e.SetAttributeRaw(key, value)
}

// SetData sets a data-* attribute on an element reached through the
// [Element] interface, prefixing the key exactly as the concrete
// SetData methods do: SetData(e, "id", v) sets data-id. The prefixed
// key is written verbatim; see [SetAttribute] for the key contract.
func SetData(e Element, key, value string) {
	e.SetAttribute("data-"+key, value)
}

// SetAria sets an aria-* attribute on an element reached through the
// [Element] interface, prefixing the key exactly as the concrete
// SetAria methods do: SetAria(e, "label", v) sets aria-label. The
// prefixed key is written verbatim; see [SetAttribute] for the key
// contract.
func SetAria(e Element, key, value string) {
	e.SetAttribute("aria-"+key, value)
}

// Attribute is a single HTML attribute pair held in an element's
// generic attribute slice. Typed attributes live in dedicated struct
// fields; this type backs the catch-all storage that
// [Element.SetAttribute], SetAria, and SetData write into.
type Attribute struct {
	// Key is the attribute name as it appears in the rendered HTML
	// (for example "class", "data-id", or "hx-get"). It is written
	// verbatim - it is code, not data, and must never be built from
	// user input. A key containing a space, quote, "=", "/" or ">"
	// changes the markup structure.
	Key string

	// Value is the attribute value as it appears in the rendered HTML,
	// written verbatim at render time. SetAttribute escapes before
	// storing, so values placed through it are already safe; a caller
	// building an Attribute literal directly, or using SetAttributeRaw,
	// is responsible for escaping untrusted input itself.
	Value string
}
