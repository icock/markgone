// This source code is dedicated to the Public Domain.

package markgone

import (
	"os"
	"testing"
	"strings"
)

var text string = `First line as the title

MarkGone uses the minimal formatting syntax from godoc.
MarkGone processes plain text instead of comments in go source code.

Paragraphs are separated by one or more blank lines.

Heading

A heading is a single line
followed by another paragraph,
beginning with a capital letter,
and containing no punctuation.

    To produce a pre-formatted blocks,
    simply indent every line of the block.
    Common indent prefix will be removed in output.

http://example.com/urls_are_auto-linked

Tags: space-separated tags on last line following at least one blank-line`

func ExampleDoc() {
	ToHTML(os.Stdout, text, nil)
	// Output:
	/*
<h2>First line as the title</h2>
<p>
MarkGone uses the minimal formatting syntax from godoc.
MarkGone processes plain text instead of comments in go source code.
</p>
<p>
Paragraphs are separated by one or more blank lines.
</p>
<h3 id="hdr-Heading">Heading</h3>
<p>
A heading is a single line
followed by another paragraph,
beginning with a capital letter,
and containing no punctuation.
</p>
<pre>To produce a pre-formatted blocks,
simply indent every line of the block.
Common indent prefix will be removed in output.
</pre>
<p>
<a href="http://example.com/urls_are_auto-linked">http://example.com/urls_are_auto-linked</a></p>
<div class="taglist">
<strong>Tags:</strong>
<a href="/tag/space-separated" rel="tag">space-separated</a>
<a href="/tag/tags" rel="tag">tags</a>
<a href="/tag/on" rel="tag">on</a>
<a href="/tag/last" rel="tag">last</a>
<a href="/tag/line" rel="tag">line</a>
<a href="/tag/following" rel="tag">following</a>
<a href="/tag/at" rel="tag">at</a>
<a href="/tag/least" rel="tag">least</a>
<a href="/tag/one" rel="tag">one</a>
<a href="/tag/blank-line" rel="tag">blank-line</a>
</div>
	*/
}


var preparedBodyTests = []struct {
	in string
	out string
}{
	{"\na leading blank line", "a leading blank line"},
	{"\n\n\nleading blank lines", "leading blank lines"},
	{"a trailing blank line\n", "a trailing blank line"},
	{"trailing blank lines\n\n\n", "trailing blank lines"},
	{"a trailing space \n", "a trailing space"},
	{"trailing spaces    \n", "trailing spaces"},
	{"\n\n\nmultiple \n\n\nlines\n\n\n", "multiple\n\n\nlines"},
}
func TestPreparedBody(t *testing.T) {
	for i, tt := range preparedBodyTests {
		outLines, _ := preparedBody(strings.Split(tt.in, "\n"))
		out := strings.Join(outLines, "\n")
		if out != tt.out {
			t.Errorf("#%d FAIL\n  Actual value: %v\nExpected value: %v", i, out, tt.out)
		}
	}
}


var formattingTitleTests = []struct{
	in string
	out string
}{
	{
		"title\n",
		"<h2>title</h2>\n",
	},
	{
		"not a title (no newline)",
		"<p>\nnot a title (no newline)</p>\n",
	},
	{
		"another title\n\n",
		"<h2>another title</h2>\n",
	},
	{
		"title\n\nbody\n",
		"<h2>title</h2>\n<p>\nbody</p>\n",
	},
	{
		"two continuous lines\nare not title\n",
		"<p>\ntwo continuous lines\nare not title</p>\n",
	},
	{
		"\n\ntitle must be the first line\n",
		"<p>\ntitle must be the first line</p>\n",
	},
}
func TestFormattingTitle(t *testing.T) {
	for i, tt := range formattingTitleTests {
		out := ToHTMLString(tt.in, nil)
		if out != tt.out {
			t.Errorf("#%d FAIL\n  Actual value: %v\nExpected value: %v", i, out, tt.out)
		}
	}
}

var taggingTests = []struct{
	in string
	out string
}{
	{
		"\nparagraph\n\nTags: one-tag",
		`<p>
paragraph</p>
<div class="taglist">
<strong>Tags:</strong>
<a href="/tag/one-tag" rel="tag">one-tag</a>
</div>
`,
	},
	{
		"\nparagraph\n\nTags: multiple tags",
		`<p>
paragraph</p>
<div class="taglist">
<strong>Tags:</strong>
<a href="/tag/multiple" rel="tag">multiple</a>
<a href="/tag/tags" rel="tag">tags</a>
</div>
`,
	},
	{
		"\nparagraph\nTags: not following blank lines",
		"<p>\nparagraph\nTags: not following blank lines</p>\n",
	},
	{
		"\nparagraph\n\nTags:no space",
		"<p>\nparagraph\n</p>\n<p>\nTags:no space</p>\n",
	},
	{
		"\nempty list\n\nTags: ",
		"<p>\nempty list\n</p>\n<p>\nTags:</p>\n",
	},
}
func TestTagging(t *testing.T) {
	for i, tt := range taggingTests {
		out := ToHTMLString(tt.in, nil)
		if out != tt.out {
			t.Errorf("#%d FAIL\n  Actual value: %v\nExpected value: %v", i, out, tt.out)
		}
	}
}

