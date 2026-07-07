package node

import "testing"

func TestEscapeAttribute(t *testing.T) {
	cases := []struct {
		name, in, want string
	}{
		{"clean", "some class", "some class"},
		{"empty", "", ""},
		{"double quote", `a"b`, "a&#34;b"},
		{"single quote", "a'b", "a&#39;b"},
		{"ampersand", "a&b", "a&amp;b"},
		{"less than", "a<b", "a&lt;b"},
		{"greater than", "a>b", "a&gt;b"},
		{"nul", "a\x00b", "a�b"},
		{"plus not escaped", "a+b", "a+b"},
		{"breakout attempt", `red" onclick="steal()`, "red&#34; onclick=&#34;steal()"},
		{"all specials", "\"'&<>", "&#34;&#39;&amp;&lt;&gt;"},
		{"leading special", "&start", "&amp;start"},
		{"trailing special", "end<", "end&lt;"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := EscapeAttribute(c.in); got != c.want {
				t.Errorf("EscapeAttribute(%q) = %q, want %q", c.in, got, c.want)
			}
		})
	}
}

// TestEscapeAttributeCleanNoAlloc pins the fast path: a value with no
// escapable character is returned unchanged with zero allocations.
func TestEscapeAttributeCleanNoAlloc(t *testing.T) {
	clean := "a fairly long but perfectly clean attribute value 0123456789"
	if n := testing.AllocsPerRun(100, func() { _ = EscapeAttribute(clean) }); n != 0 {
		t.Errorf("clean-value escape allocated %v times, want 0", n)
	}
}

func BenchmarkEscapeAttributeClean(b *testing.B) {
	s := "card highlighted col-6 responsive"
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkStr = EscapeAttribute(s)
	}
}

func BenchmarkEscapeAttributeDirty(b *testing.B) {
	s := `red" onclick="steal(document.cookie)" data-x="`
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkStr = EscapeAttribute(s)
	}
}

var sinkStr string
