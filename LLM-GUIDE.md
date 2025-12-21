# Fluent LLM Guide

Guide to help LLMs work with the Fluent HTML generation framework.

## Core Concepts

### Everything is a Node

The `node.Node` interface is Fluent's foundation. Every renderable piece implements it: HTML elements, text nodes, conditionals (`node.Condition`), and function wrappers (`node.Func`, `node.FuncNodes`). This enables arbitrary composition - any `node.Node` can be a child of any element.

**When in doubt, return `node.Node`** - it's always safe and provides maximum flexibility:

```go
func MyComponent(showHeader bool) node.Node {
    if showHeader {
        return header.New(h1.Text("Welcome"))
    }
    return nil  // nil nodes are safely skipped during rendering
}
```

Returning concrete types (like `*div.Element`) allows method chaining after the call. See [Component Pattern](#component-pattern) for detailed guidance on when to use each approach.

### Static vs Text Rendering

**Static()** - Immutable content that is known at template definition time. It is useful to use Static() when working with the JIT optimisations
```go
div.Static("Copyright 2024")  // Content never changes
```

**Text()** - HTML-escaped dynamic content
```go
div.Text(user.Name)  // Escaped at runtime, value can change per render
```

**Textf()** - HTML-escaped dynamic content with formatting
```go
div.Textf("Hello %s, you have %d messages", user.Name, count)  // Escaped, formatted
```

**RawText()** - Unescaped HTML content
```go
div.RawText("<em>Bold</em>")  // Not escaped, use carefully
```

**RawTextf()** - Unescaped HTML content with formatting
```go
div.RawTextf("<span class=\"%s\">%s</span>", className, content)  // Not escaped, formatted
```

**Rule:** Use `Static()` for unchanging content (labels, headings, boilerplate). Use `Text()` or `Textf()` for user input or values that change between renders. Use `RawText()` or `RawTextf()` only when you need to inject HTML and trust the source.

### Security Package

`Text()` and `Textf()` use Go's `html.EscapeString()` for basic HTML escaping (prevents `<`, `>`, `&`, quotes from being interpreted as HTML). For content injected into `<script>` or `<style>` blocks, use the `security` package which detects dangerous patterns:

```go
import "github.com/jpl-au/fluent/security"

// Sanitise rendered output - returns empty if dangerous patterns found
security.Sanitise(scriptComponent).Render()

// Sanitise with error fallback - renders error message if invalid
security.Sanitise(scriptComponent).Error()

// Validate a string directly
if err := security.Validate(userInput); err != nil {
    // Contains dangerous pattern
}

// Safe wrappers for script/style content
security.SafeScript(jsCode)  // Returns sanitised <script> or error comment
security.SafeStyle(cssCode)  // Returns sanitised <style> or error comment
```

**Detected patterns:** `</script>`, `</style>`, `<script`, `onclick=` (and other event handlers), `javascript:`, `eval(`, `document.`, `window.`, `expression(`, and their HTML-encoded equivalents.

**Rule:** Use `Text()`/`Textf()` for general content. Use the `security` package when injecting content into `<script>` or `<style>` blocks.

### Type Safety

Fluent uses typed constants for attributes with enumerated values. Methods like `InputType()` accept a typed constant (e.g., `inputtype.Email`), not a string - so `input.New().InputType("emial")` won't compile.

Each attribute package provides a `Custom()` function as an escape hatch for edge cases or future HTML specifications not yet covered by predefined constants:

```go
// Standard usage - typed constants
input.New().InputType(inputtype.Email)

// Escape hatch for edge cases
input.New().InputType(inputtype.Custom("future-type"))
```

This approach guides developers towards valid values while allowing flexibility when needed.

### Reserved Keyword Alternatives

**CRITICAL:** Some HTML elements use Go reserved keywords. Fluent provides alternative package names:

| HTML Element | Reserved Keyword | Fluent Package | Import Path |
|--------------|------------------|----------------|-------------|
| `<select>`   | `select`         | `dropdown`     | `github.com/jpl-au/fluent/html5/dropdown` |
| `<main>`     | `main`           | `primary`      | `github.com/jpl-au/fluent/html5/primary` |
| `<var>`      | `var`            | `variable`     | `github.com/jpl-au/fluent/html5/variable` |

**Usage:**
```go
// CORRECT - Use alternative package names
import (
    "github.com/jpl-au/fluent/html5/dropdown"  // <select>
    "github.com/jpl-au/fluent/html5/primary"   // <main>
    "github.com/jpl-au/fluent/html5/variable"  // <var>
)

dropdown.New(...)  // Renders <select>...</select>
primary.New(...)   // Renders <main>...</main>
variable.New(...)  // Renders <var>...</var>

// INCORRECT - Don't use these
import "github.com/jpl-au/fluent/html5/select"  // Does not exist
import "github.com/jpl-au/fluent/html5/main"    // Does not exist
import "github.com/jpl-au/fluent/html5/var"     // Does not exist
```

**Rule:** Always use the Fluent package name (dropdown, primary, variable), never the HTML element name when it's a Go reserved keyword.

## Element Construction

All elements follow consistent constructor patterns:

```go
// Basic element
div.New()                              // <div></div>

// With text content (escaped) - constructor style
div.Text("Hello")                      // <div>Hello</div>

// With text content (escaped) - method chain style
div.New().Text("Hello")                // <div>Hello</div>

// With static content (JIT-optimisable)
div.Static("Footer")                   // <div>Footer</div>
div.New().Static("Footer")             // <div>Footer</div>

// With raw HTML (unescaped)
div.RawText("<em>Bold</em>")           // <div><em>Bold</em></div>
div.New().RawText("<em>Bold</em>")     // <div><em>Bold</em></div>

// With formatted text
div.Textf("Hello %s", name)            // <div>Hello John</div>
div.New().Textf("Hello %s", name)      // <div>Hello John</div>

// With child nodes
div.New(
    p.Text("Paragraph"),
    span.Text("Inline"),
)

// Chained attributes and content
div.New().Class("container").ID("main").Text("Content")
```

### Content Methods

All non-self-closing elements have these chainable content methods:

- `.Text(s)` - adds escaped text content
- `.Textf(format, args...)` - adds formatted escaped text
- `.Static(s)` - adds static text (JIT-optimisable)
- `.RawText(s)` - adds unescaped HTML content
- `.RawTextf(format, args...)` - adds formatted unescaped HTML

```go
// Flexible chaining - add content at any point
div.New().Class("foo").Text("Hello").ID("bar")
p.New().Text("Line 1").Text(" Line 2")  // Multiple text nodes
style.New().RawText("body { color: red; }")
```

### Node Management Methods

- `.Add(nodes...)` - appends child nodes to the element
- `.Replace(nodes...)` - replaces all child nodes with the provided nodes

```go
// Add children after construction
container := div.New().Class("container")
container.Add(h1.Text("Title"), p.Text("Content"))

// Replace all children
container.Replace(span.Text("New content"))
```

### Attribute Methods

**CRITICAL:** Use the correct method for setting attributes:

**SetAria(key, value)** - For ARIA attributes (automatically adds "aria-" prefix)
```go
// CORRECT
button.New().SetAria("label", "Close dialog")
// Renders: <button aria-label="Close dialog"></button>

div.New().SetAria("hidden", "true")
// Renders: <div aria-hidden="true"></div>

// INCORRECT - Don't use SetAttribute for ARIA
button.New().SetAttribute("aria-label", "Close")  // Don't do this
```

**SetData(key, value)** - For data attributes (automatically adds "data-" prefix)
```go
// CORRECT
div.New().SetData("id", "123")
// Renders: <div data-id="123"></div>

button.New().SetData("action", "submit")
// Renders: <button data-action="submit"></button>

// INCORRECT - Don't use SetAttribute for data attributes
div.New().SetAttribute("data-id", "123")  // Don't do this
```

**SetAttribute(key, value)** - Only for custom/non-standard attributes
```go
// Use only when SetAria() and SetData() don't apply
div.New().SetAttribute("custom-attr", "value")
div.New().SetAttribute("x-on:click", "handler")  // Alpine.js
```

**IMPORTANT RULES:**
1. **ALL standard HTML attributes have proper typed methods** - Never use `SetAttribute()` for standard attributes like `defer`, `async`, `disabled`, `required`, `title`, `alt`, `href`, `src`, etc. These all have dedicated methods (e.g., `.Defer()`, `.Async()`, `.Disabled()`, `.Required()`, `.Title()`, `.Alt()`, `.Href()`, `.Src()`). If you need to set custom attributes, the correct method is `.SetAttribute()`, but first verify no dedicated method exists.
2. Always prefer `SetAria()` for ARIA attributes and `SetData()` for data attributes.
3. Only use `SetAttribute()` for truly custom/non-standard attributes (e.g., Alpine.js directives, custom framework attributes).
4. If you find yourself using `SetAttribute()` for a standard HTML attribute, you're doing it wrong - check the element's available methods first.

### HTML Document Construction

```go
// Complete document with DOCTYPE (default)
html.New(
    head.New(title.Text("Page")),
    body.New(div.Text("Content")),
)
// Renders: <!DOCTYPE html><html>...</html>

// Fragment without DOCTYPE (rare)
html.Fragment(...)
// Renders: <html>...</html>
```

## Dynamic Content

### Conditional Rendering

`node.Condition()` provides inline conditional rendering with `True()` and `False()` branches:

```go
// Both branches
node.Condition(user.IsLoggedIn).
    True(div.Text("Welcome back!")).
    False(div.Text("Please log in"))
```

For single-branch conditions, `When()` and `Unless()` provide concise shorthand:

```go
// Render only when condition is true
node.When(user.IsAdmin, span.Static("Admin"))

// Render only when condition is false
node.Unless(user.IsLoggedIn, a.New().Href("/login").Text("Sign in"))
```

Conditions can be nested since they return `node.Node`:

```go
node.Condition(user.IsLoggedIn).
    True(
        node.Condition(user.IsAdmin).
            True(span.Static("Admin Dashboard")).
            False(span.Static("User Dashboard")),
    ).
    False(a.New().Href("/login").Text("Sign in"))
```

For multiple branches, `node.Func()` is cleaner than deeply nested conditions:

```go
node.Func(func() node.Node {
    if !user.IsLoggedIn {
        return a.New().Href("/login").Text("Sign in")
    }
    if user.IsAdmin {
        return span.Static("Admin Dashboard")
    }
    return span.Static("User Dashboard")
})
```

**Summary:**
- `Condition(bool).True(node).False(node)` - full conditional with both branches
- `Condition(bool).True(node)` - render only when true
- `Condition(bool).False(node)` - render only when false
- `When(bool, node)` - shorthand for `Condition(bool).True(node)`
- `Unless(bool, node)` - shorthand for `Condition(!bool).True(node)`

Nil nodes are safely ignored - if `nil` is passed to `True()` or `False()`, nothing renders for that path.

### Function Component (Single Node)

`node.Func()` executes a function during render that returns a single node:

```go
node.Func(func() node.Node {
    if user.Role == "admin" {
        return div.Text("Admin Panel")
    }
    return div.Text("User Dashboard")
})
```

### Function Component (Multiple Nodes)

`node.FuncNodes()` executes a function during render that returns multiple nodes:

```go
node.FuncNodes(func() []node.Node {
    nodes := []node.Node{}
    for _, item := range items {
        nodes = append(nodes, li.Text(item.Name))
    }
    return nodes
})
```

This is useful for generating lists without a wrapper element. Nil nodes in the returned slice are safely ignored.

## Component Pattern

Components are functions returning either `node.Node` (interface) or a concrete element type (e.g., `*div.Element`, `*span.Element`).

### Return Types: Interface vs Concrete

**`node.Node` (interface)** - Use when:
- The function may return different element types conditionally
- The component is a final building block (callers won't chain additional methods)
- You want maximum flexibility

**`*element.Element` (concrete type)** - Use when:
- Callers should be able to chain additional methods after the call
- The function always returns the same element type
- You want IDE auto-completion for the returned element's methods

Each HTML element package exports an `Element` type alias (e.g., `div.Element`, `span.Element`, `a.Element`). These are the concrete types returned by constructors like `div.New()`, `span.Text()`, etc.

### Examples

```go
// Return node.Node - flexible but no chaining after call
func Card(title, content string) node.Node {
    return div.New(
        h2.Text(title),
        p.Text(content),
    ).Class("card")
}

// Usage: Card is a complete unit
div.New(Card("Welcome", "Hello!"))
```

```go
// Return *div.Element - allows continued chaining
func Card(title, content string) *div.Element {
    return div.New(
        h2.Text(title),
        p.Text(content),
    ).Class("card")
}

// Usage: Caller can add more attributes
Card("Welcome", "Hello!").ID("welcome-card").Class("highlighted")
```

```go
// Return node.Node when conditionally returning different types
func ActionButton(user User) node.Node {
    if user.IsAdmin {
        return button.Text("Delete").Class("btn-danger")
    }
    return span.Text("No action available")
}
```

```go
// Return concrete type for consistent element with chaining
func NavLink(href, text string) *a.Element {
    return a.New().Href(href).Text(text).Class("nav-link")
}

// Caller can extend:
NavLink("/about", "About Us").Class("active").SetData("section", "about")
```

### Common Element Types

All element packages export their concrete type as `Element`:
- `*div.Element`, `*span.Element`, `*p.Element`
- `*a.Element`, `*button.Element`, `*input.Element`
- `*form.Element`, `*table.Element`, `*ul.Element`, `*li.Element`
- `*h1.Element`, `*h2.Element`, ... `*h6.Element`
- `*img.Element`, `*link.Element`, `*script.Element`
- And all other HTML5 elements

### When to Use Each

| Scenario | Return Type | Reason |
|----------|-------------|--------|
| Reusable UI component (card, modal) | `*div.Element` | Allows caller customisation |
| Navigation link builder | `*a.Element` | Caller may add classes, data attributes |
| Conditional element (may be div or span) | `node.Node` | Different types possible |
| Layout wrapper | `node.Node` | Usually complete, no chaining needed |
| Form field with label + input | `node.Node` | Composite, returns wrapper div |
| Single styled button | `*button.Element` | Allows caller to add handlers, classes |

### Component with Internal Conditionals

```go
func UserProfile(user User) node.Node {
    return div.New(
        img.New().Src(user.Avatar).Alt(user.Name),
        h3.Text(user.Name),
        node.When(user.IsVerified, span.Static("✓ Verified").Class("badge")),
    ).Class("profile")
}
```

**Rule of thumb:** If the component always returns the same element type and callers might want to customise it, return the concrete type. If it's a complete unit or may return different types, return `node.Node`.

## Common Patterns

### Layout with Dynamic Content

```go
func Layout(title string, content node.Node) node.Node {
    return html.New(
        head.New(
            title.Text(title),
            link.New().Rel(rel.Stylesheet).Href("/app.css"),
        ),
        body.New(
            header.Static("My Site"),
            primary.New(content),
            footer.Static("© 2024"),
        ),
    )
}
```

### Conditional Attributes

```go
func Button(text string, isPrimary bool) node.Node {
    btn := button.Text(text)

    if isPrimary {
        btn.Class("btn-primary")
    } else {
        btn.Class("btn-secondary")
    }

    return btn
}
```

### List Rendering

```go
func ProductList(products []Product) node.Node {
    items := make([]node.Node, len(products))
    for i, p := range products {
        items[i] = li.New(
            h3.Text(p.Name),
            p.Text(p.Description),
            span.Textf("$%.2f", p.Price),
        )
    }
    return ul.New(items...)
}
```

## Performance

- Use `Static()` for unchanging content (enables JIT optimisation if used)
- Build components for reuse
- Buffer pooling is enabled by default and handled automatically
- Avoid string concatenation in hot paths - use `RenderBuilder()`
- Reuse components vs recreating nodes

## JIT Optimisation

For high-throughput applications, [Fluent JIT](https://github.com/jpl-au/fluent-jit) provides additional optimisation strategies:

- **Compile** - Pre-render static portions, re-evaluate dynamic content via path navigation
- **Tune** - Adaptive buffer sizing that learns optimal sizes over time
- **Flatten** - Pre-render fully static content to raw bytes

The base Fluent API performs well with automatic buffer pooling. Apply JIT selectively after profiling to identify actual bottlenecks.

See the [Fluent JIT LLM Guide](https://github.com/jpl-au/fluent-jit/blob/main/LLM-GUIDE.md) for detailed API reference and usage patterns

## Buffer Management

Fluent uses buffer pooling for allocation efficiency. Pooling is enabled by default - when you call `Render(w)` with a writer, pooled buffers are used automatically.

```go
func handler(w http.ResponseWriter, r *http.Request) {
    page := html.New(
        head.New(title.Text("My Page")),
        body.New(div.Text("Hello")),
    )
    page.Render(w)  // Pooled buffer used automatically
}
```

Without a hint, renders still benefit from pooling - buffers are retrieved from the small pool, which over time will contain pre-grown buffers from previous renders.

### BufferHint (Optional)

Each element has a `BufferHint()` method to help determine which pool will be used when retrieving a bytes.Buffer and to grow the buffer to the appropriate hint. After `Render(w)`, the hint is updated to reflect the actual rendered size, which you can retrieve and reuse:

```go
// First render - set a hint if you know approximate size
page := html.New(...)
page.BufferHint(8192)  // Hint at 8KB
page.Render(w)

// Get the actual size for reuse on similar pages
actualSize := page.BufferHint()  // e.g., 6543

// Use that hint for a new page with similar content
anotherPage := html.New(...)
anotherPage.BufferHint(actualSize)
anotherPage.Render(w)
```

### Direct Buffer Access

When implementing custom `node.Node` types or for advanced use cases:

```go
buf := fluent.NewBuffer(hint)  // Get buffer from pool with optional hint
element.RenderBuilder(buf)
output := buf.Bytes()
fluent.PutBuffer(buf)              // Return buffer to pool
```

### Pool Configuration

Configure globally via the `pool` package:
```go
pool.SetEnabled(false)                    // Disable pooling entirely
pool.SetPoolThreshold(4096)               // Small vs large pool threshold
pool.SetMaxPoolSize(65536, true)          // Max size to pool, discard oversized
```

**Defaults:** Enabled: true, Threshold: 4KB, Max: 256KB, Discard oversized: true

### How the Two-Tier Pool Works

Fluent uses two separate `sync.Pool` instances: a small pool and a large pool. The threshold (default 4KB) determines which pool a buffer is routed to.

**Get behaviour:**
1. Based on the size hint, retrieves from either smallPool (hint < threshold) or largePool (hint >= threshold)
2. Calls `Reset()` on the buffer (sets length to 0, **capacity unchanged**)
3. Calls `Grow(hint)` if the hint exceeds current capacity

**Put behaviour:**
1. Routes back to pool based on **actual capacity** (not original hint)
2. A buffer pulled from smallPool that grows beyond threshold will be routed to largePool on return
3. Buffers exceeding maxPoolSize are discarded (if discardOversized is true)

**Trade-offs:**

*Memory retention:* When a buffer renders content and is returned to the pool, it retains its grown capacity. A 512-byte buffer that grows to 3KB during use stays at 3KB capacity when pooled. This means:
- Subsequent renders reuse the pre-grown buffer without reallocation
- Small renders using a pre-grown buffer have "wasted" capacity
- Buffers naturally migrate toward the size of their largest render

*Threshold tuning:* The threshold affects memory efficiency:
- **Too high (e.g., 8KB):** Small fragments may pull buffers with large unused capacity
- **Too low (e.g., 1KB):** Mid-size content constantly migrates to largePool, reducing small pool effectiveness
- **4KB default:** Balances typical fragment sizes (~500 bytes to 4KB) against occasional larger renders

*Practical impact:* For applications serving many small fragments alongside occasional full page renders, the two-tier approach keeps small and large buffers separated, preventing fragment renders from inheriting oversized buffers from page renders

## Extending Fluent

Implement `node.Node` interface for custom elements or components.

### Node Interface

```go
type Node interface {
    Render(w ...io.Writer) []byte
    RenderBuilder(*bytes.Buffer)
    Nodes() []Node
    SetAttribute(key string, value string)
}
```

### Implementation Example

When implementing custom elements, leverage the pre-defined `[]byte` constants from `github.com/jpl-au/fluent/html5` for optimal performance. Using `buf.Write()` with byte slice constants is more efficient than `buf.WriteString()` as it avoids runtime string-to-bytes conversion.

```go
import (
    "bytes"
    "io"
    "github.com/jpl-au/fluent"
    "github.com/jpl-au/fluent/html5"
    "github.com/jpl-au/fluent/node"
)

type CustomCard struct {
    title   string
    content string
    nodes   []node.Node
}

func (c *CustomCard) Render(w ...io.Writer) []byte {
    buf := fluent.NewBuffer()
    c.RenderBuilder(buf)

    if len(w) > 0 && w[0] != nil {
        buf.WriteTo(w[0])
        fluent.PutBuffer(buf)
        return nil
    }

    output := buf.Bytes()
    fluent.PutBuffer(buf)
    return output
}

func (c *CustomCard) RenderBuilder(buf *bytes.Buffer) {
    // Use constants for tags and markup - more efficient than WriteString
    buf.Write(html5.TagDiv)
    buf.Write(html5.AttrClass)
    buf.WriteString("card")
    buf.Write(html5.MarkupQuote)
    buf.Write(html5.MarkupCloseTag)

    // Opening h2 tag
    buf.Write(html5.TagH2)
    buf.Write(html5.MarkupCloseTag)
    buf.WriteString(c.title)
    buf.Write(html5.TagH2Close)

    // Opening p tag
    buf.Write(html5.TagP)
    buf.Write(html5.MarkupCloseTag)
    buf.WriteString(c.content)
    buf.Write(html5.TagPClose)

    // Render children
    for _, child := range c.nodes {
        child.RenderBuilder(buf)
    }

    buf.Write(html5.TagDivClose)
}

func (c *CustomCard) Nodes() []node.Node {
    return c.nodes
}

func (c *CustomCard) SetAttribute(key string, value string) {
    // Store attributes if needed
}
```

**Constants in `html5` package:**
- Tags: `TagDiv`, `TagH1`, `TagP`, `TagDivClose`, etc.
- Markup: `MarkupCloseTag` (`>`), `MarkupSelfCloseTag` (` />`), `MarkupQuote`, `MarkupEquals`, `MarkupSpace`
- Attributes: `AttrClass`, `AttrID`, `AttrStyle`, `AttrHref`, `AttrSrc`, `AttrHidden`, `AttrRequired`, etc.

Use `buf.Write(html5.TagDiv)` vs `buf.WriteString("<div")` - pre-allocated `[]byte` avoids string-to-byte conversion.

### Composite Components

Combine multiple HTML elements with type-safe attributes in fluent API.

**Email Field Component:**

```go
package field

import (
    "bytes"
    "io"
    "github.com/jpl-au/fluent/html5/div"
    "github.com/jpl-au/fluent/html5/label"
    "github.com/jpl-au/fluent/html5/input"
    "github.com/jpl-au/fluent/html5/attr/autocomplete"
    "github.com/jpl-au/fluent/node"
)

// EmailField renders a complete form field with label and input
type EmailField struct {
    labelText   string
    id          string
    name        string
    placeholder string
    required    bool
    class       string
    inputClass  string
    labelClass  string
    attr        *[]node.Attribute
}

// Email creates a new email field component
func Email(name string, labelText string) *EmailField {
    return &EmailField{
        labelText: labelText,
        id:        name,
        name:      name,
    }
}

// Fluent API methods
func (f *EmailField) Placeholder(text string) *EmailField {
    f.placeholder = text
    return f
}

func (f *EmailField) Required() *EmailField {
    f.required = true
    return f
}

func (f *EmailField) Class(class string) *EmailField {
    f.class = class
    return f
}

func (f *EmailField) InputClass(class string) *EmailField {
    f.inputClass = class
    return f
}

func (f *EmailField) LabelClass(class string) *EmailField {
    f.labelClass = class
    return f
}

// Implement node.Node interface
func (f *EmailField) Render(w ...io.Writer) []byte {
    var buf bytes.Buffer
    f.RenderBuilder(&buf)
    if len(w) > 0 && w[0] != nil {
        w[0].Write(buf.Bytes())
        return nil
    }
    return buf.Bytes()
}

func (f *EmailField) RenderBuilder(buf *bytes.Buffer) {
    // Build label with type-safe attributes
    labelElem := label.For(f.id, f.labelText)
    if f.labelClass != "" {
        labelElem.Class(f.labelClass)
    }

    // Build input using type-safe constants
    inputElem := input.Email(f.name).
        ID(f.id).
        AutoComplete(autocomplete.Email)

    if f.placeholder != "" {
        inputElem.Placeholder(f.placeholder)
    }

    if f.required {
        inputElem.Required()
    }

    if f.inputClass != "" {
        inputElem.Class(f.inputClass)
    }

    // Combine into container
    container := div.New(labelElem, inputElem)

    if f.class != "" {
        container.Class(f.class)
    }

    // Apply custom attributes
    if f.attr != nil {
        for _, attr := range *f.attr {
            container.SetAttribute(attr.Key, attr.Value)
        }
    }

    container.RenderBuilder(buf)
}

func (f *EmailField) Nodes() []node.Node {
    return nil
}

func (f *EmailField) SetAttribute(key string, value string) {
    if f.attr == nil {
        f.attr = &[]node.Attribute{}
    }
    for i, attr := range *f.attr {
        if attr.Key == key {
            (*f.attr)[i].Value = value
            return
        }
    }
    *f.attr = append(*f.attr, node.Attribute{Key: key, Value: value})
}
```

**Usage:**

```go
// Create complete form field with fluent API
emailField := field.Email("email", "Email Address").
    Placeholder("Enter your email address").
    Required().
    Class("form-group").
    InputClass("form-control").
    LabelClass("form-label")

// Renders:
// <div class="form-group">
//   <label for="email" class="form-label">Email Address</label>
//   <input type="email" id="email" name="email" autocomplete="email"
//          placeholder="Enter your email address" required class="form-control" />
// </div>

// Use as a node.Node
body.New(
    h1.Text("Registration"),
    emailField,
    field.Email("confirm", "Confirm Email").Required(),
)
```

### Type-Safe Composite Components

Leverage Fluent's type-safe constants when building composites:

```go
import (
    "github.com/jpl-au/fluent/html5/attr/inputtype"
    "github.com/jpl-au/fluent/html5/attr/autocomplete"
)

type TextField struct {
    label       string
    name        string
    inputType   inputtype.InputType
    autocomplete autocomplete.AutoComplete
}

func Text(name, label string) *TextField {
    return &TextField{
        label:        label,
        name:         name,
        inputType:    inputtype.Text,
        autocomplete: autocomplete.Off,
    }
}

func (t *TextField) AutoComplete(ac autocomplete.AutoComplete) *TextField {
    t.autocomplete = ac
    return t
}

func (t *TextField) RenderBuilder(buf *bytes.Buffer) {
    labelElem := label.For(t.name, t.label)
    inputElem := input.New().
        Type(t.inputType).
        Name(t.name).
        ID(t.name).
        AutoComplete(t.autocomplete)

    div.New(labelElem, inputElem).RenderBuilder(buf)
}
```

**Benefits:** Encapsulation, type safety, consistency, reusability, single source of truth, fluent API chaining.

### Best Practises

- Store state, not rendered HTML
- Build during `RenderBuilder()`, not construction
- Flexible constructors with sensible defaults
- Methods can update multiple internal elements (e.g., ID updates both label `for` and input `id`)

### Wrapper Pattern

For adding attributes to existing nodes:

```go
type AttributeWrapper struct {
    node node.Node
    attr *[]node.Attribute
}

func Wrap(n node.Node) *AttributeWrapper {
    return &AttributeWrapper{node: n}
}

func (w *AttributeWrapper) SetAttribute(key, value string) {
    if w.attr == nil {
        slice := make([]node.Attribute, 0, 1)
        w.attr = &slice
    }

    // Update existing or add new
    for i, existing := range *w.attr {
        if existing.Key == key {
            (*w.attr)[i].Value = value
            return
        }
    }

    *w.attr = append(*w.attr, node.Attribute{Key: key, Value: value})
}

func (w *AttributeWrapper) RenderBuilder(buf *bytes.Buffer) {
    // Apply wrapper attributes to node
    if w.attr != nil {
        for _, attr := range *w.attr {
            w.node.SetAttribute(attr.Key, attr.Value)
        }
    }
    w.node.RenderBuilder(buf)
}

func (w *AttributeWrapper) Render(wr ...io.Writer) []byte {
    return w.node.Render(wr...)
}

func (w *AttributeWrapper) Nodes() []node.Node {
    return w.node.Nodes()
}
```

**Use:** Custom elements, third-party libraries, framework wrappers (HTMX, Alpine.js), specialised rendering.

## Typed Attributes Reference

Elements with typed attribute constants:

### Input Elements
- `input.Accept()` - File types (accept package: ImageJPEG, ImagePNG, VideMP4, AudioMP3, Pdf, Docx, etc.)
- `input.AutoComplete()` - Autocomplete hints (autocomplete package: On, Off, Name, Email, Username, etc.)
- `input.InputType()` - Input types (inputtype package: Text, Email, Password, Number, Tel, URL, etc.)
- `input.Capture()` - Camera capture (capture package: User, Environment)

### Form Elements
- `form.Method()` - HTTP methods (method package: Get, Post, Dialog)
- `form.Enctype()` - Encoding types (enctype package: URLEncoded, Multipart, TextPlain)
- `button.FormMethod()` - Form submission method (formmethod package: Get, Post)

### Link Elements
- `link.As()` - Resource type hints (as package: Script, Style, Image, Font, Fetch, etc.)
- `link.CrossOrigin()` - CORS settings (crossorigin package: Anonymous, UseCredentials)
- `link.FetchPriority()` - Loading priority (fetchpriority package: High, Low, Auto)
- `link.ReferrerPolicy()` - Referrer policies (referrerpolicy package: NoReferrer, Origin, StrictOrigin, etc.)
- `link.Rel()` - Link relationships (rel package: Stylesheet, Icon, Preload, Prefetch, etc.)

### Image/Media Elements
- `img.Decoding()` - Image decode (decoding package: Sync, Async, Auto)
- `img.Loading()` - Lazy loading (loading package: Lazy, Eager)
- `video.Preload()` - Media preload (preload package: None, Metadata, Auto)

### Global Attributes
- `*.AutoCapitalize()` - Text capitalisation (autocapitalize package: Off, None, On, Sentences, Words, Characters)
- `*.AutoCorrect()` - Auto-correction (autocorrect package: On, Off)
- `*.ContentEditable()` - Editable content (contenteditable package: True, False, PlaintextOnly)
- `*.Dir()` - Text direction (dir package: Ltr, Rtl, Auto)
- `*.EnterKeyHint()` - Virtual keyboard hint (enterkeyhint package: Enter, Done, Go, Next, Previous, Search, Send)
- `*.InputMode()` - Virtual keyboard type (inputmode package: None, Text, Tel, URL, Email, Numeric, Decimal, Search)
- `*.Popover()` - Popover behaviour (popover package: Auto, Manual)
- `*.SpellCheck()` - Spell checking (spellcheck package: True, False)
- `*.Translate()` - Translation hint (translate package: Yes, No)
- `*.VirtualKeyboardPolicy()` - Keyboard behaviour (virtualkeyboardpolicy package: Auto, Manual)
- `*.WritingSuggestions()` - Writing suggestions (writingsuggestions package: True, False)

### Specific Elements
- `meta.Charset()` - Character encoding (charset package: UTF8, ISO88591, Windows1252)
- `ol.ListType()` - List numbering (listtype package: Decimal, LowerAlpha, UpperAlpha, LowerRoman, UpperRoman)
- `area.Shape()` - Image map shape (shape package: Rect, Circle, Poly, Default)
- `iframe.Sandbox()` - Security restrictions (sandbox package: AllowForms, AllowScripts, AllowSameOrigin, etc.)
- `img.Sizes()` - Responsive image sizes (sizes package: predefined breakpoints)
- `button.PopoverTargetAction()` - Popover control (popovertargetaction package: Toggle, Show, Hide)
- `video.CrossOrigin()` - CORS for media (crossorigin package: Anonymous, UseCredentials)

**Usage:**
```go
input.New().
    InputType(inputtype.Email).
    AutoComplete(autocomplete.Email).
    Required()

link.New().
    Rel(rel.Stylesheet).
    Href("/style.css").
    CrossOrigin(crossorigin.Anonymous)

img.New().
    Src("/photo.jpg").
    Loading(loading.Lazy).
    Decoding(decoding.Async)
```

## Dot Import (Convenience Alternative)

For cleaner syntax without package prefixes, the dot import is available as an alternative:

```go
import (
    . "github.com/jpl-au/fluent/dot"
    "github.com/jpl-au/fluent/html5/meta"
)

func render() node.Node {
    return Html(
        Head(
            meta.UTF8(),
            Title().Text("My Page"),
        ),
        Body(
            Div(
                H1().Text("Welcome"),
                P().Text("Hello, world!"),
            ).Class("container"),
        ),
    )
}
```

**Notes:**
- Provides wrapper functions (e.g., `Div()`, `H1()`, `P()`) that call the underlying package constructors
- Specialised constructors like `meta.UTF8()` still require direct package import
- The package-based approach (`div.New()`, `p.Text()`) is the primary API
