// Package text provides the leaf text nodes of the Fluent render tree.
// It backs the Static/Text/Textf/RawText/RawTextf constructors on every
// HTML element and the like-named methods used to add content after
// construction.
//
// Three escaping models cover trusted, untrusted, and pre-escaped
// content:
//
//   - [Static] never escapes and is marked non-dynamic so the JIT can
//     pre-render it. Use only with string literals you control.
//   - [Text] and [Textf] escape via html.EscapeString and are marked
//     dynamic. Use these for variables and user input.
//   - [RawText] and [RawTextf] never escape and are marked dynamic.
//     Use only with HTML you have already sanitised - pair with
//     fluent-security for untrusted markup.
//
// All non-[Static] nodes report IsDynamic() == true via the
// node.Dynamic interface so the diff engine can track them across
// renders.
//
// This package covers text *content*. The other injection surface,
// attribute *values*, is escaped separately and automatically by the
// generated setters via [node.EscapeAttribute] (with URL sinks
// scheme-filtered via [node.FilterURL]); [node.SetAttributeRaw] is the
// attribute mirror of [RawText]. Together they are one escaping system:
// content and attributes are both safe by default, with a raw hatch on
// each side for values you have already sanitised.
package text

import (
	"bytes"
	"fmt"
	"html"
	"io"

	"github.com/jpl-au/fluent"
	"github.com/jpl-au/fluent/node"
)

// Node represents text content that can be either HTML-escaped (safe) or raw (unescaped).
// It implements the node.Node interface and is used internally by Text() and RawText()
// constructor functions to handle different security models.
type Node struct {
	content string // The text content, possibly HTML-escaped
	dynamic bool   // Whether the content is dynamically generated
}

// Static creates a text node for compile-time constant strings. The content
// is NOT HTML-escaped and is marked as non-dynamic, allowing the JIT to
// pre-render it. Only use with string literals you control - never with
// user input or dynamic values, as this would create an XSS vulnerability.
//
// For dynamic or user-provided content, use [Text] or [Textf] instead.
//
// Example:
//
//	text.Static("Copyright 2024") // Renders as: Copyright 2024
func Static(str string) *Node {
	return &Node{
		content: str,
		dynamic: false,
	}
}

// Text creates a dynamic text node with automatic HTML escaping. The
// characters <, >, &, and quotes are escaped via html.EscapeString so
// user-supplied content cannot inject markup.
//
// Example:
//
//	text.Text(userName)
func Text(str string) *Node {
	return &Node{
		content: html.EscapeString(str),
		dynamic: true,
	}
}

// RawText creates a dynamic, unescaped text node for trusted HTML
// content. Use only with content you control, such as pre-built HTML
// strings; pair with fluent-security to sanitise untrusted markup
// before passing it here.
//
// Example:
//
//	text.RawText(htmlContent)
func RawText(str string) *Node {
	return &Node{
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
func Textf(format string, a ...any) *Node {
	return &Node{
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
func RawTextf(format string, a ...any) *Node {
	return &Node{
		content: fmt.Sprintf(format, a...),
		dynamic: true,
	}
}

// RenderBuilder writes the text content directly to the provided buffer.
// This method provides efficient rendering for large node trees.
func (tn *Node) RenderBuilder(buf *bytes.Buffer) {
	buf.WriteString(tn.content)
}

// Render writes the text content to w.
// Write errors are intentionally discarded; see [node.Node] for rationale.
func (tn *Node) Render(w io.Writer) {
	_, _ = tn.WriteTo(w)
}

// WriteTo writes the text content to w, returning the byte count and
// any write error. Satisfies [io.WriterTo].
func (tn *Node) WriteTo(w io.Writer) (int64, error) {
	buf := fluent.NewBuffer()
	tn.RenderBuilder(buf)
	n, err := buf.WriteTo(w)
	fluent.PutBuffer(buf)
	return n, err
}

// RenderBytes returns the text content as a byte slice.
func (tn *Node) RenderBytes() []byte {
	var buf bytes.Buffer
	tn.RenderBuilder(&buf)
	return buf.Bytes()
}

// Nodes returns an empty slice as text nodes do not have children.
func (tn *Node) Nodes() []node.Node {
	return []node.Node{}
}

// IsDynamic returns true if this text content is dynamically generated (created with Text, Textf, RawText, or RawTextf).
// Static content (created with Static) returns false, allowing JIT to pre-render it.
func (tn *Node) IsDynamic() bool {
	return tn.dynamic
}

// DynamicKey returns an empty string - text nodes do not carry tracking keys.
func (tn *Node) DynamicKey() string {
	return ""
}

// String returns the render-ready content: exactly the bytes this node
// will emit when rendered. That means the escaped form for [Text] and
// [Textf] nodes - escaping happens at construction, not at render - and
// the verbatim input for [Static], [RawText] and [RawTextf] nodes.
// Printing a Text node with %s therefore shows &lt; rather than the
// original <. Intended for tests and debugging; use the Render methods
// for output.
func (tn *Node) String() string {
	return tn.content
}
