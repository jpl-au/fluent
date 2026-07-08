package security_test

import (
	"strings"
	"testing"

	"github.com/jpl-au/fluent/html5/a"
	"github.com/jpl-au/fluent/html5/attr/target"
	"github.com/jpl-au/fluent/html5/div"
	"github.com/jpl-au/fluent/html5/iframe"
	"github.com/jpl-au/fluent/html5/img"
	"github.com/jpl-au/fluent/html5/svg"
	"github.com/jpl-au/fluent/node"
)

// hostile closes a double-quoted attribute and adds an event handler; if any path
// stores it raw, the rendered tag gains a live onmouseover handler.
const hostile = `x" onmouseover="alert(1)`

// escaped is hostile after node.EscapeAttribute: the double quotes become &#34;,
// so the value can no longer break out of the attribute.
const escaped = `x&#34; onmouseover=&#34;alert(1)`

func render(n interface{ RenderBytes() []byte }) string { return string(n.RenderBytes()) }

// TestSettersEscape proves the typed setters, SetAttribute, Dynamic keys and enum
// Custom all neutralise a breakout string.
func TestSettersEscape(t *testing.T) {
	setAttr := div.New()
	setAttr.SetAttribute("data-x", hostile)
	cases := []struct {
		name string
		got  string
	}{
		{"Class", render(div.New().Class(hostile))},
		{"SetAttribute", render(setAttr)},
		{"Dynamic", render(div.New().Dynamic(hostile))},
		{"enum Custom", render(a.New().Target(target.Custom(hostile)))},
	}
	for _, c := range cases {
		if strings.Contains(c.got, `"`+" onmouseover") || !strings.Contains(c.got, escaped) {
			t.Errorf("%s did not escape the breakout: %q", c.name, c.got)
		}
	}
}

// TestSetAttributeRaw is the escape hatch: it must NOT escape, so the breakout
// survives verbatim. This is the deliberate trusted-value path.
func TestSetAttributeRaw(t *testing.T) {
	e := div.New()
	e.SetAttributeRaw("data-x", hostile)
	got := render(e)
	if !strings.Contains(got, hostile) {
		t.Errorf("SetAttributeRaw should store the value verbatim, got %q", got)
	}
}

// TestURLFilterRejects proves navigation sinks reject off-allowlist schemes with
// the inert sentinel, and that iframe (a plain nav sink) rejects data: too.
func TestURLFilterRejects(t *testing.T) {
	cases := []struct {
		name string
		got  string
	}{
		{"a.Href javascript:", render(a.New().Href("javascript:alert(1)"))},
		{"a.Href vbscript:", render(a.New().Href("vbscript:msgbox(1)"))},
		{"iframe.Src data:", render(iframe.New().Src("data:text/html,<script>alert(1)</script>"))},
	}
	for _, c := range cases {
		if !strings.Contains(c.got, node.UnsafeURL) {
			t.Errorf("%s should be filtered to the sentinel, got %q", c.name, c.got)
		}
	}
}

// TestURLFilterAllows proves clean URLs pass unchanged, and that image sinks
// (escaped only, never filtered) keep their data: URL.
func TestURLFilterAllows(t *testing.T) {
	if got := render(a.New().Href("https://example.com/x?a=b")); strings.Contains(got, node.UnsafeURL) {
		t.Errorf("clean https URL was filtered: %q", got)
	}
	// img.Src is an image sink: escaped, never filtered, so a data:image URL renders.
	if got := render(img.New().Src("data:image/png;base64,AAAA")); strings.Contains(got, node.UnsafeURL) {
		t.Errorf("img data:image URL should not be filtered: %q", got)
	}
}

// TestSVGEscapes proves the svg shapes share the escaping (svg has no byte golden,
// so it needs an explicit probe).
func TestSVGEscapes(t *testing.T) {
	svgSet := svg.Rect()
	svgSet.SetAttribute("x", hostile)
	cases := []string{
		render(svg.Rect().Class(hostile)),
		render(svgSet),
	}
	for _, got := range cases {
		if strings.Contains(got, `"`+" onmouseover") || !strings.Contains(got, escaped) {
			t.Errorf("svg did not escape the breakout: %q", got)
		}
	}
}
