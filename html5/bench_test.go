package html5_test

import (
	"io"
	"testing"

	"github.com/jpl-au/fluent/html5/div"
	"github.com/jpl-au/fluent/html5/li"
	"github.com/jpl-au/fluent/html5/svg"
	"github.com/jpl-au/fluent/html5/ul"
	"github.com/jpl-au/fluent/text"
)

// These benchmarks cover the construction and rendering shapes that a change
// to the element layout trades against each other: a text element, an empty
// container, a container given a separate text node, a list built from data,
// an SVG drawing, and a retained element rendered again. Compare all of them
// when changing how elements store their children; a gain on one shape has
// cost the others before.

var sink any

func BenchmarkTextElement(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sink = div.Text("hello")
	}
}

func BenchmarkContainer(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sink = div.New()
	}
}

func BenchmarkContainerWithTextNode(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sink = div.New(text.Text("hello"))
	}
}

func BenchmarkListRender(b *testing.B) {
	items := make([]int, 20)
	for i := range items {
		items[i] = i
	}
	b.ReportAllocs()
	for b.Loop() {
		sink = ul.ItemsOf(items, func(int) *li.Element { return li.Text("item") }).RenderBytes()
	}
}

func BenchmarkSVGRender(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		shapes := make([]svg.Shape, 10)
		for i := range shapes {
			shapes[i] = svg.Rect().X("0").Y("0").Width("10").Height("10").Fill("currentColor")
		}
		sink = svg.New(shapes...).ViewBox("0 0 100 100").RenderBytes()
	}
}

func BenchmarkRetainedWriteTo(b *testing.B) {
	el := div.Text("hello")
	b.ReportAllocs()
	for b.Loop() {
		if _, err := el.WriteTo(io.Discard); err != nil {
			b.Fatal(err)
		}
	}
}
