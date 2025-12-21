package text

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestTextNode(t *testing.T) {
	tests := []struct {
		name     string
		node     *TextNode
		expected string
	}{
		{
			name:     "Static",
			node:     Static("Hello"),
			expected: "Hello",
		},
		{
			name:     "Text escaped",
			node:     Text("<script>"),
			expected: "&lt;script&gt;",
		},
		{
			name:     "RawText unescaped",
			node:     RawText("<b>Bold</b>"),
			expected: "<b>Bold</b>",
		},
		{
			name:     "Textf escaped",
			node:     Textf("Hello %s", "<World>"),
			expected: "Hello &lt;World&gt;",
		},
		{
			name:     "RawTextf unescaped",
			node:     RawTextf("Hello %s", "<b>World</b>"),
			expected: "Hello <b>World</b>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test Render()
			if string(tt.node.Render()) != tt.expected {
				t.Errorf("Render() = %q, want %q", string(tt.node.Render()), tt.expected)
			}

			// Test RenderBuilder()
			var buf bytes.Buffer
			tt.node.RenderBuilder(&buf)
			if buf.String() != tt.expected {
				t.Errorf("RenderBuilder() = %q, want %q", buf.String(), tt.expected)
			}
		})
	}
}

func TestTextNode_Render_Writer(t *testing.T) {
	node := Text("Hello")
	var buf bytes.Buffer
	node.Render(&buf)
	if buf.String() != "Hello" {
		t.Errorf("Render(writer) = %q, want %q", buf.String(), "Hello")
	}
}

func BenchmarkRenderWriter(b *testing.B) {
	node := Text("Hello World")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		node.Render(io.Discard)
	}
}

func BenchmarkRenderWriterLarge(b *testing.B) {
	node := Text(strings.Repeat("Hello World ", 100))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		node.Render(io.Discard)
	}
}
