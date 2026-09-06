package svg_test

import (
	"fmt"

	"github.com/jpl-au/fluent/html5/attr/compositeoperator"
	"github.com/jpl-au/fluent/html5/attr/units"
	"github.com/jpl-au/fluent/html5/svg"
)

// Title gives the circle a short name; Desc provides a longer description.
// Nesting them inside the circle associates the text with that shape.
func Example_labelledShape() {
	drawing := svg.New(
		svg.Circle(
			svg.Title("Q3 revenue"),
			svg.Desc("Revenue increased by 12 percent"),
		).Cx("50").Cy("50").R("20"),
	).ViewBox("0 0 100 100")

	fmt.Println(string(drawing.RenderBytes()))
	// Output:
	// <svg viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg"><circle cx="50" cy="50" r="20"><title>Q3 revenue</title><desc>Revenue increased by 12 percent</desc></circle></svg>
}

// Defs holds a gradient without drawing it. Both rectangles reference the same
// definition through their fill. ObjectBoundingBox scales the gradient to each
// rectangle, so each shows the full colour range despite their different widths.
func Example_gradient() {
	drawing := svg.New(
		svg.Defs(
			// Keep this ID unique in the containing document when embedding inline SVG.
			svg.LinearGradient(
				svg.Stop().Offset("0%").StopColor("#2563eb"),
				svg.Stop().Offset("100%").StopColor("#7c3aed"),
			).ID("revenue-gradient").GradientUnits(units.ObjectBoundingBox).
				X1("0%").Y1("0%").X2("100%").Y2("0%"),
		),
		svg.Rect().X("10").Y("10").Width("80").Height("30").Fill("url(#revenue-gradient)"),
		svg.Rect().X("10").Y("50").Width("50").Height("30").Fill("url(#revenue-gradient)"),
	).ViewBox("0 0 100 100")

	fmt.Println(string(drawing.RenderBytes()))
	// Output:
	// <svg viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg"><defs><linearGradient x1="0%" y1="0%" x2="100%" y2="0%" gradientUnits="objectBoundingBox" id="revenue-gradient"><stop offset="0%" stop-color="#2563eb"></stop><stop offset="100%" stop-color="#7c3aed"></stop></linearGradient></defs><rect x="10" y="10" width="80" height="30" fill="url(#revenue-gradient)"></rect><rect x="10" y="50" width="50" height="30" fill="url(#revenue-gradient)"></rect></svg>
}

// The circle defines a clipping region; it is not drawn as a visible circle.
// Applying its reference to G clips both rectangles. UserSpaceOnUse expresses
// the clipping geometry in the same coordinates as the group's contents.
func Example_clippingPath() {
	drawing := svg.New(
		svg.Defs(
			svg.ClipPath(
				svg.Circle().Cx("50").Cy("50").R("40"),
			).ID("chart-clip").ClipPathUnits(units.UserSpaceOnUse),
		),
		svg.G(
			svg.Rect().Width("100").Height("50").Fill("#2563eb"),
			svg.Rect().Y("50").Width("100").Height("50").Fill("#7c3aed"),
		).ClipPath("url(#chart-clip)"),
	).ViewBox("0 0 100 100")

	fmt.Println(string(drawing.RenderBytes()))
	// Output:
	// <svg viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg"><defs><clipPath clipPathUnits="userSpaceOnUse" id="chart-clip"><circle cx="50" cy="50" r="40"></circle></clipPath></defs><g clip-path="url(#chart-clip)"><rect width="100" height="50" fill="#2563eb"></rect><rect y="50" width="100" height="50" fill="#7c3aed"></rect></g></svg>
}

// A filter builds a shadow by blurring the shape's alpha channel, offsetting
// that result, and drawing the original graphic over it. Result names connect
// stages within this filter; unlike the filter's ID, they are not URL references.
// The filter region is expanded to leave room for the shadow beyond the shape.
func Example_filterChain() {
	drawing := svg.New(
		svg.Defs(
			svg.Filter(
				// SourceAlpha and SourceGraphic are predefined inputs for the filtered shape.
				svg.FeGaussianBlur().In("SourceAlpha").StdDeviation("2").Result("blur"),
				svg.FeOffset().In("blur").Dx("2").Dy("4").Result("shadow"),
				// In is the foreground; In2 is the background for the Over operator.
				svg.FeComposite().In("SourceGraphic").In2("shadow").Operator(compositeoperator.Over),
			).ID("chart-shadow").X("-20%").Y("-20%").Width("140%").Height("140%"),
		),
		svg.Rect().X("10").Y("10").Width("80").Height("80").
			Fill("#2563eb").Filter("url(#chart-shadow)"),
	).ViewBox("0 0 100 100")

	fmt.Println(string(drawing.RenderBytes()))
	// Output:
	// <svg viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg"><defs><filter x="-20%" y="-20%" width="140%" height="140%" id="chart-shadow"><feGaussianBlur in="SourceAlpha" stdDeviation="2" result="blur"></feGaussianBlur><feOffset in="blur" dx="2" dy="4" result="shadow"></feOffset><feComposite in="SourceGraphic" in2="shadow" operator="over"></feComposite></filter></defs><rect x="10" y="10" width="80" height="80" fill="#2563eb" filter="url(#chart-shadow)"></rect></svg>
}
