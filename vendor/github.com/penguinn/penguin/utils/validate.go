package utils

import (
	"unicode"
	"unicode/utf8"
)

func ISHan(str string) bool {
	strRune, _ := utf8.DecodeRuneInString(str)
	return unicode.Is(unicode.Han, strRune)
}
