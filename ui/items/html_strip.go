package items

import (
	"html"
	"regexp"
	"strings"
)

var (
	reBlock = regexp.MustCompile(`</(p|div|br|h[1-6]|li|tr)>`)
	reTags  = regexp.MustCompile(`<[^>]*>`)
)

func StripHTML(raw string) string {
	out := reBlock.ReplaceAllString(raw, "\n")
	out = reTags.ReplaceAllString(out, "")
	out = html.UnescapeString(out)
	out = regexp.MustCompile(`\n{3,}`).ReplaceAllString(out, "\n\n")
	return strings.TrimSpace(out)
}
