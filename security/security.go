package security

import (
	"bytes"
	"fmt"
	"io"
	"regexp"

	"github.com/jpl-au/fluent/node"
	"github.com/jpl-au/fluent/text"
)

// dangerousPatterns contains a single compiled regular expression for known harmful or suspicious constructs
// that must be blocked from within <script> or <style> blocks.
var pattern = regexp.MustCompile(`(?i)(</\s*script\s*>)|(<\s*script)|(</\s*style\s*>)|(on\w+\s*=)|(javascript\s*:)|(eval\s*\()|(document\s*\.)|(window\s*\.)|(expression\s*\()|(url\s*\(\s*['"]?\s*javascript:)|(&lt;\s*/\s*script\s*&gt;)|(&lt;\s*/\s*style\s*&gt;)`)

// SanitiseBuilder allows fluent validation with error handling
type SanitiseBuilder struct {
	component node.Node
	err       error
	content   []byte
}

// Sanitise creates a new validation builder for the given component.
// It performs security validation checks against potentially unsafe content within the rendered output.
// This is useful when rendering raw content into <script> or <style> blocks where injection is possible.
//
// Basic usage (renders component if valid, nothing if invalid):
//
//	security.Sanitise(scriptComponent)
//
// With error fallback (renders an error message if invalid):
//
//	security.Sanitise(scriptComponent).Error()
func Sanitise(comp node.Node) *SanitiseBuilder {
	sb := &SanitiseBuilder{
		component: comp,
	}

	if comp == nil {
		sb.err = fmt.Errorf("component is nil")
		return sb
	}

	// Get the rendered content for inspection
	// We cache this content to avoid re-rendering later
	sb.content = comp.Render()

	// Apply sanitisation rule
	if pattern.Match(sb.content) {
		sb.err = fmt.Errorf("content contains disallowed pattern")
		return sb
	}

	// If no patterns matched, consider it valid
	return sb
}

// Sanitize is an alias for Sanitise
func Sanitize(comp node.Node) *SanitiseBuilder {
	return Sanitise(comp)
}

// Error returns a fallback component that renders the original component if valid,
// or an error message component if validation failed
func (sb *SanitiseBuilder) Error() node.Node {
	if sb.err != nil {
		return text.Text("Validation Error: " + sb.err.Error())
	}
	return sb.component
}

// Render renders the original component if valid, or an empty byte slice if invalid
func (sb *SanitiseBuilder) Render(w ...io.Writer) []byte {
	if sb.err != nil {
		if len(w) > 0 {
			return nil
		}
		return []byte{}
	}

	if len(w) > 0 {
		w[0].Write(sb.content)
		return nil
	}
	return sb.content
}

// RenderBuilder renders the original component into the given buffer if valid,
// or writes nothing if validation failed
func (sb *SanitiseBuilder) RenderBuilder(buf *bytes.Buffer) {
	if sb.err != nil {
		return
	}
	buf.Write(sb.content)
}

// Nodes returns the children of the wrapped component if valid,
// or an empty slice if validation failed
func (sb *SanitiseBuilder) Nodes() []node.Node {
	if sb.err != nil {
		return []node.Node{}
	}
	return sb.component.Nodes()
}

// Validate performs sanitisation check on the given content string directly
// without wrapping it in a component. Returns an error if validation fails.
func Validate(content string) error {
	if pattern.MatchString(content) {
		return fmt.Errorf("content contains disallowed pattern")
	}
	return nil
}

// Safe creates a sanitised text node from the given content.
// If the content fails validation, it returns an error text node instead.
func Safe(content string) node.Node {
	if err := Validate(content); err != nil {
		return text.Text("Validation Error: " + err.Error())
	}
	return text.RawText(content)
}

// SafeScript creates a sanitised script element with inline JavaScript.
// If the JavaScript content fails validation, it returns an error comment instead.
func SafeScript(js string) node.Node {
	if err := Validate(js); err != nil {
		return text.Text("<!-- Validation Error: " + err.Error() + " -->")
	}
	return text.RawText("<script>" + js + "</script>")
}

// SafeStyle creates a sanitised style element with inline CSS.
// If the CSS content fails validation, it returns an error comment instead.
func SafeStyle(css string) node.Node {
	if err := Validate(css); err != nil {
		return text.Text("<!-- Validation Error: " + err.Error() + " -->")
	}
	return text.RawText("<style>" + css + "</style>")
}

// Dynamic returns true as SanitiseBuilder performs runtime validation.
func (sb *SanitiseBuilder) Dynamic() bool {
	return true
}
