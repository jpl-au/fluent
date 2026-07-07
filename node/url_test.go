package node

import "testing"

func TestFilterURL(t *testing.T) {
	safe := []string{
		"/path",
		"./relative",
		"path/to/x",
		"?redirect=a:b",       // colon after '?' is not a scheme
		"#t:30",               // colon after '#' is not a scheme
		"//cdn.example.com/x", // protocol-relative
		"",                    // empty passes through
		"http://example.com",
		"https://example.com/a?b=c",
		"HTTPS://EXAMPLE.COM",   // scheme compare is case-insensitive
		"MailTo:john@x.com",     // allowed, any casing
		"tel:+61400000000",      // allowed
		"sms:+61400000000",      // allowed
		" https://example.com ", // surrounding whitespace trimmed for the decision
	}
	for _, s := range safe {
		if got := FilterURL(s); got != s {
			t.Errorf("FilterURL(%q) = %q, want it unchanged (safe)", s, got)
		}
	}

	unsafe := []string{
		"javascript:alert(1)",
		"JaVaScRiPt:alert(1)",   // casing does not help the attacker
		" javascript:alert(1)",  // leading space trimmed, then rejected
		"java\tscript:alert(1)", // interior tab: scheme not on allowlist, rejected
		"vbscript:msgbox(1)",
		"data:text/html,<script>alert(1)</script>",
		"data:image/png;base64,AAAA", // data: is not allowed on nav sinks
		"file:///etc/passwd",
		"ftp://host/file",
		"blob:https://x/uuid",
		"intent://scan/#Intent;scheme=x;end",
	}
	for _, s := range unsafe {
		if got := FilterURL(s); got != UnsafeURL {
			t.Errorf("FilterURL(%q) = %q, want %q (rejected)", s, got, UnsafeURL)
		}
	}
}

// TestFilterURLCleanNoAlloc pins that a safe URL is returned unchanged
// without allocating.
func TestFilterURLCleanNoAlloc(t *testing.T) {
	u := "https://example.com/assets/app.js?v=abc123"
	if n := testing.AllocsPerRun(100, func() { _ = FilterURL(u) }); n != 0 {
		t.Errorf("safe-URL filter allocated %v times, want 0", n)
	}
}

func BenchmarkFilterURLSafe(b *testing.B) {
	u := "https://example.com/path?query=value"
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkStr = FilterURL(u)
	}
}

func BenchmarkFilterURLReject(b *testing.B) {
	u := "javascript:alert(document.cookie)"
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkStr = FilterURL(u)
	}
}
