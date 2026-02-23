# Fluent — HTML Generation for Go

Fluent is a type-safe, composable HTML generation library for Go. Every HTML element is a Go package (e.g. `div`, `a`, `input`). Elements are constructed with `New()`, configured with chainable methods, and rendered with `Render()`.

## Build & Test

```bash
go build ./...
go test ./...
go vet ./...
```

## Methods That Do NOT Exist

**CRITICAL:** These methods do not exist anywhere in Fluent. Do not use them. Do not hallucinate them. They have never existed.

| Non-existent method | What to use instead |
|---------------------|---------------------|
| `.Attr()` | Use the dedicated typed method (e.g. `.Class()`, `.Href()`, `.Src()`) |
| `.SetAttr()` | Use `.SetAttribute()` for custom attributes only |
| `.Attribute()` | Use the dedicated typed method or `.SetAttribute()` |
| `.Attrs()` | No bulk attribute setter exists — set each attribute individually |
| `.WithAttr()` | Use the dedicated typed method or `.SetAttribute()` |

**The correct approach to setting attributes has three levels:**

1. **Dedicated typed methods (use first)** — Every standard HTML attribute has a chainable method. For example: `.Class()`, `.ID()`, `.Href()`, `.Src()`, `.Alt()`, `.Title()`, `.Disabled()`, `.Required()`, `.Placeholder()`, `.Name()`, `.Value()`, `.Type()`, etc.
2. **SetAria(key, value)** — For ARIA attributes. Automatically adds the `aria-` prefix.
3. **SetData(key, value)** — For data attributes. Automatically adds the `data-` prefix.
4. **SetAttribute(key, value)** — Only for truly custom or non-standard attributes (e.g. Alpine.js directives, HTMX attributes).

```go
// WRONG — these methods do not exist
div.New().Attr("class", "container")           // NO
div.New().SetAttr("id", "main")                // NO
button.New().Attribute("disabled", "")         // NO

// RIGHT — use dedicated typed methods
div.New().Class("container")                   // YES
div.New().ID("main")                           // YES
button.New().Disabled()                        // YES

// RIGHT — use SetAria for ARIA attributes
button.New().SetAria("label", "Close dialog")  // YES — renders aria-label="Close dialog"

// RIGHT — use SetData for data attributes
div.New().SetData("id", "123")                 // YES — renders data-id="123"

// RIGHT — use SetAttribute only for custom/non-standard attributes
div.New().SetAttribute("x-on:click", "handler")  // YES — Alpine.js directive
div.New().SetAttribute("hx-get", "/items")        // YES — HTMX attribute
```

**Important:** `SetAttribute()` does not return the element for chaining. `SetAria()` and `SetData()` do return the element for chaining.

## Core Concepts

### Node and Element Interfaces

The `node.Node` interface is Fluent's foundation. Every renderable piece implements it: HTML elements, text nodes, conditionals (`node.Condition`), and function wrappers (`node.Func`, `node.FuncNodes`).

```go
type Node interface {
    Render(w ...io.Writer) []byte
    RenderBuilder(*bytes.Buffer)
    Nodes() []Node
}
```

HTML elements also implement `node.Element`, which extends `Node` with `SetAttribute()`, `RenderOpen()`, and `RenderClose()`. Text nodes, function components, and conditionals are **not** elements — they don't have attributes or tags.

```go
type Element interface {
    Node
    SetAttribute(key string, value string)
    RenderOpen(buf *bytes.Buffer)
    RenderClose(buf *bytes.Buffer)
}
```

Extensions like fluent-htmx accept `node.Element` rather than `node.Node` because they need to set attributes.

**When in doubt, return `node.Node`** — it's always safe and provides maximum flexibility:

```go
func MyComponent(showHeader bool) node.Node {
    if showHeader {
        return header.New(h1.Text("Welcome"))
    }
    return nil  // nil nodes are safely skipped during rendering
}
```

### Reserved Keyword Alternatives

**CRITICAL:** Some HTML elements use Go reserved keywords. Fluent provides alternative package names:

| HTML Element | Reserved Keyword | Fluent Package | Import Path |
|--------------|------------------|----------------|-------------|
| `<select>`   | `select`         | `dropdown`     | `github.com/jpl-au/fluent/html5/dropdown` |
| `<main>`     | `main`           | `primary`      | `github.com/jpl-au/fluent/html5/primary` |
| `<var>`      | `var`            | `variable`     | `github.com/jpl-au/fluent/html5/variable` |

```go
// CORRECT
dropdown.New(...)  // Renders <select>...</select>
primary.New(...)   // Renders <main>...</main>
variable.New(...)  // Renders <var>...</var>

// WRONG — these packages do not exist
import "github.com/jpl-au/fluent/html5/select"
import "github.com/jpl-au/fluent/html5/main"
import "github.com/jpl-au/fluent/html5/var"
```

### Static vs Text Rendering

**Static()** — Immutable content known at template definition time. JIT-optimisable.
```go
div.Static("Copyright 2024")
```

**Text()** — HTML-escaped dynamic content.
```go
div.Text(user.Name)  // Escaped at runtime
```

**Textf()** — HTML-escaped dynamic content with formatting.
```go
div.Textf("Hello %s, you have %d messages", user.Name, count)
```

**RawText()** — Unescaped HTML content.
```go
div.RawText("<em>Bold</em>")  // Not escaped, use carefully
```

**RawTextf()** — Unescaped HTML content with formatting.
```go
div.RawTextf("<span class=\"%s\">%s</span>", className, content)
```

**Rule:** Use `Static()` for unchanging content (labels, headings, boilerplate). Use `Text()` or `Textf()` for user input or values that change between renders. Use `RawText()` or `RawTextf()` only when you need to inject HTML and trust the source.

### Security Package

`Text()` and `Textf()` use Go's `html.EscapeString()` for basic HTML escaping. For content injected into `<script>` or `<style>` blocks, use the `security` package which detects dangerous patterns:

```go
import "github.com/jpl-au/fluent/security"

security.Sanitise(scriptComponent).Render()   // Returns empty if dangerous
security.Sanitise(scriptComponent).Error()    // Renders error message if invalid
security.SafeScript(jsCode)                   // Sanitised <script> or error comment
security.SafeStyle(cssCode)                   // Sanitised <style> or error comment
```

## Element Construction

### Constructors

All elements follow consistent constructor patterns:

```go
div.New()                              // <div></div>
div.Text("Hello")                      // <div>Hello</div>
div.Static("Footer")                   // <div>Footer</div>
div.RawText("<em>Bold</em>")           // <div><em>Bold</em></div>
div.Textf("Hello %s", name)            // <div>Hello John</div>

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

- `.Text(s)` — adds escaped text content
- `.Textf(format, args...)` — adds formatted escaped text
- `.Static(s)` — adds static text (JIT-optimisable)
- `.RawText(s)` — adds unescaped HTML content
- `.RawTextf(format, args...)` — adds formatted unescaped HTML

```go
div.New().Class("foo").Text("Hello").ID("bar")
p.New().Text("Line 1").Text(" Line 2")  // Multiple text nodes
style.New().RawText("body { color: red; }")
```

### Node Management

- `.Add(nodes...)` — appends child nodes to the element
- `.Replace(nodes...)` — replaces all child nodes with the provided nodes

```go
container := div.New().Class("container")
container.Add(h1.Text("Title"), p.Text("Content"))
container.Replace(span.Text("New content"))
```

## Setting Attributes

### Typed Attribute Methods (Primary API)

Every standard HTML attribute has a dedicated, chainable method on its element. Always use these first.

**Global attributes** (available on all elements):
- `.Class(class)`, `.ID(id)`, `.Style(css)`, `.Title(text)`
- `.Role(role)`, `.Lang(language)`, `.AccessKey(key)`, `.AriaLabel(label)`
- `.Hidden()`, `.TabIndex(index)`, `.AutoFocus()`, `.Draggable()`, `.Inert()`
- `.Nonce(value)`, `.Slot(name)`, `.Is(element)`
- `.AutoCapitalize(value)`, `.AutoCorrect(value)`, `.ContentEditable(value)`
- `.Dir(direction)`, `.EnterKeyHint(hint)`, `.InputMode(mode)`
- `.Popover(value)`, `.SpellCheck(value)`, `.Translate(value)`
- `.VirtualKeyboardPolicy(policy)`, `.WritingSuggestions(value)`
- `.ExportParts(parts)`, `.ItemId(id)`, `.ItemProp(properties)`, `.ItemRef(refs)`, `.ItemScope()`, `.ItemType(itemType)`, `.Part(names)`

**Event attributes** (available on all elements):
- `.OnClick(handler)`, `.OnChange(handler)`, `.OnInput(handler)`
- `.OnFocus(handler)`, `.OnBlur(handler)`, `.OnSubmit(handler)`
- `.OnLoad(handler)`, `.OnError(handler)`, `.OnKeyDown(handler)`, `.OnKeyUp(handler)`
- `.SetEvent(key, value)` — for custom event attributes

**Element-specific attributes** — each element has its own methods. Examples:

| Element | Methods |
|---------|---------|
| `a` | `.Href()`, `.Download()`, `.HrefLang()`, `.Ping()`, `.ReferrerPolicy()`, `.Rel()`, `.Target()`, `.Type()` |
| `img` | `.Src()`, `.Alt()`, `.Width()`, `.Height()`, `.Loading()`, `.Decoding()`, `.Sizes()` |
| `input` | `.Name()`, `.Value()`, `.Placeholder()`, `.InputType()`, `.AutoComplete()`, `.Disabled()`, `.Required()`, `.ReadOnly()`, `.Multiple()`, `.Checked()`, `.MaxLength()`, `.Accept()`, `.Capture()` |
| `form` | `.Action()`, `.Method()`, `.Enctype()`, `.Target()` |
| `link` | `.Href()`, `.Rel()`, `.As()`, `.CrossOrigin()`, `.FetchPriority()`, `.Media()` |
| `button` | `.Disabled()`, `.FormAction()`, `.FormMethod()`, `.PopoverTarget()`, `.PopoverTargetAction()` |
| `script` | `.Src()`, `.Async()`, `.Defer()`, `.Type()`, `.CrossOrigin()`, `.Integrity()` |

If a standard HTML attribute has a method on the element, use that method. Do not use `SetAttribute()` for standard attributes.

### SetAria — ARIA Attributes

`SetAria(key, value)` sets ARIA attributes. It automatically adds the `aria-` prefix. Returns the element for chaining.

```go
button.New().SetAria("label", "Close dialog")
// Renders: <button aria-label="Close dialog"></button>

div.New().SetAria("hidden", "true")
// Renders: <div aria-hidden="true"></div>

nav.New().SetAria("expanded", "false").SetAria("controls", "menu")
// Renders: <nav aria-expanded="false" aria-controls="menu"></nav>
```

Note: `.AriaLabel()` is a convenience method equivalent to `.SetAria("label", value)`. For all other ARIA attributes, use `SetAria()`.

### SetData — Data Attributes

`SetData(key, value)` sets data attributes. It automatically adds the `data-` prefix. Returns the element for chaining.

```go
div.New().SetData("id", "123")
// Renders: <div data-id="123"></div>

button.New().SetData("action", "submit").SetData("confirm", "true")
// Renders: <button data-action="submit" data-confirm="true"></button>
```

### SetAttribute — Custom/Non-Standard Only

`SetAttribute(key, value)` sets arbitrary attributes. **Does NOT return the element** (it satisfies the `node.Element` interface). Use only for custom or non-standard attributes.

```go
div.New().SetAttribute("x-on:click", "handler")  // Alpine.js
div.New().SetAttribute("hx-get", "/items")        // HTMX
div.New().SetAttribute("custom-attr", "value")    // Custom
```

**Do not use SetAttribute for:**
- Standard HTML attributes (use the typed method instead)
- ARIA attributes (use `SetAria()` instead)
- Data attributes (use `SetData()` instead)

### Type-Safe Constants

Fluent uses typed constants for attributes with enumerated values. Methods accept a typed constant, not a string — so typos cause compile errors.

```go
// Typed constant — compile error on typo
input.New().InputType(inputtype.Email)

// Escape hatch for edge cases
input.New().InputType(inputtype.Custom("future-type"))
```

Each attribute package provides a `Custom()` function for values not yet covered by predefined constants.

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
```

## Dynamic Content

### Conditional Rendering

`node.Condition()` provides inline conditional rendering:

```go
node.Condition(user.IsLoggedIn).
    True(div.Text("Welcome back!")).
    False(div.Text("Please log in"))
```

Shorthand forms:

```go
node.When(user.IsAdmin, span.Static("Admin"))                        // Render when true
node.Unless(user.IsLoggedIn, a.New().Href("/login").Text("Sign in")) // Render when false
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
- `Condition(bool).True(node).False(node)` — both branches
- `When(bool, node)` — shorthand for true-only branch
- `Unless(bool, node)` — shorthand for false-only branch
- Nil nodes are safely ignored

### Function Components

**Single node** — `node.Func()`:

```go
node.Func(func() node.Node {
    if user.Role == "admin" {
        return div.Text("Admin Panel")
    }
    return div.Text("User Dashboard")
})
```

**Multiple nodes** — `node.FuncNodes()`:

```go
node.FuncNodes(func() []node.Node {
    nodes := []node.Node{}
    for _, item := range items {
        nodes = append(nodes, li.Text(item.Name))
    }
    return nodes
})
```

### Dynamic Interface

All nodes implement the `node.Dynamic` interface, which reports whether a node produces different output across renders:

```go
type Dynamic interface {
    IsDynamic() bool
    DynamicKey() string
}
```

**`IsDynamic()`** returns `true` if the node's output may change between renders:
- `Text()`, `RawText()`, `Textf()`, `RawTextf()` nodes — always dynamic
- `Static()` nodes — never dynamic
- `node.Condition`, `node.Func`, `node.FuncNodes` — always dynamic
- HTML elements — dynamic if marked with `.Dynamic()` or if any child is dynamic

**`DynamicKey()`** returns the developer-assigned key for reactive tracking, or an empty string if unset.

### Reactive Tracking with .Dynamic()

The `.Dynamic(key)` method on HTML elements marks them for reactive tracking by the [Fluent Poly](#ecosystem) diff engine. The key identifies the element across renders so the diff engine can detect changes and send targeted DOM patches.

```go
// Mark an element for reactive tracking
span.Textf("Count: %d", state.Count).Dynamic("count")

// The key renders as a data-poly-key attribute
// <span data-poly-key="count">Count: 42</span>

// Keys must be unique within a render tree
p.Text(state.ErrorMsg).Dynamic("error-message")
table.New(rows...).Dynamic("data-table")

// Mark as dynamic without a tracking key (rare)
div.New(children...).Dynamic()
```

`.Dynamic()` is chainable and follows the same pattern as `.Class()`, `.SetData()`, etc. It is used by [Fluent JIT](https://github.com/jpl-au/fluent-jit) to identify which segments need re-evaluation and by [Fluent Poly](https://github.com/jpl-au/fluent-poly) for targeted DOM patching over WebSocket.

Elements without `.Dynamic()` are not tracked — the diff engine only examines keyed nodes.

## Component Pattern

### Return Types: Interface vs Concrete

**`node.Node` (interface)** — use when the function may return different element types, or the component is a final building block.

**`*element.Element` (concrete type, e.g. `*div.Element`)** — use when callers should be able to chain additional methods.

```go
// Return node.Node — flexible, no chaining after call
func Card(title, content string) node.Node {
    return div.New(
        h2.Text(title),
        p.Text(content),
    ).Class("card")
}

// Return *div.Element — allows continued chaining
func Card(title, content string) *div.Element {
    return div.New(
        h2.Text(title),
        p.Text(content),
    ).Class("card")
}

// Caller can chain: Card("Hi", "Hello!").ID("welcome").Class("highlighted")
```

All element packages export their concrete type as `Element`: `*div.Element`, `*a.Element`, `*input.Element`, etc.

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
            span.Textf("$%.2f", p.Price),
        )
    }
    return ul.New(items...)
}
```

## Performance

- Use `Static()` for unchanging content (enables JIT optimisation)
- Buffer pooling is enabled by default and handled automatically
- Each element has a `BufferHint()` method for optional buffer size hints
- For high-throughput applications, [Fluent JIT](https://github.com/jpl-au/fluent-jit) provides additional optimisation (Compile, Tune, Flatten)

## Extending Fluent

Implement `node.Node` for composite components, or `node.Element` for custom HTML elements that need attributes and open/close tags.

### Composite Component Example

A composite component composes elements internally. It satisfies `node.Node` (not `node.Element`) because it doesn't have its own tag or attributes.

```go
type EmailField struct {
    labelText   string
    name        string
    placeholder string
    required    bool
    class       string
}

func Email(name string, labelText string) *EmailField {
    return &EmailField{labelText: labelText, name: name}
}

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
    labelElem := label.For(f.name, f.labelText)
    inputElem := input.Email(f.name).
        ID(f.name).
        AutoComplete(autocomplete.Email)

    if f.placeholder != "" {
        inputElem.Placeholder(f.placeholder)
    }
    if f.required {
        inputElem.Required()
    }

    container := div.New(labelElem, inputElem)
    if f.class != "" {
        container.Class(f.class)
    }
    container.RenderBuilder(buf)
}

func (f *EmailField) Nodes() []node.Node {
    return nil
}
```

## Typed Attributes Reference

Elements with typed attribute constants:

### Input Elements
- `input.Accept()` — File types (accept: ImageJPEG, ImagePNG, VideoMP4, AudioMP3, Pdf, Docx, etc.)
- `input.AutoComplete()` — Autocomplete hints (autocomplete: On, Off, Name, Email, Username, etc.)
- `input.InputType()` — Input types (inputtype: Text, Email, Password, Number, Tel, URL, etc.)
- `input.Capture()` — Camera capture (capture: User, Environment)

### Form Elements
- `form.Method()` — HTTP methods (method: Get, Post, Dialog)
- `form.Enctype()` — Encoding types (enctype: URLEncoded, Multipart, TextPlain)
- `button.FormMethod()` — Form submission method (formmethod: Get, Post)

### Link Elements
- `link.As()` — Resource type hints (as: Script, Style, Image, Font, Fetch, etc.)
- `link.CrossOrigin()` — CORS settings (crossorigin: Anonymous, UseCredentials)
- `link.FetchPriority()` — Loading priority (fetchpriority: High, Low, Auto)
- `link.ReferrerPolicy()` — Referrer policies (referrerpolicy: NoReferrer, Origin, StrictOrigin, etc.)
- `link.Rel()` — Link relationships (rel: Stylesheet, Icon, Preload, Prefetch, etc.)

### Image/Media Elements
- `img.Decoding()` — Image decode (decoding: Sync, Async, Auto)
- `img.Loading()` — Lazy loading (loading: Lazy, Eager)
- `video.Preload()` — Media preload (preload: None, Metadata, Auto)

### Global Attributes (on all elements)
- `*.AutoCapitalize()` — (autocapitalize: Off, None, On, Sentences, Words, Characters)
- `*.AutoCorrect()` — (autocorrect: On, Off)
- `*.ContentEditable()` — (contenteditable: True, False, PlaintextOnly)
- `*.Dir()` — (dir: Ltr, Rtl, Auto)
- `*.EnterKeyHint()` — (enterkeyhint: Enter, Done, Go, Next, Previous, Search, Send)
- `*.InputMode()` — (inputmode: None, Text, Tel, URL, Email, Numeric, Decimal, Search)
- `*.Popover()` — (popover: Auto, Manual)
- `*.SpellCheck()` — (spellcheck: True, False)
- `*.Translate()` — (translate: Yes, No)
- `*.VirtualKeyboardPolicy()` — (virtualkeyboardpolicy: Auto, Manual)
- `*.WritingSuggestions()` — (writingsuggestions: True, False)

### Specific Elements
- `meta.Charset()` — (charset: UTF8, ISO88591, Windows1252)
- `ol.ListType()` — (listtype: Decimal, LowerAlpha, UpperAlpha, LowerRoman, UpperRoman)
- `area.Shape()` — (shape: Rect, Circle, Poly, Default)
- `iframe.Sandbox()` — (sandbox: AllowForms, AllowScripts, AllowSameOrigin, etc.)
- `img.Sizes()` — (sizes: predefined breakpoints)
- `button.PopoverTargetAction()` — (popovertargetaction: Toggle, Show, Hide)

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

## Ecosystem

Fluent has companion packages that extend its capabilities. All are optional — Fluent works standalone for static HTML generation.

| Package | Description |
|---------|-------------|
| [Fluent JIT](https://github.com/jpl-au/fluent-jit) | Performance optimisation. **Compile** pre-renders static portions and re-evaluates dynamic content. **Tune** provides adaptive buffer sizing. **Flatten** pre-renders fully static content to raw bytes. Also provides the **Diff** engine that compares renders by dynamic key and produces `[]Patch`. |
| [Fluent HTMX](https://github.com/jpl-au/fluent-htmx) | HTMX integration. Accepts `node.Element` to set HTMX attributes (`hx-get`, `hx-post`, `hx-swap`, etc.) on any Fluent element. |
| [Fluent Poly](https://github.com/jpl-au/fluent-poly) | Server-driven reactive UI. Manages sessions, WebSocket transport, and a client-side runtime that applies targeted DOM patches. Mark elements with `.Dynamic("key")` and Poly handles diffing, patching, and event handling. Uses the JIT diff engine internally. |

**How the packages relate:**
- **fluent** (this package) — core HTML generation, `node.Node`/`node.Element` interfaces, `.Dynamic()` method
- **fluent-jit** — rendering optimisation, diff engine produces `[]Patch` from two tree states
- **fluent-poly** — connection lifecycle, calls `jit.Differ.Diff()` and sends patches over WebSocket
- **fluent-htmx** — attribute wrapper for HTMX, independent of JIT and Poly

## Dot Import (Convenience Alternative)

For cleaner syntax without package prefixes:

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

The package-based approach (`div.New()`, `p.Text()`) is the primary API. Specialised constructors like `meta.UTF8()` still require direct package import.

## Profile-Guided Optimization (PGO)

Applications using Fluent benefit from [PGO](https://go.dev/doc/pgo) (Go 1.21+). Collect a CPU profile from production, place it as `default.pgo` in the main package, and `go build` applies it automatically. Expect 10-20% speed improvements across the rendering pipeline with no code changes. Allocations are unaffected — PGO improves inlining decisions only. Collect fresh profiles periodically as code evolves.
