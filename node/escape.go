package node

import "strings"

// EscapeAttribute returns s with the characters that could break out of a
// double-quoted HTML attribute value replaced by entities. Fluent always
// renders attribute values inside double quotes, so escaping the quote is
// what prevents an attacker-supplied value from closing the attribute and
// injecting new attributes or event handlers.
//
// It escapes the same six characters as the standard library's attribute
// escaper - " ' & < > and NUL - using the short numeric entities, so its
// output is byte-identical to html/template and gomponents for the same
// input. The clean-value fast path allocates nothing: a value with none of
// these characters (the overwhelming common case) is returned unchanged.
//
// It deliberately does not escape '+': the UTF-7 charset-sniffing attack that
// motivated escaping it died with legacy browsers, and escaping it would put
// fluent's output out of parity with gomponents and the benchmark suite.
//
// Escaping happens at set time, once, when a value is stored - not at render
// time - so the render hot path is untouched and the JIT bakes escaped static
// attributes into compiled fragments.
func EscapeAttribute(s string) string {
	if !strings.ContainsAny(s, "\"'&<>\x00") {
		return s
	}

	var b strings.Builder
	b.Grow(len(s) + 8)
	last := 0
	for i := 0; i < len(s); i++ {
		var repl string
		switch s[i] {
		case '"':
			repl = "&#34;"
		case '\'':
			repl = "&#39;"
		case '&':
			repl = "&amp;"
		case '<':
			repl = "&lt;"
		case '>':
			repl = "&gt;"
		case 0:
			repl = "�"
		default:
			continue
		}
		b.WriteString(s[last:i])
		b.WriteString(repl)
		last = i + 1
	}
	b.WriteString(s[last:])
	return b.String()
}
