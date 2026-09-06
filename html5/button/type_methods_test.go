package button_test

import (
	"bytes"
	"testing"

	"github.com/jpl-au/fluent/html5/button"
	"github.com/jpl-au/fluent/html5/span"
	"github.com/jpl-au/fluent/html5/svg"
)

func TestButtonTypeMethods(t *testing.T) {
	setters := []struct {
		name, value string
		set         func(*button.Element) *button.Element
	}{
		{"Submit", "submit", (*button.Element).Submit},
		{"Reset", "reset", (*button.Element).Reset},
		{"Button", "button", (*button.Element).Button},
	}
	for _, setter := range setters {
		t.Run(setter.name, func(t *testing.T) {
			icon := svg.New(svg.Circle().R("5"))
			label := span.Text("Action")
			e := button.New(icon, label)
			if got := setter.set(e); got != e {
				t.Fatal("setter returned a different element")
			}
			if children := e.Nodes(); len(children) != 2 || children[0] != icon || children[1] != label {
				t.Fatal("setter changed the children")
			}
			want := `<button type="` + setter.value + `">` + string(icon.RenderBytes()) + `<span>Action</span></button>`
			if got := string(e.RenderBytes()); got != want {
				t.Fatalf("got %q, want %q", got, want)
			}
			var split bytes.Buffer
			e.RenderOpen(&split)
			for _, child := range e.Nodes() {
				child.RenderBuilder(&split)
			}
			e.RenderClose(&split)
			if split.String() != want {
				t.Fatalf("split render got %q, want %q", split.String(), want)
			}
		})
	}
}

func TestButtonTypeLastSetterWins(t *testing.T) {
	tests := []struct {
		name, want string
		build      func() *button.Element
	}{
		{"constructor override", "submit", func() *button.Element { return button.Reset("Action").Submit() }},
		{"submit to reset", "reset", func() *button.Element { return button.Text("Action").Submit().Reset() }},
		{"reset to button", "button", func() *button.Element { return button.Text("Action").Reset().Button() }},
		{"button to submit", "submit", func() *button.Element { return button.Text("Action").Button().Submit() }},
		{"type to submit", "submit", func() *button.Element { return button.Text("Action").Type("reset").Submit() }},
		{"submit to type", "button", func() *button.Element { return button.Text("Action").Submit().Type("button") }},
		{"repeated setter", "reset", func() *button.Element { return button.Text("Action").Reset().Reset() }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want := `<button type="` + tt.want + `">Action</button>`
			if got := string(tt.build().RenderBytes()); got != want {
				t.Fatalf("got %q, want %q", got, want)
			}
		})
	}
}
