package form_test

import (
	"bytes"
	"testing"

	"github.com/jpl-au/fluent/html5/form"
	"github.com/jpl-au/fluent/html5/input"
)

func TestMultipartChildren(t *testing.T) {
	e := form.Multipart("/upload")
	wantEmpty := `<form action="/upload" method="post" enctype="multipart/form-data"></form>`
	if got := string(e.RenderBytes()); got != wantEmpty {
		t.Fatalf("got %q, want %q", got, wantEmpty)
	}
	e.Add(input.File("document"), input.Hidden("token", "abc"))
	want := `<form action="/upload" method="post" enctype="multipart/form-data"><input name="document" type="file"><input name="token" value="abc" type="hidden"></form>`
	if got := string(e.RenderBytes()); got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
	var split bytes.Buffer
	e.RenderOpen(&split)
	for _, child := range e.Nodes() {
		child.RenderBuilder(&split)
	}
	e.RenderClose(&split)
	if split.String() != want {
		t.Fatalf("split render got %q, want %q", split.String(), want)
	}
	e.Replace()
	if got := string(e.RenderBytes()); got != wantEmpty {
		t.Fatalf("after Replace got %q, want %q", got, wantEmpty)
	}
}
