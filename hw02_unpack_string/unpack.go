package hw02unpackstring

import (
	"errors"
)

var ErrInvalidString = errors.New("invalid string")

func runeIsNumber(r rune) bool {
	// руны - символы 48 - 0, 57 - 9
	return r >= 48 && r <= 57
}

func Unpack(str string) (string, error) {
	isPrevEscaped := false
	var prev rune
	var result string

	for _, i := range str {
		if prev == 0 {
			if runeIsNumber(i) {
				return "", ErrInvalidString
			}
			result += string(i)
			prev = i
			continue
		}

		if runeIsNumber(i) && runeIsNumber(prev) && !isPrevEscaped {
			return "", ErrInvalidString
		}

		switch {
		case runeIsNumber(i) && (prev != '\\' || prev == '\\' && isPrevEscaped):
			if i == rune('0') {
				asRune := []rune(result)
				result = string(asRune[0 : len(asRune)-1])
			}

			for j := 1; j < int(i)-'0'; j++ {
				result += string(prev)
			}
		case i != '\\':
			result += string(i)
		case i == '\\' && prev == '\\' && !isPrevEscaped:
			result += string(i)
		}

		isPrevEscaped = prev == '\\' && !isPrevEscaped
		prev = i
	}

	return result, nil
}
