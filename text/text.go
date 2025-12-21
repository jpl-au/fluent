package text

import (
	"bytes"
	"fmt"
	"html"
	"io"

	"github.com/jpl-au/fluent"
	"github.com/jpl-au/fluent/node"
)

// TextNode represents text content that can be either HTML-escaped (safe) or raw (unescaped).
// It implements the node.Node interface and is used internally by Text() and RawText()
// constructor functions to handle different security models.
type TextNode struct {
	content string // The text content, possibly HTML-escaped
	dynamic bool   // Whether the content is dynamically generated
}

// Static creates a text component that is explicitly marked as static content.
// This should be used for content that never changes, allowing JIT optimisation.
// Warning: Static content is not html encoded/escaped.
//
// Example:
//
//	text.Static("Copyright 2024") // Renders as: Copyright 2024
func Static(str string) *TextNode {
	return &TextNode{
		content: str,
		dynamic: false,
	}
}

// Text creates a safe text component with automatic HTML escaping for security.
// This is now marked as dynamic by default to handle variables and user content.
//
// Special HTML characters like <, >, &, and quotes are automatically escaped.
//
// Example:
//
//	text.Text(userName) // Renders with HTML escaping, marked as dynamic
func Text(str string) *TextNode {
	return &TextNode{
		content: html.EscapeString(str),
		dynamic: true,
	}
}

// RawText creates an unescaped text component for trusted HTML content.
// This is now marked as dynamic by default to handle variables and dynamic content.
// Use this ONLY for content you control, such as pre-built HTML strings.
//
// Example:
//
//	text.RawText(htmlContent) // Renders unescaped, marked as dynamic
func RawText(str string) *TextNode {
	return &TextNode{
		content: str,
		dynamic: true,
	}
}

// Textf creates a safe, formatted text component with automatic HTML escaping.
// It works like fmt.Sprintf but ensures the final string is properly escaped
// to prevent XSS attacks.
//
// Example:
//
//	text.Textf("Hello, %s!", "<world>") // Renders as: Hello, &lt;world&gt;!
func Textf(format string, a ...any) *TextNode {
	return &TextNode{
		content: html.EscapeString(fmt.Sprintf(format, a...)),
		dynamic: true,
	}
}

// RawTextf creates a formatted text component without HTML escaping.
// It should only be used with trusted format strings and arguments.
//
// Example:
//
//	text.RawTextf("<a href='%s'>%s</a>", "/home", "Home") // Renders as: <a href='/home'>Home</a>
func RawTextf(format string, a ...any) *TextNode {
	return &TextNode{
		content: fmt.Sprintf(format, a...),
		dynamic: true,
	}
}

// RenderBuilder writes the text content directly to the provided buffer.
// This method provides efficient rendering for large node trees.
func (tn *TextNode) RenderBuilder(buf *bytes.Buffer) {
	buf.WriteString(tn.content)
}

// Render returns the text content as a byte slice or writes to the provided writer.
func (tn *TextNode) Render(w ...io.Writer) []byte {
	buf := fluent.NewBuffer()
	tn.RenderBuilder(buf)

	if len(w) > 0 && w[0] != nil {
		buf.WriteTo(w[0])
		fluent.PutBuffer(buf)
		return nil
	}
	return buf.Bytes()
}

// Nodes returns an empty slice as text nodes do not have children.
func (tn *TextNode) Nodes() []node.Node {
	return []node.Node{}
}

// Dynamic returns true if this text content is dynamically generated (created with Textf or RawTextf)
func (tn *TextNode) Dynamic() bool {
	return tn.dynamic
}

// Base returns nil as RawText nodes do not have attributes.
// SetAttribute is a no-op for TextNode as it does not have attributes.
func (tn *TextNode) SetAttribute(key string, value string) {
	// TextNode does not support attributes
}

// String returns the text content as a string.
// This allows RawText to be used in contexts that require string values.
func (tn *TextNode) String() string {
	return tn.content
}
