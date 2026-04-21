# Fluent

HTML5 components in Go using a Fluent API.

**Quick links:** [Why Fluent?](#why-fluent) · [Install](#install) · [Quick Start](#quick-start) · [Flint (linter & element info)](#flint) · [Reserved Keywords](#reserved-keywords) · [Static vs Dynamic Content](#static-vs-dynamic-content) · [Typed Constructors](#typed-constructors) · [Conditional Rendering](#conditional-rendering) · [Functional Processing](#functional-processing) · [Building Components](#building-components) · [Type-Safe Attributes](#type-safe-attributes) · [Architecture](#architecture) · [Performance](#performance) · [Generator](#generator) · [Ecosystem](#ecosystem) · [PGO](#profile-guided-optimization-pgo)

## Why Fluent?

**No template language to learn.** Write HTML using Go code. Get IDE auto-completion, type checking, and refactoring support for free. [AGENTS.md](AGENTS.md) makes it trivial for AI agents to do the hard work for you.

**Built for developers.** Thoughtful around the developer experience: attributes use native Go types - set a `width` with an `int`, a `volume` with a `float64`. Fluent handles the conversion. Type-safe constants for enumerated values catch typos like `type="emial"`.

**Type-safe nesting.** Typed constructors enforce correct HTML parent-child relationships at compile time. `ul.Items()` only accepts `*li.Element`, `tr.Cells()` only accepts `*td.Element` - the compiler catches nesting mistakes that would otherwise become silent bugs. `New()` remains available as the flexible escape hatch. [See Typed Constructors](#typed-constructors).

**HTML escaping by default.** `Text()` and `Textf()` automatically escape `<`, `>`, `&`, and quotes. For untrusted HTML that needs to render *as* HTML (rendered markdown, rich-text input), reach for the opt-in [fluent-security](https://github.com/jpl-au/fluent-security) package, which wraps [bluemonday](https://github.com/microcosm-cc/bluemonday) and returns Fluent nodes directly.

**Performance considered.** Buffer pooling and efficient rendering for high-throughput applications. Don't want to use `sync.Pool`? Just turn it off.

**Extensible.** Interface approach with methods to work with the underlying attributes allows any element in Fluent to be extended to work with any framework (htmx, Turbo) or to rewrite elements entirely to build web components.

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

[Flint](https://github.com/jpl-au/flint) is the companion linter and introspection CLI for Fluent. It catches mistakes before they compile, and doubles as a command-line reference for every Fluent element.

```bash
go install github.com/jpl-au/flint/cmd/flint@latest
```

Lint your code:

```bash
flint ./...              # check all Go files recursively
flint ./views            # check a specific directory
flint views/home.go      # check a single file
```

Look up an element's full API - constructors, typed methods, valid children, attribute mappings, and the typed constants each method accepts:

```bash
flint -info div          # everything about <div>
flint -info input        # everything about <input> (including every typed constant)
flint -info ol           # list constructors and typed variants
```

### What Flint catches

- **Hallucinated APIs.** Every function, method, and type reference is validated against a generated Fluent registry. Catches `node.Fragment()`, `.Href()` on the wrong element, `inputtype.Telephone` (does not exist), and similar.
- **Typed constant enforcement.** Flags raw strings passed to methods that require typed constants. `input.New().Type("email")` is flagged with a fix pointing at `inputtype.Email`.
- **Unsafe `Static()` and `RawText()`.** `Static()` requires a string literal so the JIT can pre-render it. `RawText()` does not HTML-escape, so dynamic values risk XSS. Both are flagged when called with a variable.
- **`New().Method()` redundancy.** Suggests the direct constructor: `div.New().Text("x")` should be `div.Text("x")`.
- **Typed constructor opportunities.** Suggests `ul.Items(...)` over `ul.New(li...)`, `tr.Cells(...)` over `tr.New(td...)`, and similar type-safe alternatives when children are uniform.
- **`SetAttribute()` misuse.** Flags chaining after `SetAttribute()` (it returns void) and flags usage where a typed method exists (`.Class()` instead of `.SetAttribute("class", ...)`).
- **Reserved keyword imports.** Points to the correct Fluent package for HTML elements that collide with Go keywords (`select` → `dropdown`, `main` → `primary`, `var` → `variable`).

Every diagnostic includes a `fix:` field with the corrected code. That makes Flint especially valuable alongside AI-generated code - the agent can read the fix and self-correct without human intervention.

## Reserved Keywords

Some HTML elements conflict with Go reserved keywords. I chose names that still feel intuitive - `dropdown` for `<select>` felt natural since that's what it renders.

| HTML Element | Fluent Package |
|--------------|----------------|
| `<select>`   | `dropdown`     |
| `<main>`     | `primary`      |
| `<var>`      | `variable`     |

## Documentation for AI agents

- `AGENTS.md` - Comprehensive guide to help AI agents work with Fluent (but it is also useful for humans who want a deeper dive into Fluent too)

## Static vs Dynamic Content

`Static()`, `Text()`, `Textf()`, `RawText()` and `RawTextf()` are both element constructors and methods. This was largely a requirement for the developer experience when using dot imports (which are optional).

The `Static()` constructor/method exists for use with [Fluent JIT](#performance) as a signal that this element is static. You should not use a variable in `Static()` as the JIT compiler will only render this node once (the first run) and subsequent calls will ignore changes. An alternative was to mark the node as dynamic, but I thought the developer experience would be hindered.

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

### Reactive tracking

Elements can be marked for reactive tracking with `.Dynamic("key")`, and subtrees can be memoised with `node.Memoise(key, func)`. These are used by [Fluent JIT](https://github.com/jpl-au/fluent-jit) for targeted diffing and by [Tether](https://github.com/jpl-au/tether) for live DOM patching. See the [Fluent JIT documentation](https://github.com/jpl-au/fluent-jit) for details.

```go
span.Textf("Count: %d", count).Dynamic("count")  // tracked by key
node.Memoise(version, func() node.Node {           // skipped when key unchanged
    return expensiveRender()
})
```

Many elements have convenience constructors for common use cases:

```go
// Form shortcuts
form.Get("/search", ...)   // <form action="/search" method="get">
form.Post("/login", ...)   // <form action="/login" method="post">

// Input types - all return *element for chaining
input.Email("email")             // <input type="email" name="email" />
input.Password("password")       // <input type="password" name="password"/>
input.Checkbox("agree", "yes")   // <input type="checkbox" name="agree" value="yes" />
input.Submit("Submit")           // <input type="submit"value="Submit"  />

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
figure.Caption("An elephant", img.New().Src("elephant.jpg"))
dl.Pair("Name", "Alice")

// New() is always available as the untyped escape hatch
table.New(caption.Text("Title"), thead.New(...), tbody.New(...))
```

## Conditional Rendering

`node.Condition()` provides inline conditional rendering. `True()` and `False()` can be used together or independently:

```go
// Both branches
node.Condition(user.IsLoggedIn).
    True(p.Text("Welcome back!")).
    False(a.New().Href("/login").Text("Sign in"))
```

For single-branch conditions, `When()` and `Unless()` provide concise shorthand:

```go
// Render only when condition is true
node.When(user.IsAdmin, span.Static("Admin"))

// Render only when condition is false
node.Unless(user.IsLoggedIn, a.New().Href("/login").Text("Sign in"))
```

Conditions can be nested since `node.Condition()` returns a `node.Node`:

```go
node.Condition(user.IsLoggedIn).
    True(
        node.Condition(user.IsAdmin).
            True(span.Static("Admin Dashboard")).
            False(span.Static("User Dashboard")),
    ).
    False(a.New().Href("/login").Text("Sign in"))
```

For multiple branches, `node.Func()` is cleaner:

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

## Functional Processing

For logic more complex than a simple condition - database lookups, error handling, or building dynamic content - `node.Func()` lets you write arbitrary Go code inline. The function executes at render time, keeping your tree structure declarative while deferring complex logic until it's needed.

```go
node.Func(func() node.Node {
    count, err := db.GetUnreadCount(userID)
    if err != nil || count == 0 {
        return nil
    }
    return div.Textf("You have %d unread messages", count).Class("notification")
})
```

The second form returns a slice `[]node.Node`:

```go
node.Funcs(func() []node.Node {
    items := make([]node.Node, len(products))
    for i, product := range products {
        items[i] = li.New(
            span.Text(product.Name),
            span.Textf("$%.2f", product.Price),
        )
    }
    return items
})
```

## Building Components

There's no special component system to learn - building your own components is handled through Go functions. You get all the benefits of Go's type system, testing, and refactoring tools. Components are just functions that return `node.Node` or a concrete element type:

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

```go
func UserGreeting(user User) node.Node {
    return div.New(
        img.New().Src(user.Avatar).Alt(user.Name),
        h3.Text(user.Name),
        node.Condition(user.IsAdmin).
            True(span.Static("Admin")).
            False(nil),
    ).Class("user-greeting")
}

// Use them like any other element
page := div.New(
    Card("Welcome", "Thanks for signing up!"),
    UserGreeting(currentUser),
)
page.Render(w)
```

## Type-Safe Attributes

Fluent provides type-safe constants for attributes with enumerated values. I wanted the IDE to do the heavy lifting. When you type `inputtype.`, your editor shows you every valid option - no more checking MDN to remember if it's `datetime-local` or `datetimeLocal`.

Methods like `InputType()` accept typed constants, not strings - so `input.New().InputType("emial")` won't compile. Each attribute package also provides a `Custom()` function for edge cases or future HTML specifications not yet covered.

```go
import (
    "github.com/jpl-au/fluent/html5/input"
    "github.com/jpl-au/fluent/html5/attr/inputtype"
    "github.com/jpl-au/fluent/html5/attr/autocomplete"
)

input.New().
    InputType(inputtype.Email).       // Typed constant, not a string
    AutoComplete(autocomplete.Email). // IDE shows all valid options
    Required()

// For edge cases or future specs
input.New().InputType(inputtype.Custom("future-type"))
```

## Architecture

Fluent is organised into several packages:

| Package | Description |
|---------|-------------|
| `node` | Core `Node` interface for all renderable types: `Render()`, `RenderBuilder()`, `Nodes()`. The `Element` interface extends `Node` with `SetAttribute()`, `RenderOpen()`, `RenderClose()` for HTML elements |
| `html5/*` | HTML5 elements, one package per element (e.g., `div`, `span`, `input`). Each provides `New()`, `Text()`, `Static()` constructors |
| `html5/attr/*` | Type-safe attribute constants (e.g., `inputtype.Email`, `autocomplete.Off`, `rel.Stylesheet`) |
| `text` | Text node implementations for `Static()`, `Text()`, `RawText()` and their formatted variants |
| `pool` | Buffer pooling configuration |
| `dot` | Optional dot import for cleaner syntax without package prefixes |

### Everything is a Node

The `node.Node` interface is the foundation of Fluent. Every renderable piece of content implements it: HTML elements, text nodes, conditionals (`node.Condition`), and function wrappers (`node.Func`). This unified interface enables arbitrary composition - any `node.Node` can be a child of any element.

HTML elements also implement `node.Element`, which extends `Node` with `SetAttribute()`, `RenderOpen()`, and `RenderClose()`. Text nodes, function components, and conditionals are not elements - they don't have attributes or tags.

When in doubt about return types for your components, `node.Node` is always safe:

```go
func MyComponent(showHeader bool) node.Node {
    if showHeader {
        return header.New(h1.Text("Welcome"))
    }
    return nil  // nil nodes are safely skipped during rendering
}
```

Returning concrete types (like `*div.Element`) allows method chaining after the call, but `node.Node` provides maximum flexibility when your component might return different element types or nil.

### Rendering

All nodes implement `node.Node`. Call `Render()` to get `[]byte`, or pass an `io.Writer` to write directly:

```go
// Get bytes
html := page.Render()

// Write to response
page.Render(w)
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

I have put it through my own benchmarking against other Go packages (Templ, Gomponents, hb) and I have been ecstatic about its performance, but as benchmarking can be quite subjective depending on how you benchmark. I've decided to not put those results up here. It would be interesting for anyone interested to write some benchmarks and publish the results.

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

You can provide a `BufferHint()` to help determine which pool will be used when retrieving a bytes.Buffer and it will grow the bytes.Buffer to the appropriate hint. After `Render(w)`, the hint is updated to reflect the actual rendered size, which you can retrieve and reuse:

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

Buffers below the threshold (default 4KB) use the small pool; larger buffers use the large pool. This two-tier approach prevents small fragment renders from inheriting oversized buffers from full page renders.

```go
import "github.com/jpl-au/fluent/pool"

pool.SetThreshold(4096)               // Small vs large pool threshold (default 4KB)
pool.SetMaxPoolSize(262144, true)     // Max pooled size, discard oversized (default 256KB)
pool.SetEnabled(false)                // Disable pooling entirely
```

For detailed mechanics and tuning guidance, see [AGENTS.md](AGENTS.md).

## Performance

The base Fluent API performs well out of the box with automatic buffer pooling. For high-throughput applications requiring additional optimisation, see [Fluent JIT](https://github.com/jpl-au/fluent-jit) which provides:

- **Compile** - Pre-render static portions, re-evaluate dynamic content via path navigation
- **Tune** - Adaptive buffer sizing that learns optimal sizes over time
- **Flatten** - Pre-render fully static content to raw bytes

Build and test without JIT first - premature optimisation is the root of all evil.

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

I also did not enjoy functions as arguments vs. the fluent API approach. It just felt awkward to me that I need to remember the function-as-attributes required to work with gomponents vs. letting the IDE give me the list of attributes (and the saftey in knowing I cannot add the wrong attribute to the wrong element unless I choose to specifically override it - which Fluent allows).

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

Perhaps the most similar framework to Fluent is [gostar](https://github.com/puregarlic/gostar) - which also uses a fluent API style with method chaining, a generator, and follows the HTML5 spec.

I have also worked with [Templ](https://github.com/a-h/templ) and while it's great, the pre-compile step just feels awkward to me, and ultimately led me to search for alternatives.

### Benchmark Performance

During the creation of Fluent I ran several benchmarks against gocomponents, gostar, hb and even templ. In comparison with the non-compiled (i.e.: not templ) solutions, Fluent seems to have better CPU and memory profiles, with significantly lower allocations due to the buffer pooling strategy. The Fluent JIT package further optimises the performance characteristics. I decided against publishing the results as benchmarking can be subjective, and the results vary depending on how and what you are measuring. I welcome the opportunity for others to create their own benchmarks and share them.

## Profile-Guided Optimization (PGO)

Go supports [Profile-Guided Optimization](https://go.dev/doc/pgo) from version 1.21+. PGO uses a CPU profile from your running application to make more aggressive inlining and optimisation decisions at compile time. Benchmarks show **10-20% speed improvements** with no code changes.

To enable PGO in your application:

1. Add profiling to your app (e.g. `import _ "net/http/pprof"`)
2. Collect a CPU profile under realistic load:
   ```bash
   curl -o default.pgo http://localhost:8080/debug/pprof/profile?seconds=30
   ```
3. Place `default.pgo` in your main package directory
4. `go build` - PGO is applied automatically

The profile captures which functions are hot in *your* application, so the compiler optimises the specific call paths you actually use - including Fluent's rendering pipeline, buffer pooling, and any JIT strategies. Allocations are unaffected; PGO improves speed only.

Collect fresh profiles periodically as your application evolves. Profiles from one platform can optimise builds for another (e.g. a Linux profile can optimise a macOS build).

## Ecosystem

Fluent has companion packages that extend its capabilities:

| Package | Description |
|---------|-------------|
| [Flint](https://github.com/jpl-au/flint) | Linter and introspection CLI for Fluent. Catches hallucinated APIs, unsafe `Static()`/`RawText()`, missed typed constructors, and raw strings where typed constants are required. `flint -info <element>` prints the full registry entry for any element. |
| [Fluent Security](https://github.com/jpl-au/fluent-security) | Opt-in security toolkit. Wraps [bluemonday](https://github.com/microcosm-cc/bluemonday) with Fluent-native helpers (`HTML`, `PlainText`) and a chainable `Cleaner` (`New`, `RichText`, `FromPolicy` + `Allow`/`AllowClasses`/`AllowAttr`) for sanitising untrusted HTML, plus `Nonce()` for Content-Security-Policy workflows with inline `<script>`/`<style>`. |
| [Fluent JIT](https://github.com/jpl-au/fluent-jit) | Performance optimisation with three strategies: **Compile** (pre-render static portions), **Tune** (adaptive buffer sizing), **Flatten** (pre-render fully static content to raw bytes). Also provides the **Diff** engine for reactive updates. |
| [Fluent HTMX](https://github.com/jpl-au/fluent-htmx) | HTMX integration. Accepts `node.Element` to set HTMX attributes (`hx-get`, `hx-post`, `hx-swap`, etc.) on any Fluent element. |
| [Tether](https://github.com/jpl-au/tether) | Server-driven reactive UI. Manages sessions, WebSocket transport, and a client-side runtime that applies targeted DOM patches using the JIT diff engine. Mark elements with `.Dynamic("key")` and Tether handles the rest. |

All companion packages are optional. Fluent works standalone for static HTML generation.

## Licence

MIT
