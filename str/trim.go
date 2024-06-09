package str

import (
	"strings"
)

const (
	whiteSpaceTrimSet = "\r\n\t "
)

// TrimRight trims the right part of a string from white spaces.
func TrimRight(s string) string {
	return strings.TrimRight(s, whiteSpaceTrimSet)
}

// TrimLeft trims the left part of a string from white spaces.
func TrimLeft(s string) string {
	return strings.TrimLeft(s, whiteSpaceTrimSet)
}
