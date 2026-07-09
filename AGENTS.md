# Fluent - HTML Generation for Go

Fluent is a type-safe, composable HTML generation library for Go. Every HTML element is a Go package (e.g. `div`, `a`, `input`). Elements are constructed with `New()`, configured with chainable methods, and rendered with `Render()`.

## Build & Test

```bash
go build ./...
go test ./...
go vet ./...
```

## Lint With Flint

After writing or editing Fluent code, run Flint. It catches Fluent-specific misuse that `go vet` cannot, and every diagnostic includes a `fix:` field with the corrected code - read the fix, apply it, and re-run until clean.

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

Use `-info` to look up an element's full API before writing against it - constructors, methods, typed parameters, and the constants each method accepts:

```bash
flint -info div          # everything about <div>
flint -info input        # every typed constant for every method
flint -info ol           # list constructors and typed variants
```

Exit codes: `0` clean, `1` diagnostics found, `2` usage or I/O error. Treat `1` as "fix and re-run", not "done".

## Methods That Do NOT Exist

**CRITICAL:** These methods do not exist anywhere in Fluent. Do not use them. Do not hallucinate them. They have never existed.

| Non-existent method | What to use instead |
|---------------------|---------------------|
| `.Attr()` | Use the dedicated typed method (e.g. `.Class()`, `.Href()`, `.Src()`) |
| `.SetAttr()` | Use `.SetAttribute()` for custom attributes only |
| `.Attribute()` | Use the dedicated typed method or `.SetAttribute()` |
| `.Attrs()` | No bulk attribute setter exists - set each attribute individually |
| `.WithAttr()` | Use the dedicated typed method or `.SetAttribute()` |

**The correct approach to setting attributes has three levels:**

1. **Dedicated typed methods (use first)** - Every standard HTML attribute has a chainable method. For example: `.Class()`, `.ID()`, `.Href()`, `.Src()`, `.Alt()`, `.Title()`, `.Disabled()`, `.Required()`, `.Placeholder()`, `.Name()`, `.Value()`, `.Type()`, etc.
2. **SetAria(key, value)** - For ARIA attributes. Automatically adds the `aria-` prefix.
3. **SetData(key, value)** - For data attributes. Automatically adds the `data-` prefix.
4. **SetAttribute(key, value)** - Only for truly custom or non-standard attributes (e.g. Alpine.js directives, HTMX attributes).

```go
// WRONG - these methods do not exist
div.New().Attr("class", "container")           // NO
div.New().SetAttr("id", "main")                // NO
button.New().Attribute("disabled", "")         // NO

// RIGHT - use dedicated typed methods
div.New().Class("container")                   // YES
div.New().ID("main")                           // YES
button.New().Disabled()                        // YES

// RIGHT - use SetAria for ARIA attributes
button.New().SetAria("label", "Close dialog")  // YES - renders aria-label="Close dialog"

// RIGHT - use SetData for data attributes
div.New().SetData("id", "123")                 // YES - renders data-id="123"

// RIGHT - use SetAttribute only for custom/non-standard attributes
div.New().SetAttribute("x-on:click", "handler")  // YES - Alpine.js directive
div.New().SetAttribute("hx-get", "/items")        // YES - HTMX attribute
```

**Important:** `SetAttribute()` does not return the element for chaining. `SetAria()` and `SetData()` do return the element for chaining.

## Common Mistakes - Read This First

**CRITICAL:** These are the most frequent errors. Every one of them causes a compile failure.

### 1. Never use raw strings where typed constants exist

Fluent uses `[]byte`-based types for enumerated HTML attribute values. The method signature **will not accept a plain string**. You must import the constant package and use its exported variable.

```go
// WRONG - does not compile
img.New().Loading("eager")              // cannot use "eager" (untyped string constant) as loading.Loading
img.New().Loading("lazy")               // same error
input.New().InputType("email")          // cannot use "email" as inputtype.InputType

// RIGHT - use the typed constant
img.New().Loading(loading.Eager)        // import "github.com/jpl-au/fluent/html5/attr/loading"
img.New().Loading(loading.Lazy)
input.New().InputType(inputtype.Email)   // import "github.com/jpl-au/fluent/html5/attr/inputtype"
```

Every attribute package also has a `Custom()` escape hatch for values not yet covered:
```go
loading.Custom("future-value")
```

### 2. Use element-specific constructors - do not build common patterns by hand

Most elements have domain-specific constructors that set multiple attributes at once. **Do not manually chain attributes when a constructor already exists.** For example:

```go
// WRONG - manual, verbose, error-prone
meta.New().SetAttribute("charset", "UTF-8")                    // SetAttribute returns void, not chainable
meta.New().SetAttribute("name", "viewport")                    // also wrong approach entirely
img.New().Src("photo.jpg").Alt("Sunset").Loading(loading.Lazy) // works but unnecessary

// RIGHT - use the constructor
meta.UTF8()                          // <meta charset="UTF-8" />
meta.Viewport("width=device-width") // <meta name="viewport" content="width=device-width" />
img.Lazy("photo.jpg", "Sunset")     // <img src="photo.jpg" alt="Sunset" loading="lazy" />
```

See the **Element-Specific Constructors** table below for every available constructor.

### 3. Do not hallucinate node-level helpers

These do **not** exist: `node.StaticText`, `node.RawNode`, `node.TextNode`, `node.HTML`. The actual helpers are:

| Exists | Usage |
|--------|-------|
| `node.Func(func() node.Node)` | Single dynamic node |
| `node.Funcs(func() []node.Node)` | Multiple dynamic nodes |
| `node.Condition(bool).True(n).False(n)` | Conditional rendering |
| `node.When(bool, node)` | Render when true |
| `node.Unless(bool, node)` | Render when false |
| `node.Empty()` | A node that renders nothing - the safe, explicit "render nothing" (prefer it over returning a nil, especially a typed nil pointer) |

Text content is always a method on the element: `div.Text(...)`, `div.Static(...)`, `div.RawText(...)`.

---

## Element-Specific Constructors

Beyond the universal constructors (`New`, `Text`, `Static`, `RawText`, `Textf`, `RawTextf`) that every element has, many elements provide domain-specific constructors. **Always prefer these over manual attribute chaining.**

**Typed constructors** enforce correct HTML nesting at compile time by accepting only the valid child element type (e.g. `ul.Items(items ...*li.Element)`). Use these instead of `New()` when the children are all the same element type. `New(...node.Node)` remains the untyped escape hatch for dynamic or mixed content.

**Cross-package constructors** create common child elements for you (e.g. `details.Summary("label", nodes...)` creates a `<summary>` child automatically).

| Package | Constructor | Description |
|---------|-------------|-------------|
| **a** | `Link(href, text)` | Anchor with href and link text |
| | `MailTo(email, text)` | `mailto:` link |
| | `JumpTo(anchor, text)` | Fragment/anchor link (`#anchor`) |
| | `Tel(number, text)` | `tel:` link |
| | `SMS(number, text)` | `sms:` link |
| | `FTP(url, text)` | `ftp:` link |
| | `DataURL(text, mime, data)` | Data URI link |
| | `Base64Data(text, mime, data)` | Base64 data URI link |
| **abbr** | `Titled(abbreviation, title)` | Abbreviation with title attribute |
| **area** | `Rect(x1, y1, x2, y2, href)` | Rectangular image map area |
| | `Circle(x, y, radius, href)` | Circular image map area |
| | `Poly(coords, href)` | Polygonal image map area |
| | `Default(href)` | Default image map area |
| **audio** | `Fallback(text)` | Audio with fallback text |
| | `Sources(nodes...)` | Audio with source elements |
| | `PreloadAuto(nodes...)` | Audio with preload=auto |
| | `PreloadMetadata(nodes...)` | Audio with preload=metadata |
| | `PreloadNone(nodes...)` | Audio with preload=none |
| **base** | `URL(href)` | Base URL for the document |
| **blockquote** | `NewCite(cite, nodes...)` | Blockquote with citation URL |
| | `TextCite(cite, text)` | Blockquote with citation and text |
| | `RawTextCite(cite, text)` | Blockquote with citation and raw HTML |
| **button** | `Submit(text)` | Submit button |
| | `Reset(text)` | Reset button |
| | `Button(text)` | Generic button (type=button) |
| **colgroup** | `Cols(cols...*col.Element)` | Type-safe column group from col elements |
| **data** | `Data(value, text)` | Data element with value |
| **datalist** | `Options(options...*option.Element)` | Type-safe datalist from option elements |
| **details** | `Summary(label, nodes...)` | Details with summary label and content |
| **dl** | `Pair(term, desc)` | Description list with dt/dd pair |
| **dropdown** | `Options(options...*option.Element)` | Type-safe select from option elements |
| **embed** | `PDF(src, w, h)` | Embedded PDF |
| | `Flash(src, w, h)` | Embedded Flash |
| | `Video(src, w, h)` | Embedded video |
| | `Audio(src, w, h)` | Embedded audio |
| **fieldset** | `Legend(caption, nodes...)` | Fieldset with legend and form controls |
| **figure** | `Caption(caption, nodes...)` | Figure with figcaption and content |
| **form** | `Get(action, nodes...)` | GET form |
| | `Post(action, nodes...)` | POST form |
| **html** | `Fragment(nodes...)` | HTML fragment without DOCTYPE |
| | `FragmentText(text)` | Fragment with text |
| | `FragmentStatic(text)` | Fragment with static text |
| | `FragmentRawText(text)` | Fragment with raw HTML |
| **iframe** | `Lazy(src)` | Iframe with loading=lazy |
| | `Eager(src)` | Iframe with loading=eager |
| **img** | `Src(src)` | Image with src only |
| | `Image(src, alt)` | Image with src and alt |
| | `Lazy(src, alt)` | Image with loading=lazy |
| | `Eager(src, alt)` | Image with loading=eager |
| **input** | `Password(name)` | Password input |
| | `Email(name)` | Email input |
| | `Search(name)` | Search input |
| | `Tel(name)` | Telephone input |
| | `URL(name)` | URL input |
| | `Number(name)` | Number input |
| | `Range(name)` | Range slider |
| | `Date(name)` | Date picker |
| | `Time(name)` | Time picker |
| | `DateTimeLocal(name)` | Datetime-local input |
| | `Month(name)` | Month picker |
| | `Week(name)` | Week picker |
| | `Checkbox(name, value)` | Checkbox |
| | `Radio(name, value)` | Radio button |
| | `File(name)` | File upload |
| | `Submit(value)` | Submit input |
| | `Button(value)` | Button input |
| | `Reset(value)` | Reset input |
| | `Hidden(name, value)` | Hidden input |
| | `Color(name)` | Colour picker |
| | `Image(name, src)` | Image input |
| **label** | `NewLabel(forID, nodes...)` | Label with for attribute and children |
| | `For(forID, text)` | Label with for attribute and text |
| **link** | `Stylesheet(href)` | CSS stylesheet link |
| | `Icon(href)` | Favicon link |
| | `Preload(href, as)` | Preload resource hint |
| **menu** | `Items(items...*li.Element)` | Type-safe menu from li elements |
| **meta** | `UTF8()` | `<meta charset="UTF-8" />` |
| | `Charset(charset)` | Meta with custom charset |
| | `Viewport(content)` | Viewport meta |
| | `Description(content)` | Meta description |
| | `Keywords(content)` | Meta keywords |
| | `Author(content)` | Meta author |
| | `Robots(content)` | Meta robots |
| | `OG(property, content)` | Open Graph meta |
| | `Refresh(seconds, url)` | Auto-refresh/redirect |
| **meter** | `ValueMax(value, max, nodes...)` | Meter with value and max |
| **object** | `PDF(data, nodes...)` | PDF object |
| | `Flash(data, nodes...)` | Flash object |
| | `Video(data, nodes...)` | Video object |
| | `Audio(data, nodes...)` | Audio object |
| **ol** | `Items(items...*li.Element)` | Type-safe ordered list from li elements |
| | `Decimal(items...*li.Element)` | Decimal numbered list (1, 2, 3) |
| | `LowerAlpha(items...*li.Element)` | Lowercase letter list (a, b, c) |
| | `UpperAlpha(items...*li.Element)` | Uppercase letter list (A, B, C) |
| | `LowerRoman(items...*li.Element)` | Lowercase Roman numeral list (i, ii, iii) |
| | `UpperRoman(items...*li.Element)` | Uppercase Roman numeral list (I, II, III) |
| **optgroup** | `Options(options...*option.Element)` | Type-safe option group |
| | `Labelled(label, options...*option.Element)` | Labelled option group |
| **option** | `Option(value, text)` | Select option with value |
| **progress** | `ValueMax(value, max, nodes...)` | Progress bar with value and max |
| **script** | `Module(src)` | ES module script |
| | `JavaScript(src)` | JavaScript src |
| | `JSON(data)` | Inline JSON (type=application/json) |
| **source** | `VideoMP4(src)` | MP4 video source |
| | `VideoWebM(src)` | WebM video source |
| | `VideoOgg(src)` | Ogg video source |
| | `AudioMP3(src)` | MP3 audio source |
| | `AudioOgg(src)` | Ogg audio source |
| | `AudioWav(src)` | WAV audio source |
| | `ImageWebP(srcset)` | WebP image source |
| | `ImageAVIF(srcset)` | AVIF image source |
| **style** | `CSS(css)` | Style element with CSS content |
| **table** | `Rows(rows...*tr.Element)` | Type-safe simple table from tr elements |
| **tbody** | `Rows(rows...*tr.Element)` | Type-safe tbody from tr elements |
| **tfoot** | `Rows(rows...*tr.Element)` | Type-safe tfoot from tr elements |
| **th** | `Col(content)` | Column header (scope=col) |
| | `Row(content)` | Row header (scope=row) |
| | `ColGroup(content)` | Column group header |
| | `RowGroup(content)` | Row group header |
| **thead** | `Rows(rows...*tr.Element)` | Type-safe thead from tr elements |
| **time** | `DateTime(datetime, content)` | Time with datetime attribute |
| **tr** | `Cells(cells...*td.Element)` | Type-safe data row from td elements |
| | `Headers(cells...*th.Element)` | Type-safe header row from th elements |
| **track** | `Subtitles(src)` | Subtitle track |
| | `Captions(src)` | Caption track |
| | `Descriptions(src)` | Description track |
| | `Chapters(src)` | Chapter track |
| | `Metadata(src)` | Metadata track |
| **ul** | `Items(items...*li.Element)` | Type-safe unordered list from li elements |
| **video** | `Src(src, nodes...)` | Video with src |
| | `PreloadAuto(nodes...)` | Video with preload=auto |
| | `PreloadMetadata(nodes...)` | Video with preload=metadata |
| | `PreloadNone(nodes...)` | Video with preload=none |

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

`Nodes()` returns the children that will actually be rendered. For conditionals this is only the active branch; for function components this is the evaluated function output. Tree walkers (including the Differ's snapshot collector) depend on `Nodes()` matching the rendered output.

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

// WRONG - these packages do not exist
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

**Rule:** Use `Static()` for unchanging content (labels, headings, boilerplate). Use `Text()` or `Textf()` for user input or values that change between renders. Use `RawText()` or `RawTextf()` only when you need to inject HTML and trust the source.

### Sanitising untrusted HTML

Fluent's `Text()` and `Textf()` HTML-escape via `html.EscapeString()` and are sufficient for plain text. For untrusted content that needs to render *as* HTML (rendered markdown, rich-text input, comment bodies), use the companion package [fluent-security](https://github.com/jpl-au/fluent-security):

```go
import "github.com/jpl-au/fluent-security"

// HTML - permissive UGC baseline. Keeps common formatting tags, strips
// scripts, event handlers, and dangerous URIs.
div.New(security.HTML(userMarkdown)).Class("body")

// PlainText - strips all tags, keeps only escaped text content.
h1.New(security.PlainText(article.RawTitle))
```

For reusable policies, hoist a `*security.Cleaner` at package scope (`security.RichText()` for the UGC baseline, `security.New()` for an empty allowlist) and extend with `Allow`, `AllowClasses`, `AllowAttr`. fluent-security also provides `Nonce()` for Content-Security-Policy workflows (inline `<script>`/`<style>` blocks via `script.Nonce(...)` / `style.Nonce(...)` matched against a CSP header). Fluent core never sanitises; that is deliberately a separate, opt-in package.

## Element Construction

### Choosing the Right Constructor

Every element has six constructors: `New`, `Text`, `Static`, `RawText`, `Textf`, `RawTextf`. **Pick the constructor that matches the element's primary content.** Only use `New()` when the element contains child nodes or has no content.

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

**When NOT to use `New()`:**
- The element's content is **just text** - use `Text()`, `Textf()`, `Static()`
- The element's content is **raw HTML** - use `RawText()`, `RawTextf()`

Attributes can be chained on any constructor. The constructor choice is about content, not attributes:
```go
h1.Text("Welcome").Class("title").ID("hero")  // text constructor + chained attributes
div.New(child1, child2).Class("grid")          // node constructor + chained attributes
```

### Content Methods (for adding to an existing element)

All non-self-closing elements also have these as **chainable methods** for adding content after construction. Use these only when you need to add content to an element that was constructed with `New()` for its child nodes:

- `.Text(s)` - adds escaped text content
- `.Textf(format, args...)` - adds formatted escaped text
- `.Static(s)` - adds static text (JIT-optimisable)
- `.RawText(s)` - adds unescaped HTML content
- `.RawTextf(format, args...)` - adds formatted unescaped HTML

```go
// Method form - adding text alongside child nodes
div.New(img.New().Src("photo.jpg")).Text("Caption")

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
- `.SetEvent(key, value)` - for custom event attributes

**Element-specific attributes** - each element has its own methods. Examples:

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

### Boolean Attributes

Boolean HTML attribute methods accept an optional `bool`. Called with no arguments they set the attribute. Called with a single `bool` they set it only when the condition is true - no subtree duplication, no intermediate variable, no break in the chain.

```go
// No argument - attribute always present
input.Email("email").Required()

// Inline condition - attribute present only when true
input.New().Checked(slices.Contains(selected, c.ID))
button.Text("Save").Disabled(!form.Valid)
script.New().Src("/app.js").Async(isProduction)
textarea.New().ReadOnly(user.IsGuest)
```

This pattern is consistent across every boolean attribute in Fluent:

| Element | Methods |
|---------|---------|
| `input` | `.Checked()`, `.Disabled()`, `.Required()`, `.ReadOnly()`, `.Multiple()`, `.AutoFocus()` |
| `button` | `.Disabled()`, `.AutoFocus()` |
| `form` | `.NoValidate()`, `.AutoFocus()` |
| `textarea` | `.Disabled()`, `.ReadOnly()`, `.Required()`, `.AutoFocus()` |
| `dropdown` (`<select>`) | `.Disabled()`, `.Multiple()`, `.Required()`, `.AutoFocus()` |
| `option` | `.Disabled()`, `.Selected()`, `.AutoFocus()` |
| `script` | `.Async()`, `.Defer()`, `.AutoFocus()` |
| `details`, `dialog` | `.Open()`, `.AutoFocus()` |
| `audio`, `video` | `.Controls()`, `.Loop()`, `.Muted()`, (`video` adds `.PlaysInline()`), `.AutoFocus()` |

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

### Attribute Escaping and URL Filtering

Attribute values are secured automatically - you do not escape them yourself:

- **Every attribute value is HTML-escaped at set time** - typed methods (`.Class()`, `.Title()`, ...), `SetAttribute()`, `SetData()`, `SetAria()`, `.Dynamic()` keys, and enum `Custom()` values all pass through `node.EscapeAttribute`. A value can never break out of its quotes.
- **URL navigation sinks are scheme-filtered.** `href`, `src`, `action`, `formaction`, and `<object>` `data` are checked against an allowlist (`http`, `https`, `mailto`, `tel`, `sms`, and relative URLs). An off-allowlist scheme (`javascript:`, `vbscript:`, `data:` in a nav sink) is replaced with the inert `node.UnsafeURL` sentinel. Register `node.OnUnsafeURL(fn)` to observe rejections. Image/media sinks (`img`/`video` `src`, `srcset`, `poster`) are escaped but not filtered, so `data:image/...` still renders.
- **Do not pre-escape.** Calling `html.EscapeString` before a setter double-escapes; pass the raw value and let the setter escape it once.
- **`SetAttributeRaw(key, value)`** is the per-value trusted-value hatch (the attribute mirror of `RawText`): it stores the value verbatim, without escaping. Use only for values you fully trust.

### Type-Safe Constants

Fluent uses typed constants for attributes with enumerated values. Methods accept a typed constant, not a string - so typos cause compile errors.

```go
// Typed constant - compile error on typo
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

The `.Dynamic(key)` method on HTML elements marks them for reactive tracking by the [Fluent JIT](https://github.com/jpl-au/fluent-jit) diff engine. The key identifies the element across renders so the diff engine can detect changes and produce targeted patches.

```go
// Mark an element for reactive tracking
span.Textf("Count: %d", state.Count).Dynamic("count")

// The key renders as a data-fluent-key attribute
// <span data-fluent-key="count">Count: 42</span>

// Keys must be unique within a render tree, and are required
p.Text(state.ErrorMsg).Dynamic("error-message")
table.New(rows...).Dynamic("data-table")
```

`.Dynamic()` is chainable and follows the same pattern as `.Class()`, `.SetData()`, etc. It is used by [Fluent JIT](https://github.com/jpl-au/fluent-jit) to identify which segments need re-evaluation.

Elements without `.Dynamic()` are not tracked - the diff engine only examines keyed nodes.

## Component Pattern

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
    items := make([]*li.Element, len(products))
    for i, p := range products {
        items[i] = li.New(
            h3.Text(p.Name),
            span.Textf("$%.2f", p.Price),
        )
    }
    return ul.Items(items...)
}
```

## Performance

- Use `Static()` for unchanging content (enables JIT optimisation)
- Buffer pooling is enabled by default and handled automatically
- Each element has a chainable `BufferHint(n)` method for optional buffer size hints, and `RenderedSize()` to read the size recorded after `Render(w)`
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

## Typed Attributes Reference

**CRITICAL:** Methods that accept typed constants will **NOT** accept raw strings - this is a compile error, not a runtime error. You must import the constant package and use its exported variable. Every package also provides `Custom(string)` for edge cases.

Import path pattern: `github.com/jpl-au/fluent/html5/attr/<package>`

```go
// WRONG - does not compile
img.New().Loading("lazy")                   // cannot use "lazy" (string) as loading.Loading
input.New().InputType("email")              // cannot use "email" (string) as inputtype.InputType

// RIGHT - use the typed constant
img.New().Loading(loading.Lazy)             // import "github.com/jpl-au/fluent/html5/attr/loading"
input.New().InputType(inputtype.Email)      // import "github.com/jpl-au/fluent/html5/attr/inputtype"
link.New().Rel(rel.Stylesheet)              // import "github.com/jpl-au/fluent/html5/attr/rel"
```

### Complete Constant Reference

| Package | Type | Values |
|---------|------|--------|
| `accept` | `Accept` | `ImageWildcard`, `VideoWildcard`, `AudioWildcard`, `ImageJPEG`, `ImagePNG`, `ImageGIF`, `ImageWebP`, `ImageSVG`, `MimePDF`, `MimeMSWord`, `MimeWordDOCX`, `TextPlain`, `TextCSV`, `JPG`, `JPEG`, `PNG`, `GIF`, `WebP`, `SVG`, `PDF`, `DOC`, `DOCX`, `TXT`, `CSV`, `XML`, `VideoMP4`, `VideoWebM`, `AudioMP3`, `AudioWAV` |
| `as` | `As` | `Audio`, `Document`, `Embed`, `Fetch`, `Font`, `Image`, `Object`, `Script`, `Style`, `Track`, `Video`, `Worker` |
| `autocapitalize` | `AutoCapitalize` | `Off`, `None`, `On`, `Sentences`, `Words`, `Characters` |
| `autocomplete` | `AutoComplete` | `On`, `Off`, `Name`, `Email`, `CurrentPassword`, `NewPassword`, `Username`, `AddressLine1`, `AddressLine2`, `Country`, `PostalCode`, `Tel`, `Url` |
| `autocorrect` | `AutoCorrect` | `On`, `Off` |
| `blocking` | `Blocking` | `Render` |
| `capture` | `Capture` | `User`, `Environment` |
| `charset` | `Charset` | `UTF8`, `ISO88591`, `Windows1252` |
| `contenteditable` | `ContentEditable` | `True`, `False`, `PlaintextOnly` |
| `controlslist` | `ControlsList` | `NoDownload`, `NoFullscreen`, `NoRemotePlayback` |
| `crossorigin` | `CrossOrigin` | `Anonymous`, `UseCredentials` |
| `decoding` | `Decoding` | `Sync`, `Async`, `Auto` |
| `dir` | `Dir` | `LeftToRight`, `RightToLeft`, `Auto` |
| `enctype` | `EncType` | `UrlEncoded`, `MultipartFormData`, `TextPlain` |
| `enterkeyhint` | `EnterKeyHint` | `Enter`, `Done`, `Go`, `Next`, `Previous`, `Search`, `Send` |
| `fetchpriority` | `FetchPriority` | `High`, `Low`, `Auto` |
| `formmethod` | `FormMethod` | `Get`, `Post` |
| `inputmode` | `InputMode` | `None`, `Text`, `Tel`, `Url`, `Email`, `Numeric`, `Decimal`, `Search` |
| `inputtype` | `InputType` | `Text`, `Password`, `Email`, `Search`, `Tel`, `Url`, `Number`, `Range`, `Date`, `Time`, `DatetimeLocal`, `Month`, `Week`, `Checkbox`, `Radio`, `File`, `Submit`, `Button`, `Reset`, `Hidden`, `Color`, `Image` |
| `listtype` | `ListType` | `LowerAlpha`, `UpperAlpha`, `LowerRoman`, `UpperRoman`, `Decimal` |
| `loading` | `Loading` | `Eager`, `Lazy` |
| `media` | `Media` | `Screen`, `Print`, `All`, `Speech`, `Mobile`, `Tablet`, `Desktop`, `SmallMobile`, `LargeMobile`, `SmallTablet`, `LargeTablet`, `LargeDesktop`, `Portrait`, `Landscape`, `Retina`, `HighDPI` |
| `method` | `Method` | `Get`, `Post`, `Dialog` |
| `popover` | `Popover` | `Auto`, `Manual` |
| `popovertargetaction` | `PopoverTargetAction` | `Toggle`, `Show`, `Hide` |
| `preload` | `Preload` | `None`, `Metadata`, `Auto` |
| `referrerpolicy` | `ReferrerPolicy` | `NoReferrer`, `NoReferrerWhenDowngrade`, `Origin`, `OriginWhenCrossOrigin`, `SameOrigin`, `StrictOrigin`, `StrictOriginWhenCrossOrigin`, `UnsafeUrl` |
| `rel` | `Rel` | `Stylesheet`, `Icon`, `Preload`, `Prefetch`, `DnsPrefetch`, `Preconnect`, `Canonical`, `Alternate`, `Prev`, `Next`, `Help`, `License`, `Manifest`, `ModulePreload`, `AppleTouchIcon` |
| `sandbox` | `Sandbox` | `AllowDownloads`, `AllowForms`, `AllowModals`, `AllowOrientationLock`, `AllowPointerLock`, `AllowPopups`, `AllowPopupsToEscapeSandbox`, `AllowPresentation`, `AllowSameOrigin`, `AllowScripts`, `AllowStorageAccessByUserActivation`, `AllowTopNavigation`, `AllowTopNavigationByUserActivation` |
| `shape` | `Shape` | `Rect`, `Circle`, `Poly`, `Default` |
| `sizes` | `Size` | `FullWidth`, `HalfWidth`, `ThirdWidth`, `QuarterWidth`, `Small`, `Medium`, `Large`, `ExtraLarge`, `MobileFullTabletHalf`, `MobileFullDesktopThird`, `ResponsiveHero`, `ResponsiveContent` |
| `spellcheck` | `Spellcheck` | `True`, `False` |
| `target` | `Target` | `Self`, `Blank`, `Parent`, `Top` |
| `translate` | `Translate` | `Yes`, `No` |
| `virtualkeyboardpolicy` | `VirtualKeyboardPolicy` | `Auto`, `Manual` |
| `writingsuggestions` | `WritingSuggestions` | `True`, `False` |

## Ecosystem

Fluent has companion packages that extend its capabilities. All are optional - Fluent works standalone for static HTML generation.

| Package | Description |
|---------|-------------|
| [Fluent JIT](https://github.com/jpl-au/fluent-jit) | Performance optimisation. **Compile** pre-renders static portions and re-evaluates dynamic content. **Tune** provides adaptive buffer sizing. **Flatten** pre-renders fully static content to raw bytes. Also provides the **Diff** engine that compares renders by dynamic key and produces `[]Patch`. |
| [Fluent HTMX](https://github.com/jpl-au/fluent-htmx) | HTMX integration. Accepts `node.Element` to set HTMX attributes (`hx-get`, `hx-post`, `hx-swap`, etc.) on any Fluent element. |

**How the packages relate:**
- **fluent** (this package) - core HTML generation, `node.Node`/`node.Element` interfaces, `.Dynamic()` method
- **fluent-jit** - rendering optimisation, diff engine produces `[]Patch` from two tree states
- **fluent-htmx** - attribute wrapper for HTMX, independent of the JIT diff engine

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

## Profile-Guided Optimisation (PGO)

Applications using Fluent benefit from [PGO](https://go.dev/doc/pgo) (Go 1.21+). Collect a CPU profile from production, place it as `default.pgo` in the main package, and `go build` applies it automatically. Expect 10-20% speed improvements across the rendering pipeline with no code changes. Allocations are unaffected - PGO improves inlining decisions only. Collect fresh profiles periodically as code evolves.
