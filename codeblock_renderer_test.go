package glamour

import (
	"os"
	"strings"
	"testing"

	"charm.land/glamour/v2/styles"
	"github.com/charmbracelet/x/ansi"
)

func renderWith(t *testing.T, md string, opts ...TermRendererOption) string {
	t.Helper()
	opts = append([]TermRendererOption{WithStandardStyle(styles.DarkStyle), WithWordWrap(0)}, opts...)
	r, err := NewTermRenderer(opts...)
	if err != nil {
		t.Fatal(err)
	}
	out, err := r.Render(md)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

func TestCodeBlockRendererLayout(t *testing.T) {
	const art = "┌─┐\n│あ│\n└─┘\n"
	draw := func(language, code string) (string, bool) {
		if language != "art" {
			return "", false
		}
		return art, true
	}

	// A rendered block must sit exactly where a plain text block with the
	// same body would.
	tests := []struct {
		name     string
		rendered string
		plain    string
	}{
		{
			name:     "backtick fence",
			rendered: "before\n\n```art\nsource\n```\n\nafter\n",
			plain:    "before\n\n```text\n" + art + "```\n\nafter\n",
		},
		{
			name:     "tilde fence",
			rendered: "~~~art\nsource\n~~~\n",
			plain:    "~~~text\n" + art + "~~~\n",
		},
		{
			name:     "in a blockquote",
			rendered: "> quote\n>\n> ```art\n> source\n> ```\n",
			plain:    "> quote\n>\n> ```text\n> ┌─┐\n> │あ│\n> └─┘\n> ```\n",
		},
		{
			name:     "in a list",
			rendered: "- item\n\n  ```art\n  source\n  ```\n",
			plain:    "- item\n\n  ```text\n  ┌─┐\n  │あ│\n  └─┘\n  ```\n",
		},
		{
			name:     "without a trailing newline",
			rendered: "```art\nsource\n```\n",
			plain:    "```text\n" + art + "```\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ansi.Strip(renderWith(t, tt.rendered, WithCodeBlockRenderer(draw)))
			want := ansi.Strip(renderWith(t, tt.plain))
			if got != want {
				t.Errorf("got\n%q\nwant\n%q", got, want)
			}
		})
	}
}

func TestCodeBlockRendererTrimmedTrailingNewline(t *testing.T) {
	trimmed := func(string, string) (string, bool) { return "┌─┐\n└─┘", true }
	got := ansi.Strip(renderWith(t, "```art\nsource\n```\n", WithCodeBlockRenderer(trimmed)))
	want := ansi.Strip(renderWith(t, "```text\n┌─┐\n└─┘\n```\n"))
	if got != want {
		t.Errorf("got\n%q\nwant\n%q", got, want)
	}
}

func TestCodeBlockRendererDeclined(t *testing.T) {
	in, err := os.ReadFile(markdown)
	if err != nil {
		t.Fatal(err)
	}
	decline := func(string, string) (string, bool) { return "unused", false }
	got := renderWith(t, string(in), WithCodeBlockRenderer(decline))
	want := renderWith(t, string(in))
	if got != want {
		t.Error("declining every block changed the output")
	}
}

func TestCodeBlockRendererCalls(t *testing.T) {
	type call struct{ language, code string }
	var calls []call
	record := func(language, code string) (string, bool) {
		calls = append(calls, call{language, code})
		return "", false
	}

	md := strings.Join([]string{
		"```Mermaid title\ngraph TD\n```",
		"~~~ go\nfunc main() {}\n~~~",
		"```\nno language\n```",
		"    indented code",
	}, "\n\n")
	renderWith(t, md, WithCodeBlockRenderer(record))

	want := []call{
		{"Mermaid", "graph TD\n"},
		{"go", "func main() {}\n"},
	}
	if len(calls) != len(want) {
		t.Fatalf("got calls %q, want %q", calls, want)
	}
	for i := range want {
		if calls[i] != want[i] {
			t.Errorf("call %d = %q, want %q", i, calls[i], want[i])
		}
	}
}
