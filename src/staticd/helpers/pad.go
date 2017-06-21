package helpers

import (
	"strings"
)

func PadLink(caption string, uri string, length int) string {
	padCount := 1 + length - len(caption)
	return `<a href="` + uri + `">` + caption + `</a>` + strings.Repeat(" ", padCount)
}

func PadText(text string, length int) string {
	padCount := 1 + length - len(text)
	return text + strings.Repeat(" ", padCount)
}
