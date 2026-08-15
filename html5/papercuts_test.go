package html5_test

import (
	"strings"
	"testing"

	"github.com/jpl-au/fluent/html5/a"
	"github.com/jpl-au/fluent/html5/area"
	"github.com/jpl-au/fluent/html5/canvas"
	"github.com/jpl-au/fluent/html5/col"
	"github.com/jpl-au/fluent/html5/colgroup"
	"github.com/jpl-au/fluent/html5/div"
	"github.com/jpl-au/fluent/html5/dropdown"
	"github.com/jpl-au/fluent/html5/embed"
	"github.com/jpl-au/fluent/html5/img"
	"github.com/jpl-au/fluent/html5/input"
	"github.com/jpl-au/fluent/html5/meter"
	"github.com/jpl-au/fluent/html5/object"
	"github.com/jpl-au/fluent/html5/ol"
	"github.com/jpl-au/fluent/html5/option"
	"github.com/jpl-au/fluent/html5/progress"
	"github.com/jpl-au/fluent/html5/source"
	"github.com/jpl-au/fluent/html5/td"
	"github.com/jpl-au/fluent/html5/textarea"
	"github.com/jpl-au/fluent/html5/th"
	"github.com/jpl-au/fluent/html5/video"
)

// TestZeroMeaningfulAttributesRender covers attributes where zero and absent
// mean different things: tabindex="0" is focusable, rowspan="0" spans to the
// end of the row group, an absent progress value is indeterminate rather than
// empty, and an absent width falls back to the intrinsic or default size.
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
		{"minlength", string(input.New().MinLength(0).RenderBytes()), ` minlength="0"`},
		{"start", string(ol.New().Start(0).RenderBytes()), ` start="0"`},
		{"canvas width", string(canvas.New().Width(0).RenderBytes()), ` width="0"`},
		{"canvas height", string(canvas.New().Height(0).RenderBytes()), ` height="0"`},
		{"img width", string(img.New().Width(0).RenderBytes()), ` width="0"`},
		{"embed height", string(embed.New().Height(0).RenderBytes()), ` height="0"`},
		{"object width", string(object.New().Width(0).RenderBytes()), ` width="0"`},
		{"source width", string(source.New().Width(0).RenderBytes()), ` width="0"`},
		{"video height", string(video.New().Height(0).RenderBytes()), ` height="0"`},
		{"progress value", string(progress.New().Value(0).RenderBytes()), ` value="0"`},
		{"meter value", string(meter.New().Value(0).RenderBytes()), ` value="0"`},
		{"meter min", string(meter.New().Min(0).RenderBytes()), ` min="0"`},
		{"meter max", string(meter.New().Max(0).RenderBytes()), ` max="0"`},
		{"meter low", string(meter.New().Low(0).RenderBytes()), ` low="0"`},
		{"meter high", string(meter.New().High(0).RenderBytes()), ` high="0"`},
		{"meter optimum", string(meter.New().Optimum(0).RenderBytes()), ` optimum="0"`},
		{"meter ctor", string(meter.ValueMax(0, 10).RenderBytes()), ` value="0"`},
	}
	for _, tc := range cases {
		if !strings.Contains(tc.html, tc.want) {
			t.Errorf("%s: zero value should render %q, got %q", tc.name, tc.want, tc.html)
		}
	}
}

// TestEmptyMeaningfulAttributesRender covers attributes where empty and absent
// mean different things: alt="" marks a decorative image, and an option with no
// value submits its own text content.
func TestEmptyMeaningfulAttributesRender(t *testing.T) {
	cases := []struct {
		name string
		html string
		want string
	}{
		{"img alt via constructor", string(img.Image("x.png", "").RenderBytes()), ` alt=""`},
		{"img alt via setter", string(img.New().Alt("").RenderBytes()), ` alt=""`},
		{"img alt via Lazy", string(img.Lazy("x.png", "").RenderBytes()), ` alt=""`},
		{"area alt", string(area.New().Href("/x").Alt("").RenderBytes()), ` alt=""`},
		{"option value", string(option.New().Value("").RenderBytes()), ` value=""`},
	}
	for _, tc := range cases {
		if !strings.Contains(tc.html, tc.want) {
			t.Errorf("%s: empty value should render %q, got %q", tc.name, tc.want, tc.html)
		}
	}
}

// TestPointerFieldsStillEscape checks a pointer field escapes its value on the
// setter and constructor paths alike. The constructor escapes into a local
// before taking its address, so the two paths reach the escaping separately.
func TestPointerFieldsStillEscape(t *testing.T) {
	const attack = `" onerror="alert(1)`

	cases := []struct {
		name string
		html string
	}{
		{"img alt via constructor", string(img.Image("x.png", attack).RenderBytes())},
		{"img alt via setter", string(img.New().Alt(attack).RenderBytes())},
		{"area alt via setter", string(area.New().Alt(attack).RenderBytes())},
		{"option value", string(option.New().Value(attack).RenderBytes())},
	}
	for _, tc := range cases {
		if strings.Contains(tc.html, `onerror="alert(1)`) {
			t.Errorf("%s: attribute value escaped out of its quotes: %q", tc.name, tc.html)
		}
		if !strings.Contains(tc.html, "&#34;") {
			t.Errorf("%s: expected the quote to be escaped, got %q", tc.name, tc.html)
		}
	}
}

// TestInvalidZeroStillRenders covers values HTML requires to be above zero.
// They render as set: an HTML validator reports the value as invalid, and a
// dropped attribute would give the caller no output and no error.
func TestInvalidZeroStillRenders(t *testing.T) {
	cases := []struct {
		name string
		html string
		want string
	}{
		{"td colspan", string(td.New().ColSpan(0).RenderBytes()), ` colspan="0"`},
		{"th colspan", string(th.New().ColSpan(0).RenderBytes()), ` colspan="0"`},
		{"col span", string(col.New().Span(0).RenderBytes()), ` span="0"`},
		{"colgroup span", string(colgroup.New().Span(0).RenderBytes()), ` span="0"`},
		{"textarea rows", string(textarea.New().Rows(0).RenderBytes()), ` rows="0"`},
		{"textarea cols", string(textarea.New().Cols(0).RenderBytes()), ` cols="0"`},
		{"input size", string(input.New().Size(0).RenderBytes()), ` size="0"`},
		{"select size", string(dropdown.New().Size(0).RenderBytes()), ` size="0"`},
		{"progress max", string(progress.New().Value(1).Max(0).RenderBytes()), ` max="0"`},
	}
	for _, tc := range cases {
		if !strings.Contains(tc.html, tc.want) {
			t.Errorf("%s: a value the caller set should render, wanted %q, got %q", tc.name, tc.want, tc.html)
		}
	}
}

// TestUnsetAttributesStayAbsent checks the other side of the pointer fields.
// Absence is the only thing that suppresses an attribute, so an element nobody
// configured renders bare.
func TestUnsetAttributesStayAbsent(t *testing.T) {
	cases := []struct {
		name string
		html string
	}{
		{"td", string(td.New().RenderBytes())},
		{"img", string(img.New().Src("a.png").RenderBytes())},
		{"canvas", string(canvas.New().RenderBytes())},
		{"meter", string(meter.New().RenderBytes())},
		{"progress", string(progress.New().RenderBytes())},
		{"textarea", string(textarea.New().RenderBytes())},
		{"option", string(option.New().RenderBytes())},
		{"area", string(area.New().RenderBytes())},
	}
	for _, tc := range cases {
		for _, attr := range []string{"colspan=", "width=", "height=", "value=", "alt=", "rows=", "cols=", "min=", "max="} {
			if strings.Contains(tc.html, attr) {
				t.Errorf("%s: unset %s should not render, got %q", tc.name, attr, tc.html)
			}
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
