# Fluent - HTML Generation for Go

Fluent is a type-safe, composable HTML generation library for Go. Every HTML element is a Go package (e.g. `div`, `a`, `input`). Elements are constructed with `New()`, configured with chainable methods, and rendered with `Render()`.

## Build & Test

```bash
go build ./...
go test ./...
go vet ./...
```

## Lint With Flint

Run Flint after you write or edit Fluent code. It catches misuse of Fluent that `go vet` cannot see. Every diagnostic includes a `fix:` field with the corrected code. Read the fix, apply it, and run Flint again until it reports nothing.

```bash
go install github.com/jpl-au/flint/cmd/flint@latest
flint ./...
```

What Flint flags:

- **Hallucinated APIs.** Every function, method, and type reference is validated against a generated registry. Catches calls like `node.Fragment()`, `div.New().Href(...)` (wrong element), `inputtype.Telephone` (does not exist).
- **Raw strings where typed constants are required.** `input.New().Type("email")` is flagged with a fix pointing at `inputtype.Email`. See the "Never use raw strings where typed constants exist" section below for the full rule.
- **Unsafe `Static()` and `RawText()`.** `Static()` requires a string literal so the JIT can pre-render it. `RawText()` does not HTML-escape. Both are flagged when called with a variable.
- **`New().Method()` redundancy.** Flags `div.New().Text("x")` and suggests `div.Text("x")` directly. See "Choosing the Right Constructor" below.
- **Typed constructor opportunities.** Suggests `ul.Items(...)` over `ul.New(li...)`, `tr.Cells(...)` over `tr.New(td...)`, and similar type-safe alternatives when children are uniform.
- **`SetAttribute()` misuse.** Flags chaining after `SetAttribute()` (it returns void), and flags usage where a typed method exists (`.Class()` instead of `.SetAttribute("class", ...)`).
- **Reserved keyword imports.** `select` → `dropdown`, `main` → `primary`, `var` → `variable`.

### Look the element up

Every constructor, method, typed parameter and constant in Fluent lives in a registry generated from the same specs as the packages, so the registry cannot drift from the code. Read it with `-info` before writing against an element. This file does not carry a constructor or constant reference, because `-info` answers both on demand and stays correct.

```bash
flint -info div                 # everything about <div>
flint -info input ctors attrs   # just constructors and attribute mappings
flint -info ol typed            # typed constructors only
```

Sections: `types`, `constructors`/`ctors`, `typed-constructors`/`typed`, `methods`, `attributes`/`attrs`, `vars`, `elements`.

A method that takes a typed constant prints the package holding its values, as `Type  (enum: inputtype)`. The output names the package, not the values. Query that package the same way to read them:

```bash
flint -info inputtype    # Vars: Button, Checkbox, Color, Date, Email, ...
flint -info loading      # Vars: Eager, Lazy
```

An element named after a Go reserved word resolves from its HTML name, so a wrong guess still lands on the right package:

```bash
flint -info select       # Element: dropdown (renders <select>; "select" is a Go reserved word)
flint -info main         # Element: primary (renders <main>; "main" is a Go reserved word)
flint -info var          # Element: variable (renders <var>; "var" is a Go reserved word)
```

Exit codes: `0` clean, `1` diagnostics found, `2` usage or I/O error. Treat `1` as "fix and re-run", not "done".

## Methods that do not exist

These methods do not exist in Fluent. Do not use them.

| Non-existent method | What to use instead |
|---------------------|---------------------|
| `.Attr()` | Use the dedicated typed method (e.g. `.Class()`, `.Href()`, `.Src()`) |
| `.SetAttr()` | Use `.SetAttribute()` for custom attributes only |
| `.Attribute()` | Use the dedicated typed method or `.SetAttribute()` |
| `.Attrs()` | No bulk attribute setter exists - set each attribute individually |
| `.WithAttr()` | Use the dedicated typed method or `.SetAttribute()` |

**The correct approach to setting attributes has four levels:**

1. **Dedicated typed methods (use first)** - Every standard HTML attribute has a chainable method. For example: `.Class()`, `.ID()`, `.Href()`, `.Src()`, `.Alt()`, `.Title()`, `.Disabled()`, `.Required()`, `.Placeholder()`, `.Name()`, `.Value()`, `.Type()`, etc.
2. **SetAria(key, value)** - For ARIA attributes. Automatically adds the `aria-` prefix.
3. **SetData(key, value)** - For data attributes. Automatically adds the `data-` prefix.
4. **SetAttribute(key, value)** - Only for truly custom or non-standard attributes (e.g. Alpine.js directives, HTMX attributes).

In every key-value setter the key is written to the rendered output verbatim. The key is code, not data: pass a fixed, developer-controlled key, and never build it from user input. A key containing a space, quote, `=`, `/` or `>` changes the markup structure. The value is escaped; the key is not.

```go
// These methods do not exist in Fluent, so none of these lines compile.
div.New().Attr("class", "container")           // wrong
div.New().SetAttr("id", "main")                // wrong
button.New().Attribute("disabled", "")         // wrong

// Use the dedicated typed method for a standard attribute.
div.New().Class("container")                   // correct
div.New().ID("main")                           // correct
button.New().Disabled()                        // correct

// Use SetAria for ARIA attributes.
button.New().SetAria("label", "Close dialog")  // correct - renders aria-label="Close dialog"

// Use SetData for data attributes.
div.New().SetData("id", "123")                 // correct - renders data-id="123"

// Use SetAttribute only for custom or non-standard attributes.
div.New().SetAttribute("x-on:click", "handler")  // correct - Alpine.js directive
div.New().SetAttribute("hx-get", "/items")        // correct - HTMX attribute
```

`SetAttribute()` does not return the element, so you cannot chain after it. `SetAria()` and `SetData()` do return the element.

## Common Mistakes

These are the most frequent errors. Each one fails to compile.

### 1. Never use raw strings where typed constants exist

Fluent uses `[]byte`-based types for enumerated HTML attribute values. The method does not accept a plain string. Import the constant package and use its exported variable.

```go
// A raw string does not compile here.
img.New().Loading("eager")              // cannot use "eager" (untyped string constant) as loading.Loading
img.New().Loading("lazy")               // same error
input.New().Type("email")               // cannot use "email" as inputtype.InputType

// Use the typed constant instead.
img.New().Loading(loading.Eager)        // import "github.com/jpl-au/fluent/html5/attr/loading"
img.New().Loading(loading.Lazy)
input.New().Type(inputtype.Email)        // import "github.com/jpl-au/fluent/html5/attr/inputtype"
```

Every attribute package also has a `Custom()` escape hatch for values not yet covered:
```go
loading.Custom("future-value")
```

Import path pattern: `github.com/jpl-au/fluent/html5/attr/<package>`. Run `flint -info <package>` to read the values a package holds, as in `flint -info loading`.

### 2. Use element-specific constructors - do not build common patterns by hand

Most elements have constructors that set several attributes at once. Use the constructor instead of chaining the attributes yourself.

```go
// Building these by hand is manual, verbose and error-prone.
meta.New().SetAttribute("charset", "UTF-8")                    // SetAttribute returns void, not chainable
meta.New().SetAttribute("name", "viewport")                    // also wrong approach entirely
img.New().Src("photo.jpg").Alt("Sunset").Loading(loading.Lazy) // works but unnecessary

// Use the constructor that already does it.
meta.UTF8()                          // <meta charset="UTF-8">
meta.Viewport("width=device-width") // <meta name="viewport" content="width=device-width">
img.Lazy("photo.jpg", "Sunset")     // <img src="photo.jpg" alt="Sunset" loading="lazy">
```

Beyond the universal constructors (`New`, `Text`, `Static`, `RawText`, `Textf`, `RawTextf`) that every element has, many elements provide domain-specific constructors. Always prefer these over manual attribute chaining.

A **typed constructor** accepts only the valid child element type, so the compiler checks the nesting. Use one whenever the children are all the same element type. `New(...node.Node)` also compiles there, because `New` accepts any child, so the nesting check is lost.

```go
// New accepts any child, so nothing here catches a <div> put inside a <ul>.
ul.New(li.Text("one"), li.Text("two"))
tr.New(td.Text("Alice"), td.Text("30"))
table.New(tr.Cells(td.Text("Alice")))

// These reject a child of the wrong type at compile time.
ul.Items(li.Text("one"), li.Text("two"))
tr.Cells(td.Text("Alice"), td.Text("30"))
table.Rows(tr.Cells(td.Text("Alice")))
```

Keep `New(...node.Node)` for mixed or dynamic content, where the children genuinely have no single element type.

**Cross-package constructors** create common child elements for you. `details.Summary("label", nodes...)` creates a `<summary>` child automatically.

Run `flint -info <element> ctors` to read what an element provides.

### 3. Do not hallucinate node-level helpers

These do not exist: `node.StaticText`, `node.RawNode`, `node.TextNode`, `node.HTML`. The helpers that do exist are:

| Exists | Usage |
|--------|-------|
| `node.Func(func() node.Node)` | Single dynamic node |
| `node.Funcs(func() []node.Node)` | Multiple dynamic nodes |
| `node.Condition(bool).True(n).False(n)` | Conditional rendering |
| `node.When(bool, node)` | Render when true |
| `node.Unless(bool, node)` | Render when false |
| `node.Empty()` | A node that renders nothing - the safe, explicit "render nothing" (prefer it over returning a nil, especially a typed nil pointer) |

Text content is always a method on the element: `div.Text(...)`, `div.Static(...)`, `div.RawText(...)`.

### 4. Do not build a typed list through node.Map

`node.Map` and `node.Funcs` return `node.Node`, so `ul.New(node.Map(...))` compiles and then accepts any child. Use the deferred typed constructor and declare the mapper to return the child element type. Each collection constructor has one: `ItemsOf` beside `Items`, `RowsOf` beside `Rows`, `CellsOf` and `HeadersOf` beside `Cells` and `Headers`, `OptionsOf` beside `Options`, and `ColsOf` beside `Cols`.

```go
ul.ItemsOf(products, func(p Product) *li.Element { ... })
```

---

## Core Concepts

### Node and Element Interfaces

The `node.Node` interface is Fluent's foundation. Every renderable piece implements it: HTML elements, text nodes, conditionals (`node.Condition`), and function wrappers (`node.Func`, `node.Funcs`).

```go
type Node interface {
    Render(w io.Writer)                // fire-and-forget; write errors discarded
    WriteTo(w io.Writer) (int64, error) // write and observe the error
    RenderBytes() []byte               // render to a new []byte
    RenderBuilder(*bytes.Buffer)
    Nodes() []Node
}
```

`Nodes()` returns the children that Fluent renders. For a conditional, this is the active branch only. For a function component, this is the output of the function. Tree walkers, such as the Differ's snapshot collector, rely on `Nodes()` matching what is rendered.

HTML elements also implement `node.Element`, which extends `Node` with `SetAttribute()`, `SetAttributeRaw()`, `RenderOpen()`, and `RenderClose()`. Text nodes, function components, and conditionals are **not** elements - they don't have attributes or tags.

```go
type Element interface {
    Node
    SetAttribute(key string, value string)    // value is HTML-escaped at set time
    SetAttributeRaw(key string, value string) // value stored verbatim (trusted values only)
    RenderOpen(buf *bytes.Buffer)
    RenderClose(buf *bytes.Buffer)
}
```

Extensions like fluent-htmx accept `node.Element` rather than `node.Node` because they need to set attributes.

**When in doubt, return `node.Node`** - it's always safe and provides maximum flexibility:

```go
func MyComponent(showHeader bool) node.Node {
    if showHeader {
        return header.New(h1.Text("Welcome"))
    }
    return nil  // nil nodes are safely skipped during rendering
}
```

### Reserved Keyword Alternatives

Some HTML element names are Go reserved keywords. Fluent gives those packages a different name:

| HTML Element | Reserved Keyword | Fluent Package | Import Path |
|--------------|------------------|----------------|-------------|
| `<select>`   | `select`         | `dropdown`     | `github.com/jpl-au/fluent/html5/dropdown` |
| `<main>`     | `main`           | `primary`      | `github.com/jpl-au/fluent/html5/primary` |
| `<var>`      | `var`            | `variable`     | `github.com/jpl-au/fluent/html5/variable` |

```go
// Use the Fluent package name.
dropdown.New(...)  // Renders <select>...</select>
primary.New(...)   // Renders <main>...</main>
variable.New(...)  // Renders <var>...</var>

// These packages do not exist, so these imports fail.
import "github.com/jpl-au/fluent/html5/select"
import "github.com/jpl-au/fluent/html5/main"
import "github.com/jpl-au/fluent/html5/var"
```

### Static vs Text Rendering

**Static()** - Immutable content known at template definition time. JIT-optimisable.
```go
div.Static("Copyright 2024")
```

**Text()** - HTML-escaped dynamic content.
```go
div.Text(user.Name)  // Escaped at runtime
```

**Textf()** - HTML-escaped dynamic content with formatting.
```go
div.Textf("Hello %s, you have %d messages", user.Name, count)
```

**RawText()** - Unescaped HTML content.
```go
div.RawText("<em>Bold</em>")  // Not escaped, use carefully
```

**RawTextf()** - Unescaped HTML content with formatting.
```go
div.RawTextf("<span class=\"%s\">%s</span>", className, content)
```

Use `Static()` for content that does not change, such as a label or a heading. Use `Text()` or `Textf()` for user input, and for any value that changes between renders. Use `RawText()` or `RawTextf()` only to insert HTML from a source you trust.

### Sanitising untrusted HTML

`Text()` and `Textf()` escape HTML with `html.EscapeString()`. That is enough for plain text. To render untrusted content *as* HTML, such as rendered markdown, rich-text input or a comment body, use the companion package [fluent-security](https://github.com/jpl-au/fluent-security):

```go
import "github.com/jpl-au/fluent-security"

// HTML - permissive UGC baseline. Keeps common formatting tags, strips
// scripts, event handlers, and dangerous URIs.
div.New(security.HTML(userMarkdown)).Class("body")

// PlainText - strips all tags, keeps only escaped text content.
h1.New(security.PlainText(article.RawTitle))
```

To reuse a policy, declare a `*security.Cleaner` at package scope. Use `security.RichText()` for the user-content baseline, or `security.New()` for an empty allowlist. Extend either with `Allow`, `AllowClasses` and `AllowAttr`.

fluent-security also provides `Nonce()` for Content-Security-Policy work. Pass the nonce to `script.Nonce(...)` or `style.Nonce(...)` and match it against your CSP header.

Fluent core does not sanitise HTML. Sanitising is a separate package that you opt in to.

## Element Construction

### Choosing the Right Constructor

Every element has six constructors: `New`, `Text`, `Static`, `RawText`, `Textf` and `RawTextf`. Pick the one that matches the element's main content. Use `New()` when the element holds child nodes, or holds nothing.

**Rule: if the element's content is text, use the text constructor.**

```go
// text constructor: the element's content is text
h1.Text("Welcome")                        // <h1>Welcome</h1>
p.Textf("Hello %s", name)                 // <p>Hello John</p>
span.Static("Copyright 2024")             // <span>Copyright 2024</span>
style.RawText("body { color: red; }")     // <style>body { color: red; }</style>

// New(): the element contains child nodes
div.New(
    p.Text("Paragraph"),
    span.Text("Inline"),
)

// New(): svg is a typed shape container; its children are svg.Shape values
svg.New(
    svg.Rect().X("0").Width("40").Fill("var(--blue)"),
    svg.Circle().Cx("60").R("50"),
)                                          // the svg package is the shape vocabulary

// New(): the element has no content, only attributes
div.New().Class("container")               // <div class="container"></div>
```

**When to use `New()`:**
- The element is a **container for child nodes**: `div.New(h1.Text("Title"), p.Text("Body"))`
- The element is **empty** with only attributes: `div.New().Class("spacer")`
- You need to **add child nodes and text**: `div.New(p.Text("child")).Class("wrapper")`

**When not to use `New()`:**
- The element's content is **just text** - use `Text()`, `Textf()`, `Static()`
- The element's content is **raw HTML** - use `RawText()`, `RawTextf()`

Attributes can be chained on any constructor. The constructor choice is about content, not attributes:
```go
h1.Text("Welcome").Class("title").ID("hero")  // text constructor + chained attributes
div.New(child1, child2).Class("grid")          // node constructor + chained attributes
```

### Content Methods (for adding to an existing element)

Every non-self-closing element also has these as chainable methods. Use a method only to add content to an element you built with `New()` for its child nodes:

- `.Text(s)` - adds escaped text content
- `.Textf(format, args...)` - adds formatted escaped text
- `.Static(s)` - adds static text (JIT-optimisable)
- `.RawText(s)` - adds unescaped HTML content
- `.RawTextf(format, args...)` - adds formatted unescaped HTML

```go
// Method form - adding text alongside child nodes
div.New(img.Src("photo.jpg")).Text("Caption")

// Multiple text nodes via method chaining
p.Text("Line 1").Text(" Line 2")
```

### Node Management

- `.Add(nodes...)` - appends child nodes to the element
- `.Replace(nodes...)` - replaces all child nodes with the provided nodes

```go
container := div.New().Class("container")
container.Add(h1.Text("Title"), p.Text("Content"))
container.Replace(span.Text("New content"))
```

## Setting Attributes

### Typed Attribute Methods (Primary API)

Every standard HTML attribute has a dedicated, chainable method on its element. Always use these first. If a standard HTML attribute has a method on the element, use that method and do not reach for `SetAttribute()`.

Three sets of methods sit on every element: the global attributes such as `.Class()`, `.ID()` and `.Style()`; the event attributes such as `.OnClick()` and `.OnSubmit()`, with `.SetEvent(key, value)` for a custom one; and the element's own attributes, such as `.Href()` on `a` or `.Src()` on `img`.

Read the list rather than recalling it:

```bash
flint -info a methods    # every method on <a>
flint -info a attrs      # HTML attribute name -> method name
```

### Boolean Attributes

A boolean HTML attribute method takes an optional `bool`. With no argument, the method sets the attribute. With one `bool`, the method sets the attribute only when the value is true. Use this form for a conditional attribute instead of building the element twice.

```go
// No argument - attribute always present
input.Email("email").Required()

// Inline condition - attribute present only when true
input.New().Checked(slices.Contains(selected, c.ID))
button.Text("Save").Disabled(!form.Valid)
script.Src("/app.js").Async(isProduction)
textarea.New().ReadOnly(user.IsGuest)
```

The pattern is consistent across every boolean attribute in Fluent, so it holds for whichever one `flint -info <element> methods` shows you.

`.Hidden()` is the exception: it accepts a typed `hidden.Hidden` value (to support the `until-found` state), not a `bool`.

The inline form composes naturally inside `node.Map` where a per-item flag toggles the attribute - see [Conditional Rendering](#conditional-rendering) and [Function Components](#function-components).

### SetAria - ARIA Attributes

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

### SetData - Data Attributes

`SetData(key, value)` sets data attributes. It automatically adds the `data-` prefix. Returns the element for chaining.

```go
div.New().SetData("id", "123")
// Renders: <div data-id="123"></div>

button.New().SetData("action", "submit").SetData("confirm", "true")
// Renders: <button data-action="submit" data-confirm="true"></button>
```

### SetAttribute - Custom/Non-Standard Only

`SetAttribute(key, value)` sets an arbitrary attribute. It does not return the element (it satisfies the `node.Element` interface). Use only for custom or non-standard attributes.

```go
div.New().SetAttribute("x-on:click", "handler")  // Alpine.js
div.New().SetAttribute("hx-get", "/items")        // HTMX
div.New().SetAttribute("custom-attr", "value")    // Custom
```

**Do not use SetAttribute for:**
- Standard HTML attributes (use the typed method instead)
- ARIA attributes (use `SetAria()` instead)
- Data attributes (use `SetData()` instead)

### Attribute Escaping and URL Filtering

Fluent escapes attribute values for you. Do not escape them yourself.

- **Constructors and setters escape the value.** This applies to typed methods such as `.Class()` and `.Title()`, and to `SetAttribute()`, `SetData()`, `SetAria()`, `.Dynamic()` keys and enum `Custom()` values.
- **Do not escape a value before you pass it.** `html.EscapeString` before a setter escapes the value twice. Pass the value as you received it.
- **`SetAttributeRaw(key, value)` does not escape.** It stores the value as given. Use it only for a value you trust. It is the attribute equivalent of `RawText`.
- **Some URL attributes accept a fixed set of schemes.** `href`, `src`, `action`, `formaction` and `<object>` `data` accept `http`, `https`, `mailto`, `tel`, `sms` and relative URLs. Fluent replaces any other scheme, such as `javascript:`, with `node.UnsafeURL`. Call `node.OnUnsafeURL(fn)` to see which values Fluent replaced.
- **Image and media URLs are escaped but not scheme-checked.** This covers `src` on `img` and `video`, and `srcset` and `poster`. A `data:image/...` URL still works.

### Type-Safe Constants

Fluent uses typed constants for attributes with enumerated values. Methods accept a typed constant, not a string - so typos cause compile errors.

```go
// Typed constant - compile error on typo
input.New().Type(inputtype.Email)

// Escape hatch for edge cases
input.New().Type(inputtype.Custom("future-type"))
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
node.Unless(user.IsLoggedIn, a.Link("/login", "Sign in")) // Render when false
```

For multiple branches, `node.Func()` is cleaner than deeply nested conditions:

```go
node.Func(func() node.Node {
    if !user.IsLoggedIn {
        return a.Link("/login", "Sign in")
    }
    if user.IsAdmin {
        return span.Static("Admin Dashboard")
    }
    return span.Static("User Dashboard")
})
```

**Summary:**
- `Condition(bool).True(node).False(node)` - both branches
- `When(bool, node)` - shorthand for true-only branch
- `Unless(bool, node)` - shorthand for false-only branch
- Nil nodes are safely ignored

### Function Components

**Single node** - `node.Func()`:

```go
node.Func(func() node.Node {
    if user.Role == "admin" {
        return div.Text("Admin Panel")
    }
    return div.Text("User Dashboard")
})
```

**Multiple nodes** - `node.Funcs()`:

```go
node.Funcs(func() []node.Node {
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
- `Text()`, `RawText()`, `Textf()`, `RawTextf()` nodes - always dynamic
- `Static()` nodes - never dynamic
- `node.Condition`, `node.Func`, `node.Funcs` - always dynamic
- HTML elements - dynamic if marked with `.Dynamic()` or if any child is dynamic

**`DynamicKey()`** returns the developer-assigned key for reactive tracking, or an empty string if unset.

### Reactive Tracking with .Dynamic()

`.Dynamic(key)` marks an element for an external render engine, such as [Fluent JIT](https://github.com/jpl-au/fluent-jit). Fluent stores the key and renders it as the element's `id` attribute, so the diff identity and the DOM identity are the same attribute. It does nothing else with it. An explicit `.ID()` on the same element must equal the key; a different value panics at render. The method is chainable, like `.Class()` and `.SetData()`.

Elements also carry `.Memoise(version)` and `.MemoiseKey()` for the same engines. Plain rendering ignores all three methods. Read the engine's own documentation for how to use them.

## Components

A component is a Go function returning `node.Node` or a concrete element type. There is no component system to learn.

### Take the record

Where a component draws something the application already has a type for, take that type. The component then decides how the record is drawn, including what an empty list looks like, and the caller passes what it loaded.

```go
type Product struct {
    ID    string
    Name  string
    Price float64
}

// One record, always present. Returns *li.Element so a typed constructor
// can check the nesting.
func Row(p Product) *li.Element {
    return li.New(
        span.Text(p.Name),
        span.Textf("$%.2f", p.Price),
    ).Class("product-row")
}

// Many records. The component owns the empty case.
func List(products []Product) node.Node {
    if len(products) == 0 {
        return div.Text("No products yet.").Class("empty")
    }
    return ul.ItemsOf(products, Row).Class("products")
}
```

`List(products)` renders the list or the empty message, so the caller writes no condition.

```go
// The caller is deciding what an empty list looks like. That is the
// component's job.
if len(products) == 0 {
    return div.Text("No products yet.")
}
return ProductList(products)
```

A condition written at the call site is repeated at every other call site, and the copies stop agreeing.

### Lists from data

Each collection constructor (`Items`, `Rows`, `Cells`, `Headers`, `Options`, `Cols`) has two deferred siblings. Styled variants such as `ol.Decimal` and `optgroup.Labelled` do not. Use the one that matches what you hold:

| You have | Use | Example |
|----------|-----|---------|
| the elements themselves | `Items` | `ul.Items(a, b, c)` |
| data and a mapper | `ItemsOf` | `ul.ItemsOf(products, Row)` |
| logic that reorders or computes the list | `ItemsFunc` | `ul.ItemsFunc(build)` |

The same three exist as `Rows` on `table`, `tbody`, `thead` and `tfoot`, `Cells` and `Headers` on `tr`, `Options` on `dropdown`, `datalist` and `optgroup`, and `Cols` on `colgroup`. Run `flint -info <element> ctors` to read them.

`ItemsOf` and `ItemsFunc` run their function at render time. The mapper passed to `ul.ItemsOf` returns `*li.Element`, so the compiler rejects any other child type. A mapper can return nil to skip an item.

```go
ul.ItemsOf(products, func(p Product) *li.Element {
    if p.Hidden {
        return nil
    }
    return Row(p)
})
```

### Return Types: Interface vs Concrete

**`node.Node` (interface)** - use when the function may return different element types, or the component is a final building block.

**`*element.Element` (concrete type, e.g. `*div.Element`)** - use when callers should be able to chain additional methods.

```go
// Return node.Node - flexible, no chaining after call
func Card(title, content string) node.Node {
    return div.New(
        h2.Text(title),
        p.Text(content),
    ).Class("card")
}

// Return *div.Element - allows continued chaining
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

### Layout with Dynamic Content

```go
func Layout(pageTitle string, content node.Node) node.Node {
    return html.New(
        head.New(
            title.Text(pageTitle),
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

## Performance

- Buffer pooling is enabled by default and handled automatically
- Each element has a chainable `BufferHint(n)` method for optional buffer size hints, and `RenderedSize()` to read the size recorded after `Render(w)`
- For high-throughput applications, [Fluent JIT](https://github.com/jpl-au/fluent-jit) provides additional optimisation (Compile, Tune, Flatten)

## Extending Fluent

Implement `node.Node` for a component that composes elements internally, or `node.Element` for a custom HTML element that needs attributes and open/close tags. Both are rare: a component is normally a plain function, as under [Components](#components).

### Implementing node.Node Directly

Most components are plain functions returning `node.Node`, as under [Components](#components) above. Implement the interface by hand only when callers need to configure the component by chaining, the way they chain onto an element. The type composes elements internally and satisfies `node.Node` (not `node.Element`) because it doesn't have its own tag or attributes.

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

// Render, WriteTo and RenderBytes are the three node.Node render entry
// points; each is a thin wrapper over RenderBuilder (mirroring text.Node).
func (f *EmailField) Render(w io.Writer) {
    _, _ = f.WriteTo(w)
}

func (f *EmailField) WriteTo(w io.Writer) (int64, error) {
    var buf bytes.Buffer
    f.RenderBuilder(&buf)
    n, err := w.Write(buf.Bytes())
    return int64(n), err
}

func (f *EmailField) RenderBytes() []byte {
    var buf bytes.Buffer
    f.RenderBuilder(&buf)
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

## Ecosystem

Fluent has companion packages that extend its capabilities. All are optional - Fluent works standalone for static HTML generation.

| Package | Description |
|---------|-------------|
| [Fluent JIT](https://github.com/jpl-au/fluent-jit) | Performance optimisation. **Compile** pre-renders static portions and re-evaluates dynamic content. **Tune** provides adaptive buffer sizing. **Flatten** pre-renders fully static content to raw bytes. Also provides two diff engines for live updates: the **Differ** (targeted patches by dynamic key) and the **Memoiser** (skips subtrees by version). |
| [Fluent HTMX](https://github.com/jpl-au/fluent-htmx) | HTMX integration. Accepts `node.Element` to set HTMX attributes (`hx-get`, `hx-post`, `hx-swap`, etc.) on any Fluent element. |

**How the packages relate:**
- **fluent** (this package) - core HTML generation, `node.Node`/`node.Element` interfaces, `.Dynamic()`/`.Memoise()` hook methods
- **fluent-jit** - rendering optimisation; its diff engines (Differ, Memoiser) produce `[]Patch` from two tree states
- **fluent-htmx** - attribute wrapper for HTMX, independent of the JIT diff engine

## Dot Import (Convenience Alternative)

The `dot` package drops the package prefixes, so `div.New()` becomes `Div()` and `p.Text("x")` becomes `P().Text("x")`.

Match the file you are editing. Use dot imports only where the surrounding file already does, and use the package-based API everywhere else. It is the primary API, and mixing the two styles across a codebase is worse than either one on its own. Specialised constructors such as `meta.UTF8()` need their own package import in both styles.

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

## Profile-Guided Optimisation (PGO)

Applications using Fluent benefit from [PGO](https://go.dev/doc/pgo) (Go 1.21+). It speeds up the rendering pipeline with no code changes. Allocations are unaffected - PGO improves inlining decisions only. Read the Go documentation for how to collect and apply a profile.
