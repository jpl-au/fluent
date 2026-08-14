package node

import (
	"strings"
	"sync/atomic"
)

// UnsafeURL is the value substituted for a URL whose scheme is not allowed.
// It is a fragment-only reference: inert for navigation and script, harmless
// in layout, and greppable in rendered output. It mirrors the mechanics of
// html/template's ZgotmplZ sentinel.
const UnsafeURL = "#fluent-unsafe-url"

// unsafeURLObserver holds the optional observer registered via OnUnsafeURL.
// It is read on the URL-setting path, so an atomic pointer keeps clean URLs
// lock-free (they never load it) and the rejection path a single atomic load.
var unsafeURLObserver atomic.Pointer[func(rejected string)]

// OnUnsafeURL registers fn to be called whenever [FilterURL] rejects a URL,
// passing the original (rejected) string. It is an observation hook for
// logging or metrics - security monitoring wants to know when a hostile URL
// was filtered - and it cannot change the outcome: the rejected URL is always
// replaced with [UnsafeURL] regardless. Pass nil to clear a previous observer.
//
// The observer runs synchronously on the goroutine that set the attribute, and
// only on the rejection path (a clean URL never calls it, so escaping-clean
// renders pay nothing). Registration is safe for concurrent use; because it is
// observe-only, it carries none of the mutable-policy risk of a runtime safety
// toggle.
func OnUnsafeURL(fn func(rejected string)) {
	if fn == nil {
		unsafeURLObserver.Store(nil)
		return
	}
	unsafeURLObserver.Store(&fn)
}

// FilterURL guards URL-valued attributes on sinks that navigate or load
// (href, src, action, ...). It returns s unchanged when the URL is safe, and
// [UnsafeURL] when the scheme is not on the allowlist.
//
// The allowlist is positive by design - only http, https, mailto, tel and sms
// (plus every relative URL, which carries no scheme) are permitted, and
// anything else is refused. An allowlist fails safe: an obfuscated or unknown
// scheme does not match the small permitted set and is rejected, so there is
// no denylist for an attacker to slip past with "java\tscript:" or casing or
// whitespace tricks. Rejecting an executable scheme (javascript:, vbscript:,
// data:) is exactly the point; rejecting a rare-but-legitimate one (a custom
// app scheme) is handled by the SetAttribute escaped-but-unfiltered path or by
// SetAttributeRaw.
//
// Leading and trailing ASCII whitespace and C0 controls are ignored for the
// scheme decision, matching how browsers trim a URL before parsing, so
// " javascript:x" is still rejected. Interior control characters are left in
// place, so "java\tscript:" yields a scheme that is not on the allowlist and
// is rejected. The returned value, when safe, is the original s untouched -
// filtering decides, it does not rewrite.
//
// FilterURL runs before [EscapeAttribute]; the survivor (or the sentinel, which
// contains nothing to escape) is then escaped for storage.
func FilterURL(s string) string {
	t := trimURLBytes(s)
	if t == "" {
		return s
	}

	colon := strings.IndexByte(t, ':')
	if colon < 0 {
		return s // no scheme: a relative URL
	}
	// A '/', '?' or '#' before the colon means the colon sits inside a path,
	// query or fragment, not a scheme - so this is a relative URL that merely
	// contains a colon (e.g. "?redirect=a:b", "#t:30"), which is safe.
	for i := range colon {
		switch t[i] {
		case '/', '?', '#':
			return s
		}
	}

	switch scheme := t[:colon]; {
	case strings.EqualFold(scheme, "http"),
		strings.EqualFold(scheme, "https"),
		strings.EqualFold(scheme, "mailto"),
		strings.EqualFold(scheme, "tel"),
		strings.EqualFold(scheme, "sms"):
		return s
	default:
		if obs := unsafeURLObserver.Load(); obs != nil {
			(*obs)(s)
		}
		return UnsafeURL
	}
}

// trimURLBytes returns s without leading or trailing ASCII whitespace and C0
// control bytes (<= 0x20), matching the WHATWG URL parser's trimming before
// scheme extraction. Interior bytes are untouched.
func trimURLBytes(s string) string {
	start, end := 0, len(s)
	for start < end && s[start] <= 0x20 {
		start++
	}
	for end > start && s[end-1] <= 0x20 {
		end--
	}
	return s[start:end]
}
