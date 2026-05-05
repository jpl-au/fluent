// Package html5 holds the cross-element machinery that the per-element
// packages under html5/ build on. It is generated from the HTML5
// specification YAML and exports three families of declarations:
//
//   - The byte-slice Tag, Markup, and Attr constants in constants.go
//     are internal rendering primitives used by every element's
//     RenderOpen and RenderClose to avoid string-to-byte conversions
//     when writing into a bytes.Buffer. They are exported so that
//     extensions and custom elements can reuse the same write path,
//     not as a way to set attributes - use the typed methods on each
//     element (.Class, .Href, .Checked, ...) or
//     [github.com/jpl-au/fluent/node.Element.SetAttribute] for that.
//
//   - [GlobalAttributes] embeds the attributes available on every
//     HTML element (style, title, lang, role, ARIA, popover, ...).
//     Each per-element struct embeds it.
//
//   - [EventAttributes] embeds the on* event handler attributes
//     available on every HTML element (onclick, onchange, ...).
//
// Application code rarely imports html5 directly - reach for the
// per-element packages (html5/div, html5/input, html5/a, ...) and the
// typed attribute packages under html5/attr instead.
package html5
