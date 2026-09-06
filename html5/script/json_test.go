package script_test

import (
	"bytes"
	"testing"

	"github.com/jpl-au/fluent/html5/script"
)

func TestJSONPreservesBytes(t *testing.T) {
	tests := []struct {
		name, data string
	}{
		{"object", `{"key": "value"}`},
		{"closing tag", `{"value":"</script><script>alert(1)</script>"}`},
		{"mixed case closing tag", `{"value":"</ScRiPt >"}`},
		{"HTML comment and script", `{"value":"<!--<script>"}`},
		{"HTML characters", `{"<key>":"<&> and &amp;"}`},
		{"already escaped", `{"value":"\u003c/script\u003e"}`},
		{"Unicode separators", "{\"value\":\"café\u2028\u2029\"}"},
		{"array", `["<",true,null,42]`},
		{"string", `"</script>"`},
		{"number", `42`},
		{"boolean", `false`},
		{"null", `null`},
		{"whitespace", "\n{\"value\": \"<\"}\n"},
		{"empty", ""},
		{"invalid JSON", "invalid JSON"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := script.JSON(tt.data)
			want := `<script type="application/json">` + tt.data + `</script>`
			if got := string(e.RenderBytes()); got != want {
				t.Fatalf("got %q, want %q", got, want)
			}

			var children bytes.Buffer
			for _, child := range e.Nodes() {
				child.RenderBuilder(&children)
			}
			if children.String() != tt.data {
				t.Fatalf("child content changed: got %q, want %q", children.String(), tt.data)
			}

			var split bytes.Buffer
			e.RenderOpen(&split)
			split.Write(children.Bytes())
			e.RenderClose(&split)
			if split.String() != want {
				t.Errorf("split render got %q, want %q", split.String(), want)
			}
			var rendered, built, written bytes.Buffer
			e.Render(&rendered)
			e.RenderBuilder(&built)
			n, err := e.WriteTo(&written)
			if err != nil || n != int64(len(want)) || written.String() != want || rendered.String() != want || built.String() != want {
				t.Errorf("render methods disagree: Render=%q, RenderBuilder=%q, WriteTo=%q (%d, %v)", rendered.String(), built.String(), written.String(), n, err)
			}
		})
	}
}

func TestJSONWithAttributes(t *testing.T) {
	data := `{"key":"value"}`
	got := string(script.JSON(data).ID("page-data").RenderBytes())
	want := `<script type="application/json" id="page-data">` + data + `</script>`
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func FuzzJSONPreservesBytes(f *testing.F) {
	for _, value := range []string{"plain", "</script>", "<!--<script>", "<&>\u2028\u2029"} {
		f.Add(value)
	}
	f.Fuzz(func(t *testing.T, data string) {
		got := string(script.JSON(data).RenderBytes())
		want := `<script type="application/json">` + data + `</script>`
		if got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
	})
}
