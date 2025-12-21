package security

import (
	"strings"
	"testing"

	"github.com/jpl-au/fluent/text"
)

func TestSanitise(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "Valid content",
			input:   "<div>Hello World</div>",
			wantErr: false,
		},
		{
			name:    "Script tag injection",
			input:   "<div><script>alert(1)</script></div>",
			wantErr: true,
		},
		{
			name:    "Style tag injection",
			input:   "<div><style>body { color: red; }</style></div>",
			wantErr: true,
		},
		{
			name:    "Event handler",
			input:   "<div onclick='alert(1)'>Click me</div>",
			wantErr: true,
		},
		{
			name:    "Javascript protocol",
			input:   "<a href='javascript:alert(1)'>Link</a>",
			wantErr: true,
		},
		{
			name:    "Encoded script tag",
			input:   "&lt;script&gt;alert(1)&lt;/script&gt;",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			comp := text.RawText(tt.input)
			sb := Sanitise(comp)

			// Check Error()
			errComp := sb.Error()
			rendered := string(errComp.Render())
			if tt.wantErr {
				if !strings.Contains(rendered, "Validation Error") {
					t.Errorf("Sanitise().Error() rendered %q, want error message", rendered)
				}
			} else {
				if strings.Contains(rendered, "Validation Error") {
					t.Errorf("Sanitise().Error() rendered error %q, want original content", rendered)
				}
				if rendered != tt.input {
					t.Errorf("Sanitise().Error() rendered %q, want %q", rendered, tt.input)
				}
			}

			// Check Render()
			renderedBytes := sb.Render()
			if tt.wantErr {
				if len(renderedBytes) != 0 {
					t.Errorf("Sanitise().Render() returned bytes for invalid content")
				}
			} else {
				if string(renderedBytes) != tt.input {
					t.Errorf("Sanitise().Render() = %q, want %q", string(renderedBytes), tt.input)
				}
			}
		})
	}
}

func TestSafeHelpers(t *testing.T) {
	t.Run("Safe", func(t *testing.T) {
		node := Safe("Hello World")
		if string(node.Render()) != "Hello World" {
			t.Error("Safe() failed for valid content")
		}

		node = Safe("<script>")
		if !strings.Contains(string(node.Render()), "Validation Error") {
			t.Error("Safe() failed to catch invalid content")
		}
	})

	t.Run("SafeScript", func(t *testing.T) {
		node := SafeScript("console.log('hello')")
		expected := "<script>console.log('hello')</script>"
		if string(node.Render()) != expected {
			t.Errorf("SafeScript() = %q, want %q", string(node.Render()), expected)
		}

		node = SafeScript("</script>")
		if !strings.Contains(string(node.Render()), "Validation Error") {
			t.Error("SafeScript() failed to catch invalid content")
		}
	})
}

func BenchmarkSanitise(b *testing.B) {
	comp := text.RawText("<div>Safe Content</div>")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Sanitise(comp).Render()
	}
}

func BenchmarkSanitiseComplex(b *testing.B) {
	comp := text.RawText(`
		<div class="container">
			<h1>Title</h1>
			<p>Some paragraph text with <b>bold</b> and <i>italic</i>.</p>
			<ul>
				<li>Item 1</li>
				<li>Item 2</li>
			</ul>
		</div>
	`)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Sanitise(comp).Render()
	}
}
