package meta_test

import (
	"bytes"
	"testing"

	"github.com/jpl-au/fluent/html5/meta"
)

func TestNameConstructorAttributes(t *testing.T) {
	tests := []struct {
		name, content, want string
	}{
		{"application-name", "My App", `<meta name="application-name" content="My App">`},
		{`custom"<&`, `A "quote" & <tag>`, `<meta name="custom&#34;&lt;&amp;" content="A &#34;quote&#34; &amp; &lt;tag&gt;">`},
		{"custom:name", "value", `<meta name="custom:name" content="value">`},
		{"", "", `<meta>`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := meta.Name(tt.name, tt.content)
			if got := string(e.RenderBytes()); got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
			//flint:allow constructors comparing the constructor with its equivalent setters
			chained := meta.New().Name(tt.name).Content(tt.content)
			if !bytes.Equal(e.RenderBytes(), chained.RenderBytes()) {
				t.Fatal("constructor differs from the equivalent setter chain")
			}
			var split bytes.Buffer
			e.RenderOpen(&split)
			e.RenderClose(&split)
			if split.String() != tt.want {
				t.Fatalf("split render got %q, want %q", split.String(), tt.want)
			}
		})
	}
}

func TestNameConstructorCanBeOverridden(t *testing.T) {
	e := meta.Name("application-name", "My App")
	e.Name("description").Content("A description")
	want := `<meta name="description" content="A description">`
	if got := string(e.RenderBytes()); got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}
