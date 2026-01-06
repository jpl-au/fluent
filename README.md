# Fluent

HTML5 components in Go using a Fluent API.

## Why Fluent?

**No template language to learn.** Write HTML using Go code. Get IDE auto-completion, type checking, and refactoring support for free. LLM-GUIDE.md makes it trivial for LLM's to do the hard work for you.

**Built for developers.** Thoughtful around the developer experience: attributes use native Go types - set a `width` with an `int`, a `volume` with a `float64`. Fluent handles the conversion. Type-safe constants for enumerated values catch typos like `type="emial"`.

**HTML escaping by default.** `Text()` and `Textf()` automatically escape `<`, `>`, `&`, and quotes. For content in `<script>` or `<style>` blocks, use the `security` package for additional sanitisation.

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

## Documentation for LLM's

- `LLM-GUIDE.md` - Comprehensive guide to help LLM's to work with Fluent (but it is also useful for humans who want a deeper dive into Fluent too)

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
    AutoComplete(autocomplete.Email)
```

## Reserved Keywords

Some HTML elements conflict with Go reserved keywords. I chose names that still feel intuitive - `dropdown` for `<select>` felt natural since that's what it renders.

| HTML Element | Fluent Package |
|--------------|----------------|
| `<select>`   | `dropdown`     |
| `<main>`     | `primary`      |
| `<var>`      | `variable`     |

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
node.FuncNodes(func() []node.Node {
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

*I was really at a stretch for what to call the method that returns a slice of node.Node `FuncNode` and if anything changes in the API, it will be this.*

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
| `node` | Core `Node` interface that all elements implement: `Render()`, `RenderBuilder()`, `Nodes()`, `SetAttribute()` |
| `html5/*` | HTML5 elements, one package per element (e.g., `div`, `span`, `input`). Each provides `New()`, `Text()`, `Static()` constructors |
| `html5/attr/*` | Type-safe attribute constants (e.g., `inputtype.Email`, `autocomplete.Off`, `rel.Stylesheet`) |
| `text` | Text node implementations for `Static()`, `Text()`, `RawText()` and their formatted variants |
| `pool` | Buffer pooling configuration |
| `security` | Sanitisation for `<script>` and `<style>` block content |
| `dot` | Optional dot import for cleaner syntax without package prefixes |

### Everything is a Node

The `node.Node` interface is the foundation of Fluent. Every renderable piece of content implements it: HTML elements, text nodes, conditionals (`node.Condition`), and function wrappers (`node.Func`). This unified interface enables arbitrary composition - any `node.Node` can be a child of any element.

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

For detailed mechanics and tuning guidance, see [LLM-GUIDE.md](LLM-GUIDE.md#how-the-two-tier-pool-works).

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

I have also worked with [Templ](https://github.com/a-h/templ) and while it's great, the pre-compile step just feels awkward to me, and ultimately led me to search for alternatives.

## Licence

MIT