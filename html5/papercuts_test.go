package html5_test

import (
	"strings"
	"testing"

	"github.com/jpl-au/fluent/html5/a"
	"github.com/jpl-au/fluent/html5/div"
	"github.com/jpl-au/fluent/html5/input"
	"github.com/jpl-au/fluent/html5/ol"
	"github.com/jpl-au/fluent/html5/td"
	"github.com/jpl-au/fluent/html5/th"
)

// TestZeroMeaningfulAttributesRender verifies that attributes whose
// zero value is meaningful HTML actually render it. The old generated
// guard was `!= 0`, which silently dropped tabindex="0" (the most
// common tabindex), rowspan="0" (span to the end of the row group),
// maxlength="0" (empty value only) and start="0" (count from zero).
func TestZeroMeaningfulAttributesRender(t *testing.T) {
	cases := []struct {
		name string
		html string
		want string
	}{
		{"tabindex", string(div.New().TabIndex(0).RenderBytes()), ` tabindex="0"`},
		{"rowspan td", string(td.New().RowSpan(0).RenderBytes()), ` rowspan="0"`},
		{"rowspan th", string(th.New().RowSpan(0).RenderBytes()), ` rowspan="0"`},
		{"maxlength", string(input.New().MaxLength(0).RenderBytes()), ` maxlength="0"`},
		{"start", string(ol.New().Start(0).RenderBytes()), ` start="0"`},
	}
	for _, tc := range cases {
		if !strings.Contains(tc.html, tc.want) {
			t.Errorf("%s: zero value should render %q, got %q", tc.name, tc.want, tc.html)
		}
	}
}

// TestDraggableIsEnumerated verifies that draggable renders as the
// enumerated attribute HTML requires. The old bare ` draggable` output
// is invalid and falls back to auto, and draggable="false" matters on
// elements that are draggable by default, like links.
func TestDraggableIsEnumerated(t *testing.T) {
	if got := string(div.New().Draggable(true).RenderBytes()); !strings.Contains(got, ` draggable="true"`) {
		t.Errorf(`Draggable(true) should render draggable="true", got %q`, got)
	}
	if got := string(a.New().Draggable(false).RenderBytes()); !strings.Contains(got, ` draggable="false"`) {
		t.Errorf(`Draggable(false) should render draggable="false", got %q`, got)
	}
}
