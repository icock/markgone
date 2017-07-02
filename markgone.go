// This source code is dedicated to the Public Domain.

/*
Package markgone provides a minimal markup language based on godoc.

MarkGone uses the minimal formatting syntax from godoc.
Unlike godoc, MarkGone processes plain text, not comments in go source code;

Philosophy

Markups are gone.
In other words, MarkGone text just looks like plain text.

Formatting Syntax

A heading is a single line
followed by another paragraph,
beginning with a capital letter,
and containing no punctuation.

Paragraphs are separated by one or more blank lines.

    To produce a pre-formatted blocks,
    simply indent every line of the block.
    Common indent prefix will be removed in output.

http://example.com/urls_are_auto-linked

Source Code HighLight

MarkGone does not highlight pre-formatted blocks.
The highlighting can be done via JavaScript.
For example, with highlight.js:

	<link rel="stylesheet" href="/path/to/styles/default.css">
	<script src="/path/to/highlight.pack.js"></script>
	<script>
	document.addEventListener('load', () => {
		hljs.configure({languages: ["go", "c", "scheme", "java", "js", "html", "css"]});
		const codes = document.querySelectorAll('pre');
		codes.forEach((code) => { hljs.highlightBlock(code); });
	});
	</script>

CSS

Since markgone only uses four HTML elements

	h2, h3, p, pre

customizing css style for it is simple.

Sample css files are provided in the css directory of source code repository:

	[]struct{
		filename string;
		description string;
	}{
		{"godoc.css":
			"mimic godoc style"},
		{"plain/blue-link":
			"mimic plain text with colored links"},
		{"plain/underline-link":
			"mimic plain text with underlined links"},
		{"plain/white-on-black":
			"mimic a console"},
		{"plain/green-on-black":
			"mimic a console with green text"},
	}

Extensions

If the first line is a single line, i.e. followed by a blank line,
then MarkGone uses it as the title.

Tags: space-separated tags on last line following at least one blank-line
 */
package markgone

import (
	"strings"
	"io"
	"go/doc"
	"html/template"
	"fmt"
	"bytes"
)

// ToHTML converts markgone text to formatted HTML.
// See doc.ToHTML for more information.
// See also ToHTMLString.
func ToHTML(w io.Writer, text string, words map[string]string) {
	title, body, tags := prepared(text)
	if title != "" {
		titleToHTML(w, title)
	}
	doc.ToHTML(w, strings.Join(body, "\n"), words)
	if tags != nil {
		tagsToHTML(w, tags)
	}
}

// ToHTMLString is like ToHTML, but returns formatted HTML as a string.
func ToHTMLString(text string, words map[string]string) string {
	var b bytes.Buffer
	ToHTML(&b, text, words)
	return b.String()
}

func prepared(text string) (title string, body []string, tags []string) {
	lines := strings.Split(text, "\n")
	length := len(lines)

	if length >= 2 && isTitle(lines[:2]) {
		title = lines[0]
		body, tags = preparedBody(lines[2:])
	} else {
		title = ""
		body, tags = preparedBody(lines)
	}

	return title, body, tags
}

func isTitle(lines []string) bool {
	if lines[0] != "" && lines[1] == "" {
		return true
	} else {
		return false
	}
}
func titleToHTML(w io.Writer, title string) {
	fmt.Fprint(w, "<h2>")
	template.HTMLEscape(w, []byte(title))
	fmt.Fprintln(w, "</h2>")
}

func preparedBody(lines []string) ([]string, []string) {
	lines = leadingBlankStripped(lines)
	lines = trailingBlankStripped(lines)
	lines = trailingSpacesStripped(lines)

	length := len(lines)
	if length <= 1 {
		return lines, nil
	} else {
		lastLine := lines[length-1]
		tagLinePrefix := "Tags: "
		if strings.HasPrefix(lastLine, tagLinePrefix) && lines[length-2] == "" {
			lastLine = lastLine[len(tagLinePrefix):]
			tags := strings.Split(lastLine, " ")
			tags = emptyStringsStripped(tags)
			return lines[:length-2], tags
		} else {
			return lines, nil
		}
	}
}
func leadingBlankStripped(lines []string) []string {
	if len(lines) == 0 {
		return lines
	} else if lines[0] != "" {
		return lines
	} else {
		return leadingBlankStripped(lines[1:])
	}
}
func trailingBlankStripped(lines []string) []string {
	length := len(lines)
	if length == 0 {
		return lines
	} else if lines[length-1] != "" {
		return lines
	} else {
		return trailingBlankStripped(lines[:length-1])
	}
}
func trailingSpacesStripped(lines []string) []string {
	for i, line := range lines {
		lines[i] = strings.TrimRight(line, " ")
	}
	return lines
}
func emptyStringsStripped(strings []string) []string {
	for i, str := range strings {
		if str == "" {
			strings = append(strings[:i], strings[i+1:]...)
		}
	}
	return strings
}

func tagsToHTML(w io.Writer, tags []string) {
	fmt.Fprintln(w, `<div class="taglist">`)
	fmt.Fprintln(w, "<strong>Tags:</strong>")
	for _, tag := range tags {
		tag = template.HTMLEscapeString(tag)
		fmt.Fprintf(w, "<a href=\"/tag/%s\" rel=\"tag\">%s</a>\n", tag, tag)
	}
	fmt.Fprintln(w, "</div>")
}
