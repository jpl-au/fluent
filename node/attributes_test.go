package node_test

import (
	"strings"
	"testing"

	"github.com/jpl-au/fluent/html5/div"
	"github.com/jpl-au/fluent/node"
)

// The package-level setters must behave identically to the concrete
// chainable methods when the same element is reached through the
// Element interface - same prefixing, same rendered output.
func TestPackageLevelSetters(t *testing.T) {
	concrete := div.New()
	concrete.SetData("id", "42").SetAria("label", "answer")
	concrete.SetAttribute("hx-get", "/next")

	var viaInterface node.Element = div.New()
	node.SetData(viaInterface, "id", "42")
	node.SetAria(viaInterface, "label", "answer")
	node.SetAttribute(viaInterface, "hx-get", "/next")

	want := string(concrete.RenderBytes())
	got := string(viaInterface.RenderBytes())
	if got != want {
		t.Errorf("interface-typed setters rendered %q, concrete methods rendered %q", got, want)
	}
	for _, attr := range []string{`data-id="42"`, `aria-label="answer"`, `hx-get="/next"`} {
		if !strings.Contains(got, attr) {
			t.Errorf("rendered output missing %s: %q", attr, got)
		}
	}
}
