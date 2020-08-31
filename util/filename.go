package util

import (
	"regexp"
)

var fileEscapeRe = regexp.MustCompile(`[/\\:\?\*]`)

func SanitizeFilename(name string) string {
	return fileEscapeRe.ReplaceAllString(name, "-")
}
