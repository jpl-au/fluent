# Fluent

HTML5 components in Go using a Fluent API.

**Quick links:** [Why Fluent?](#why-fluent) · [Install](#install) · [Quick Start](#quick-start) · [Flint (linter & element info)](#flint) · [Reserved Keywords](#reserved-keywords) · [Static vs Dynamic Content](#static-vs-dynamic-content) · [Typed Constructors](#typed-constructors) · [Conditional Rendering](#conditional-rendering) · [Functional Processing](#functional-processing) · [Building Components](#building-components) · [Type-Safe Attributes](#type-safe-attributes) · [Architecture](#architecture) · [Performance](#performance) · [Generator](#generator) · [Ecosystem](#ecosystem) · [PGO](#profile-guided-optimisation-pgo)

## Why Fluent?

**No template language to learn.** Write HTML using Go code. Get IDE auto-completion, type checking, and refactoring support for free.

**Built for developers.** Thoughtful around the developer experience: attributes use native Go types - set a `width` with an `int`, a `volume` with a `float64`. Fluent handles the conversion. Type-safe constants for enumerated values catch typos like `type="emial"`.

**Type-safe nesting.** Typed constructors enforce correct HTML parent-child relationships at compile time. `ul.Items()` only accepts `*li.Element`, `tr.Cells()` only accepts `*td.Element` - the compiler catches nesting mistakes that would otherwise become silent bugs. `New()` remains available as the flexible escape hatch. [See Typed Constructors](#typed-constructors).

**HTML escaping by default.** `Text()` and `Textf()` automatically escape `<`, `>`, `&`, and quotes, and every attribute value is escaped (and URL sinks scheme-filtered) at set time too - see [Attribute escaping](#attribute-escaping). For untrusted HTML that needs to render as HTML (rendered markdown, rich-text input), reach for the opt-in [fluent-security](https://github.com/jpl-au/fluent-security) package, which wraps [bluemonday](https://github.com/microcosm-cc/bluemonday) and returns Fluent nodes directly.

**Performance considered.** Buffer pooling and efficient rendering for high-throughput applications. Don't want to use `sync.Pool`? Just turn it off.

**Extensible.** Every element exposes its attributes through an interface. Use that interface to extend an element for a framework such as htmx, Datastar or Turbo. The same interface lets you replace an element and build a web component.

**Optional JIT optimisations.** Three strategies (Compile, Tune, Flatten) available via a separate package for high-throughput applications. [See Performance](#performance).

**HTML5 spec aligned.** Elements and attributes follow the HTML5 specification, generated from YAML definitions. [See Generator](#generator).

*Interested in why I created another HTML rendering library for Go? [See my motivations](#why-fluent-my-motivations).*

## Install

```bash
go get github.com/jpl-au/fluent
```

## Quick Start

```go
package main

import (
    "net/http"

    "github.com/jpl-au/fluent/html5/body"
    "github.com/jpl-au/fluent/html5/div"
    "github.com/jpl-au/fluent/html5/h1"
    "github.com/jpl-au/fluent/html5/head"
    "github.com/jpl-au/fluent/html5/html"
    "github.com/jpl-au/fluent/html5/p"
    "github.com/jpl-au/fluent/html5/title"
)

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        html.New(
            head.New(
                title.Text("home"),
            ),
            body.New(
                div.New(
                    h1.Text("Hello, World"),
                    p.Text("Built with Fluent."),
                ).Class("container"),
            ),
        ).Render(w)
    })
    http.ListenAndServe(":8080", mux)
}
```

## Flint

[Flint](https://github.com/jpl-au/flint) is the companion linter and element reference for Fluent. It checks Fluent calls against a registry generated from the same specs as the packages, and every diagnostic carries a `fix:` line with the correction.

```bash
go install github.com/jpl-au/flint/cmd/flint@latest
flint ./...              # lint
flint -info input        # an element's constructors, methods and attribute mappings
flint -info inputtype    # the values a typed method accepts
```

An element named after a Go reserved word resolves from its HTML name, so `flint -info select` lands on `dropdown`. The full list of checks is in the [Flint README](https://github.com/jpl-au/flint#what-it-checks).

## Reserved Keywords

Some HTML elements conflict with Go reserved keywords. I chose names that still feel intuitive - `dropdown` for `<select>` felt natural since that's what it renders.

| HTML Element | Fluent Package |
|--------------|----------------|
| `<select>`   | `dropdown`     |
| `<main>`     | `primary`      |
| `<var>`      | `variable`     |

## Documentation for AI agents

- `AGENTS.md` - a guide for AI agents working with Fluent. It also serves as a detailed reference for developers.

## Static vs Dynamic Content

`Static()`, `Text()`, `Textf()`, `RawText()` and `RawTextf()` are both element constructors and methods. This was largely a requirement for the developer experience when using dot imports (which are optional).

The `Static()` constructor/method exists for use with [Fluent JIT](#performance) as a signal that this element is static. You should not use a variable in `Static()` as the JIT compiler will only render this node once (the first run) and subsequent calls will ignore changes. An alternative was to mark the node as dynamic, but I thought the developer experience would be hindered. Nothing fails at runtime when you get this wrong. The content goes stale, so Flint's `static` check flags a `Static()` call given anything other than a string literal.

```go
// Static - content known at definition time (see Performance section)
div.Static("Copyright 2024")

// Dynamic - HTML-escaped, safe for user input
div.Text(user.Name)
div.Textf("Hello %s, you have %d messages", user.Name, count)

// Raw - unescaped, use only for trusted HTML
div.RawText("<em>Bold</em>")
div.RawTextf("<span class=\"%s\">%s</span>", className, content)
```

### Attribute escaping

Text content is not the only injection surface - attribute *values* are too. Fluent escapes an attribute value when you set it, so the value stays inside its quotes instead of becoming markup or an event handler. This is on by default. It covers constructors and setters alike: typed setters (`.Class()`, `.Title()`, `.Value()`, ...), `.SetAttribute()`, `.SetData()`, `.Dynamic()` keys, and enum `Custom(...)` values.

```go
// The value is escaped for you - the rendered attribute stays a single,
// inert attribute no matter what the user typed.
div.New().SetAttribute("data-note", userInput)
```

Because escaping is automatic, do not pre-escape values yourself. Calling `html.EscapeString` before a setter double-escapes the value: the browser decodes one layer on `getAttribute` and the caller sees stray `&amp;`/`&#34;` artefacts (and JSON stored in a data attribute stops parsing). Pass the raw value and let the setter escape it once.

For a value you have already sanitised and trust verbatim, `SetAttributeRaw` is the per-value hatch - the attribute mirror of `RawText`:

```go
// Stored verbatim, no escaping. Use only for values you fully trust.
el.SetAttributeRaw("data-html", trustedFragment)
```

**URL attributes are also scheme-checked.** This covers `href`, `src`, `action`, `formaction` and `<object>` `data`. These attributes accept `http`, `https`, `mailto`, `tel`, `sms` and relative URLs. Fluent replaces any other scheme, such as `javascript:` or `vbscript:`, with `node.UnsafeURL`, which the browser does not act on. Call `node.OnUnsafeURL` to see which values Fluent replaced. Image and media attributes are escaped but not scheme-checked, so a `data:image/...` URL still renders. This covers `src` on `img` and `video`, and `srcset` and `poster`.

The whole layer - escaping and URL filtering - is baked into the generated setters at generation time, so it costs nothing to branch on at runtime.

### Reactive tracking

Every element carries chainable hook methods - `.Dynamic("key")`, `.Memoise(version)`, `.MemoiseKey()` - for external render engines such as [Fluent JIT](https://github.com/jpl-au/fluent-jit). Fluent's part is mechanical: a `.Dynamic` key renders as the element's `id` attribute, nothing more (an explicit `.ID()` that differs from the key panics at render). Plain rendering ignores all three, so they cost nothing unless an engine consumes them; usage and semantics are documented by the engine.

```go
span.Textf("Count: %d", count).Dynamic("count")
```

The semantics - targeted patches, subtree skipping, cross-session caching - are documented by the consuming engine; see the [Fluent JIT documentation](https://github.com/jpl-au/fluent-jit).

### Convenience constructors

Many elements have convenience constructors for common use cases:

```go
// Form shortcuts
form.Get("/search", ...)   // <form action="/search" method="get">
form.Post("/login", ...)   // <form action="/login" method="post">

// Input types - all return *element for chaining
input.Email("email")             // <input type="email" name="email">
input.Password("password")       // <input type="password" name="password">
input.Checkbox("agree", "yes")   // <input type="checkbox" name="agree" value="yes">
input.Submit("Submit")           // <input type="submit" value="Submit">

// Chain additional attributes as needed
input.Email("email").
    Placeholder("you@example.com").
    Required().
    AutoComplete(autocomplete.Email) // typed constant, not a string
```

## Typed Constructors

Constructors that enforce correct HTML nesting at compile time. These accept only the valid child element type - `New()` remains the untyped escape hatch for dynamic or mixed content.

```go
// Lists - only accept *li.Element
ul.Items(li.Text("one"), li.Text("two"))
ol.Decimal(li.Text("first"), li.Text("second"))

// Tables - only accept the correct child type
table.Rows(
    tr.Headers(th.Col("Name"), th.Col("Age")),
    tr.Cells(td.Text("Alice"), td.Text("30")),
)
thead.Rows(tr.Headers(th.Col("Name")))
tbody.Rows(tr.Cells(td.Text("Alice")))

// Select/datalist - only accept *option.Element
dropdown.Options(option.Option("au", "Australia"), option.Option("nz", "New Zealand"))
optgroup.Labelled("Oceania", option.Option("au", "Australia"))

// Cross-package constructors - create child elements for you
details.Summary("Click to expand", p.Text("Hidden content"))
fieldset.Legend("Address", input.Text("street", ""))
figure.Caption("An elephant", img.Src("elephant.jpg"))
dl.Pair("Name", "Alice")

// New() is always available as the untyped escape hatch
table.New(caption.Text("Title"), thead.New(...), tbody.New(...))
```

Most elements carry their own construction helpers beyond the ones shown here, both typed constructors like these and the convenience constructors above. Run `flint -info <element> ctors` to see what any element offers before building it by hand.

### Lists that come from data

When the children come from a slice, each collection constructor (`Items`, `Rows`, `Cells`, `Headers`, `Options`, `Cols`) has two deferred siblings that keep the compile-time check. Styled variants such as `ol.Decimal` do not:

```go
// ItemsOf maps a slice. The mapper returns the child element type,
// so the compiler still rejects a wrong child. A nil result skips the item.
ul.ItemsOf(products, func(p Product) *li.Element {
    return li.Text(p.Name)
})

// ItemsFunc takes a function, for logic that reorders or computes the list.
ul.ItemsFunc(func() []*li.Element { ... })

// The same pair on every collection constructor
table.RowsOf(people, rowFor)      // fn func(T) *tr.Element
tr.CellsOf(values, cellFor)       // fn func(T) *td.Element
dropdown.OptionsOf(countries, at) // fn func(T) *option.Element
```

Both forms run their function at render time.

## Conditional Rendering

`node.Condition()` provides inline conditional rendering. `True()` and `False()` can be used together or independently:

```go
// Both branches
node.Condition(user.IsLoggedIn).
    True(p.Text("Welcome back!")).
    False(a.Link("/login", "Sign in"))
```

For single-branch conditions, `When()` and `Unless()` provide concise shorthand:

```go
// Render only when condition is true
node.When(user.IsAdmin, span.Static("Admin"))

// Render only when condition is false
node.Unless(user.IsLoggedIn, a.Link("/login", "Sign in"))
```

`When()` and `Unless()` return nodes, so they slot directly into the children of a parent element. This is how you add conditional content inside a larger tree without duplicating the surrounding subtree. Remember that `.Text(...)` on an element is sugar for appending a `text.Text(...)` child - when you need a text node to appear conditionally, build it as a child via the `text` package:

```go
import "github.com/jpl-au/fluent/text"

span.New(
    node.When(len(foo) > 0, text.Textf("foo: %d", foo)),
    span.New(text.Text("g")).Class("weight"),
).Class("amount")
```

The outer `span.amount` always renders. The `foo: %d` text appears only when the condition holds. No `Condition(...).True(...).False(...)` mirror, no duplicated children.

Conditions can be nested since `node.Condition()` returns a `node.Node`:

```go
node.Condition(user.IsLoggedIn).
    True(
        node.Condition(user.IsAdmin).
            True(span.Static("Admin Dashboard")).
            False(span.Static("User Dashboard")),
    ).
    False(a.Link("/login", "Sign in"))
```

For multiple branches, `node.Func()` is cleaner:

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

## Functional Processing

For logic more complex than a simple condition - database lookups, error handling, or building dynamic content - `node.Func()` lets you write arbitrary Go code inline. The function executes at render time, keeping your tree structure declarative while deferring complex logic until it's needed.

```go
node.Func(func() node.Node {
    count, err := db.GetUnreadCount(userID)
    if err != nil || count == 0 {
        return node.Empty()
    }
    return div.Textf("You have %d unread messages", count).Class("notification")
})
```

`node.Empty()` returns a node that renders nothing. Use it wherever a `node.Node` is expected and there is nothing to render: a `node.Func` body, a conditional branch, or a children slice. An untyped `nil` also renders nothing, because the render loops skip nil children, but a typed nil pointer does not. A function declared to return `*div.Element` that returns `nil` hands the render loop a non-nil interface holding a nil pointer, which panics. `node.Empty()` avoids that.

For the common case of one node per slice element, `node.Map` removes the loop boilerplate:

```go
node.Map(products, func(p Product) node.Node {
    return li.New(
        span.Text(p.Name),
        span.Textf("$%.2f", p.Price),
    )
})
```

`Map` is for readability, not speed. It wraps `node.Funcs`, so the loop runs at render time inside a closure, and that closure and the component it returns are both extra allocations. The overhead is fixed rather than per element, so it is negligible on a long list and proportionally significant on a short one. Build the slice by hand on a hot path with few children.

Boolean attribute methods (`.Checked`, `.Disabled`, `.Required`, `.ReadOnly`, `.Multiple`, `.Selected`, `.Async`, `.Defer`, `.Open`, and friends) accept an optional `bool`: call them with no arguments to set the attribute, or pass a condition to set it only when the condition is true. This keeps conditionals inline inside the mapped element, with no subtree duplication:

```go
node.Map(categories, func(c Category) node.Node {
    return label.New(
        input.New().Type(inputtype.Checkbox).
            Name("category").
            Value(strconv.Itoa(c.ID)).
            Checked(slices.Contains(selected, c.ID)),
        text.Text(c.Name),
    )
})
```

When the per-item logic is more involved than a one-liner - database lookups, error handling, building different node shapes per item - drop down to `node.Funcs` and write the loop by hand:

```go
node.Funcs(func() []node.Node {
    items := make([]node.Node, 0, len(products))
    for _, p := range products {
        if !p.Visible {
            continue
        }
        items = append(items, li.New(
            span.Text(p.Name),
            span.Textf("$%.2f", p.Price),
        ))
    }
    return items
})
```

## Building Components

There's no special component system to learn - building your own components is handled through Go functions. You get all the benefits of Go's type system, testing, and refactoring tools. Components are just functions that return `node.Node` or a concrete element type.

Where a component draws something your application already has a type for, take that type. The component then decides how the record is drawn, including what an empty list looks like, and the caller passes what it loaded:

```go
func UserGreeting(user User) node.Node {
    return div.New(
        img.Image(user.Avatar, user.Name),
        h3.Text(user.Name),
        // One branch: nothing renders when the condition is false.
        node.When(user.IsAdmin, span.Static("Admin")),
        // Two branches: something renders either way.
        node.Condition(user.Online).
            True(span.Static("Online")).
            False(span.Static("Offline")),
    ).Class("user-greeting")
}

func UserList(users []User) node.Node {
    if len(users) == 0 {
        return p.Static("Nobody here yet.").Class("empty")
    }
    return ul.ItemsOf(users, func(u User) *li.Element {
        return li.New(UserGreeting(u))
    }).Class("users")
}
```

`UserList` owns the empty case, so no caller writes `if len(users) == 0`. A condition about the shape of the data belongs in the component, or it is repeated at every call site.

Take scalar parameters where the inputs are unrelated values your application has no type for:

```go
// Return node.Node for flexibility - can return different element types
func Card(heading string, content string) node.Node {
    return div.New(
        h2.Text(heading),
        p.Text(content),
    ).Class("card")
}

// Return concrete type to allow continued chaining after the call
func Card(heading string, content string) *div.Element {
    return div.New(
        h2.Text(heading),
        p.Text(content),
    ).Class("card")
}

// With concrete return type, callers can chain additional methods:
Card("Welcome", "Hello!").ID("welcome-card").Class("highlighted")
```

Use them like any other element:

```go
page := div.New(
    Card("Welcome", "Thanks for signing up!"),
    UserList(users),
)
page.Render(w)
```

## Type-Safe Attributes

Fluent provides type-safe constants for attributes with enumerated values. I wanted the IDE to do the heavy lifting. When you type `inputtype.`, your editor shows you every valid option - no more checking MDN to remember if it's `datetime-local` or `datetimeLocal`.

Methods like `Type()` accept typed constants, not strings - so `input.New().Type("emial")` won't compile. Each attribute package also provides a `Custom()` function for edge cases or future HTML specifications not yet covered.

```go
import (
    "github.com/jpl-au/fluent/html5/input"
    "github.com/jpl-au/fluent/html5/attr/inputtype"
    "github.com/jpl-au/fluent/html5/attr/autocomplete"
)

input.New().
    Type(inputtype.Email).            // Typed constant, not a string
    AutoComplete(autocomplete.Email). // IDE shows all valid options
    Required()

// For edge cases or future specs
input.New().Type(inputtype.Custom("future-type"))
```

## Architecture

Fluent is organised into several packages:

| Package | Description |
|---------|-------------|
| `node` | Core `Node` interface for all renderable types: `Render(w)`, `WriteTo(w)`, `RenderBytes()`, `RenderBuilder()`, `Nodes()`. The `Element` interface extends `Node` with `SetAttribute()`, `SetAttributeRaw()`, `RenderOpen()`, `RenderClose()` for HTML elements |
| `html5/*` | HTML5 elements, one package per element (e.g., `div`, `span`, `input`). Each provides `New()`, `Text()`, `Static()` constructors |
| `html5/attr/*` | Type-safe attribute constants (e.g., `inputtype.Email`, `autocomplete.Off`, `rel.Stylesheet`) |
| `text` | Text node implementations for `Static()`, `Text()`, `RawText()` and their formatted variants |
| `pool` | Buffer pooling configuration |
| `dot` | Optional dot import for cleaner syntax without package prefixes |

### Everything is a Node

The `node.Node` interface is the foundation of Fluent. Every renderable piece of content implements it: HTML elements, text nodes, conditionals (`node.Condition`), function wrappers (`node.Func`), and the empty node (`node.Empty`). This unified interface enables arbitrary composition - any `node.Node` can be a child of any element.

HTML elements also implement `node.Element`, which extends `Node` with `SetAttribute()`, `SetAttributeRaw()`, `RenderOpen()`, and `RenderClose()`. Text nodes, function components, and conditionals are not elements - they don't have attributes or tags.

When in doubt about return types for your components, `node.Node` is always safe:

```go
func MyComponent(showHeader bool) node.Node {
    if showHeader {
        return header.New(h1.Text("Welcome"))
    }
    return node.Empty()  // renders nothing
}
```

Returning concrete types (like `*div.Element`) allows method chaining after the call, but `node.Node` provides maximum flexibility when your component might return different element types or nothing at all. A `nil` return also renders nothing. Prefer `node.Empty()`, which is a real node and cannot become a nil pointer inside the interface.

### Rendering

All nodes implement `node.Node`. The usual path is `Render(w)` - write straight to an `io.Writer` such as an `http.ResponseWriter`. Reach for `WriteTo(w)` when you need the write error, or `RenderBytes()` when you actually need a `[]byte`:

```go
// Write to a writer (the common case; write errors discarded)
page.Render(w)

// Same, but observe the byte count and any write error
n, err := page.WriteTo(w)

// Allocate and return the bytes (tests, caching, string building)
html := page.RenderBytes()
```

For building complex trees efficiently, `RenderBuilder(*bytes.Buffer)` writes directly to a shared buffer.

### Attributes

Fluent uses a tiered approach for attributes, based on MDN documentation:

- **Inlined fields** - Very common attributes (`class`, `id`, `style`) are direct struct fields for efficient access
- **Global attributes** - Attributes available on all elements (e.g., `hidden`, `tabindex`, `title`) are embedded via a shared struct
- **Event attributes** - Event handlers (`onclick`, `onchange`, ...) are embedded via a separate shared struct
- **Element-specific** - Attributes unique to an element (e.g., `href` on `<a>`, `src` on `<img>`) are direct struct fields
- **Generic slice** - Any additional or custom attributes are stored in an attributes slice

This design balances memory efficiency with access speed for the most commonly used attributes.

### Constants

Fluent uses `[]byte` variables for common patterns which removes the need for string to `[]byte` conversions when writing to the `bytes.Buffer`. It's a small optimisation, but it's there. These are stored in a constants.go file. The same principle can be used for extensions.

## Is it stable?

I've been using Fluent in production for building HTML since July 2025, which has helped me iron out many of the bugs and issues prior to releasing it to the public. In all fairness, I have not put it through its paces nor started using any of the JIT optimisation features (that whole premature optimisation == root of all evil concept).

## Advanced

### Buffer Pooling

Fluent uses a two-tier buffer pool (`sync.Pool`) to balance memory efficiency across different render sizes. When you call `Render(w)` with a writer, pooled buffers are used automatically:

```go
func handler(w http.ResponseWriter, r *http.Request) {
    page := html.New(
        head.New(title.Text("My Page")),
        body.New(div.Text("Hello")),
    )
    page.Render(w)  // Pooled buffer used automatically
}
```
Key behaviour: buffers retain their capacity when returned to the pool. A 512-byte buffer that grows to 3KB stays at 3KB capacity. This means subsequent renders reuse pre-grown buffers without reallocation.

Without a hint, renders still benefit from pooling - buffers are retrieved from the small pool, which over time will contain pre-grown buffers from previous renders.

#### BufferHint (Optional)

You can provide a `BufferHint(n)` to help determine which pool will be used when retrieving a bytes.Buffer and it will grow the bytes.Buffer to the appropriate hint. It returns the element, so it chains like any other method. After `Render(w)`, the element records the actual rendered size, which you can read with `RenderedSize()` and reuse:

```go
// First render - set a hint if you know approximate size
page := html.New(...).BufferHint(8192)  // Hint at 8KB
page.Render(w)

// Get the actual size for reuse on similar pages
actualSize := page.RenderedSize()  // e.g., 6543

// Use that hint for a new page with similar content
anotherPage := html.New(...).BufferHint(actualSize)
anotherPage.Render(w)
```

Buffers below the threshold (default 4KB) use the small pool; larger buffers use the large pool. This two-tier approach prevents small fragment renders from inheriting oversized buffers from full page renders.

```go
import "github.com/jpl-au/fluent/pool"

pool.SetThreshold(4096)               // Small vs large pool threshold (default 4KB)
pool.SetMaxPoolSize(262144)           // Max pooled size; oversized buffers are discarded (default 256KB)
pool.Disable()                        // Disable pooling entirely (pool.Enable() turns it back on)
```

For detailed mechanics and tuning guidance, see [AGENTS.md](AGENTS.md).

## Performance

The base Fluent API performs well out of the box with automatic buffer pooling. For high-throughput applications requiring additional optimisation, see [Fluent JIT](https://github.com/jpl-au/fluent-jit) which provides:

- **Compile** - Pre-render static portions, re-evaluate dynamic content via path navigation
- **Tune** - Adaptive buffer sizing that learns optimal sizes over time
- **Flatten** - Pre-render fully static content to raw bytes

For live updates, the same package provides the Differ and Memoiser diff engines, which consume the [reactive tracking](#reactive-tracking) hooks.

Build and test without JIT first.

## Generator

The `html5` and `dot` packages are generated from YAML definitions that follow the HTML5 specification. This keeps the API consistent with the spec and makes updates straightforward as HTML evolves.

The generator is still in early alpha, but I plan to release it allowing customisations to how you want to prioritise attributes. Figuring out how to create directives based on YAML structs and then use that to write out the files was largely an LLM-driven experience (thanks Claude Code). Having said that, it is also one of the reasons I am more tentative about putting the generator up as a repo given an LLM does not quite always follow its CLAUDE.md properly.

## Why Fluent (my motivations)

I created Fluent for a few reasons:

### Intellectual Curiosity

Fluent has been a great learning perspective, if nothing else. Building something from scratch is a great way to push your understanding of any programming language. From defining the architecture, experimenting with different packages, prototyping code... There is a lot of behind the scenes work. What you see today is the result of many months (March 2025) of building prototypes, testing, benchmarking, and experimentation.

During the course of it's design, I've built several discard prototypes before I settled on the architecture you see today. Some of my early prototype work involved the use of `strings.Builder` as part of the rendering pipeline, but I wasn't happy with the benchmark results. I did research into using a variety of alternatives (including some usage of `unsafe` pointers to manipulate the internals of some data structures) before I settled on the humble `bytes.Buffer`

Other prototypes focused on the use of generics and embedded structs, but the performance characteristics weren't quite what I had in mind. Exploring `sync.Pool` along with pprof helped me to analyse how the pool was working for me (or sometimes against me) but ultimately made me think of the two-pool aproach to cater for the small vs large buffer size (fragment vs. full page renders), as well as discarding over-sized buffers. It's been a fun and rewarding experience.

### I didn't like the alternatives

While I've built a few personal projects with [gomponents](https://github.com/maragudk/gomponents) and it is in all honesty the original inspiration for Fluent.

I'm not a fan of dot imports personally, but I know some developers prefer the syntax they provide. Fluent also includes the `dot` package as an optional way to interact with Fluent, but you still need to use the Fluent API regardless of which style you choose.

I also did not enjoy functions as arguments vs. the fluent API approach. It just felt awkward to me that I need to remember the function-as-attributes required to work with gomponents vs. letting the IDE give me the list of attributes (and the safety in knowing I cannot add the wrong attribute to the wrong element unless I choose to specifically override it - which Fluent allows).

```go
func Card(title, text string) Node {
	return Div(Class("card"),
		H2(Class("card-title"), g.Text(title)),
		P(Class("card-text"), g.Text(text)),
	)
}
```

The same component in Fluent:

```go
func Card(title, text string) node.Node {
	return div.New(
		h2.Text(title).Class("card-title"),
		p.Text(text).Class("card-text"),
	).Class("card")
}
```

I know Fluent's approach leads to a more verbose import declaration area, but `goimports` exists and can automatically handle this (as can your IDE). There are trade-offs to either approach, and I cannot say one is better than the other.

Another framework I looked into quite a while into the development of Fluent is [hb](https://github.com/dracory/hb) - and it is in many ways practically similar in syntax to Fluent. As I'd already started with Fluent, it gave me an alternate framework to work against in my internal benchmarking. I also don't think it is great that you have to import all extensions (htmx, Alpine, Swal, ...) as I always prefer an opt-in approach that tries to keep your code lean.

Perhaps the most similar framework to Fluent is [gostar](https://github.com/delaneyj/gostar) - which also uses a fluent API style with method chaining, a generator, and follows the HTML5 spec.

I have also worked with [Templ](https://github.com/a-h/templ) and while it's great, the pre-compile step just feels awkward to me, and ultimately led me to search for alternatives.

### Benchmark Performance

During the creation of Fluent I ran several benchmarks against gomponents, gostar, hb and even templ. In comparison with the non-compiled (i.e.: not templ) solutions, Fluent seems to have better CPU and memory profiles, with significantly lower allocations due to the buffer pooling strategy. The Fluent JIT package further optimises the performance characteristics. I decided against publishing the results as benchmarking can be subjective, and the results vary depending on how and what you are measuring. I welcome the opportunity for others to create their own benchmarks and share them.

## Profile-Guided Optimisation (PGO)

Go supports [Profile-Guided Optimisation](https://go.dev/doc/pgo) from version 1.21+. PGO uses a CPU profile from your running application to make more aggressive inlining and optimisation decisions at compile time, with no code changes.

The profile records which functions your application runs most. The compiler then optimises those call paths, which include Fluent's rendering pipeline, its buffer pooling, and any JIT strategy you use. PGO improves speed only. It does not change the number of allocations. Collect fresh profiles periodically as your application evolves.

## Ecosystem

Fluent has companion packages that extend its capabilities:

| Package | Description |
|---------|-------------|
| [Flint](https://github.com/jpl-au/flint) | Linter and introspection CLI for Fluent. Catches hallucinated APIs, unsafe `Static()`/`RawText()`, missed typed constructors, and raw strings where typed constants are required. `flint -info <element>` prints the full registry entry for any element. |
| [Fluent Security](https://github.com/jpl-au/fluent-security) | Opt-in security toolkit. Wraps [bluemonday](https://github.com/microcosm-cc/bluemonday) with Fluent-native helpers (`HTML`, `PlainText`) and a chainable `Cleaner` (`New`, `RichText`, `FromPolicy` + `Allow`/`AllowClasses`/`AllowAttr`) for sanitising untrusted HTML, plus `Nonce()` for Content-Security-Policy workflows with inline `<script>`/`<style>`. |
| [Fluent JIT](https://github.com/jpl-au/fluent-jit) | Performance optimisation with three strategies: **Compile** (pre-render static portions), **Tune** (adaptive buffer sizing), **Flatten** (pre-render fully static content to raw bytes). Also provides two diff engines for live updates: the **Differ** (targeted patches) and the **Memoiser** (subtree skipping by version). |
| [Fluent HTMX](https://github.com/jpl-au/fluent-htmx) | HTMX integration. Accepts `node.Element` to set HTMX attributes (`hx-get`, `hx-post`, `hx-swap`, etc.) on any Fluent element. |
| [Fluent Datastar](https://github.com/jpl-au/fluent-datastar) | [Datastar](https://data-star.dev) integration, targeting Datastar v1.0.2. Wrap any `node.Element` with `datastar.New()` to add `data-*` attributes through method chaining. Includes a server-side SSE generator for patching elements and signals. Covers the free, open-source attributes; the Pro attributes are not implemented. |
| [Tether](https://github.com/jpl-au/tether) | Reactive server-driven UI. Connects Fluent node trees to the browser over WebSocket, with an SSE fallback, and sends only targeted patches when state changes, so the client morphs the DOM in place and keeps input focus, scroll position, and form state. Three update modes: server rendering, signals, and client directives. The API is not yet stable. |

All companion packages are optional. Fluent works standalone for static HTML generation.

## Licence

MIT
