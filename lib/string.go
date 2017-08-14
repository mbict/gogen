package lib

import (
	"bytes"
	"fmt"
	"strconv"
	"unicode"
	"unicode/utf8"
)

// UnTitle Make the first charater lower case
func UnTitle(s string) string {
	if s == "" {
		return ""
	}
	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToLower(r)) + s[n:]
}

// SnakeCase produces the snake_case of a CamelCased string
func SnakeCase(name string) string {
	return xCase(name, '_')
}

// KebabCase produces the kebab-case of a CamelCased string
func KebabCase(name string) string {
	return xCase(name, '-')
}

func xCase(name string, replace rune) string {
	var b bytes.Buffer
	var lastUnderscore bool
	ln := len(name)
	if ln == 0 {
		return ""
	}
	b.WriteRune(unicode.ToLower(rune(name[0])))
	for i := 1; i < ln; i++ {
		r := rune(name[i])
		nextIsLower := false
		if i < ln-1 {
			n := rune(name[i+1])
			nextIsLower = unicode.IsLower(n) && unicode.IsLetter(n)
		}
		if unicode.IsUpper(r) {
			if !lastUnderscore && nextIsLower {
				b.WriteRune(replace)
				lastUnderscore = true
			}
			b.WriteRune(unicode.ToLower(r))
		} else {
			b.WriteRune(r)
			lastUnderscore = false
		}
	}
	return b.String()
}

func PadLeft(l int, str string) string {
	return fmt.Sprintf("%-"+strconv.Itoa(l)+"s", str)
}

func PadRight(l int, str string) string {
	return fmt.Sprintf("%"+strconv.Itoa(l)+"s", str)
}
