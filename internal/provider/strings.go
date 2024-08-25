package provider

import (
	"unicode"
	"unicode/utf8"
)

func uppercaseFirstCharacter(s string) string {
	if s == "" {
		return s
	}

	// Convert the first character to uppercase
	r, size := utf8.DecodeRuneInString(s)
	if r == utf8.RuneError && size == 1 {
		return s
	}

	upper := unicode.ToUpper(r)
	if upper == r {
		return s
	}

	// Reconstruct the string with the first character uppercase
	return string(upper) + s[size:]
}
