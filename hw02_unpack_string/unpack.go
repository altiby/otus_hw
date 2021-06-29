package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

// ErrInvalidString error invalid string
var ErrInvalidString = errors.New("invalid string")

func Unpack(input string) (string, error) {
	builder := strings.Builder{}
	runes := []rune(input)
	lexes, err := getLexemes(runes)
	if err != nil {
		return "", err
	}
	for _, lex := range lexes {
		builder.WriteString(getStringFormLexeme(lex))
	}
	return builder.String(), nil
}

func isSlash(r rune) bool {
	return r == '\\'
}

func getStringFormLexeme(lex []rune) string {
	var r rune
	var counter int
	firstSlash := false
	if len(lex) > 0 && isSlash(lex[0]) {
		firstSlash = true
	}
	switch {
	case len(lex) == 0:
		return ""
	case len(lex) == 1:
		return string(lex[0])
	case len(lex) == 2 && firstSlash:
		return string(lex[1])
	case len(lex) == 2:
		r = lex[0]
		counter, _ = strconv.Atoi(string(lex[1]))
		if counter == 0 {
			return ""
		}
		return strings.Repeat(string(r), counter)
	case len(lex) == 3:
		r = lex[1]
		counter, _ = strconv.Atoi(string(lex[2]))
		if counter == 0 {
			return ""
		}
		return strings.Repeat(string(r), counter)
	default:
		return ""
	}
}

func getLexemes(runes []rune) (result [][]rune, err error) {
	startPosition := 0
	for curPosition := 0; curPosition < len(runes); curPosition++ {
		if unicode.IsDigit(runes[startPosition]) {
			return nil, ErrInvalidString
		}
		if isEndOfLex(runes, curPosition) {
			result = append(result, runes[startPosition:curPosition+1])
			startPosition = curPosition + 1
		}
	}
	return
}

func isEndOfLex(runes []rune, index int) bool {
	if index+1 >= len(runes) {
		return true
	}
	isSlashedSlash := false
	isCounter := false
	nextSym := runes[index+1]
	curSym := runes[index]
	if index > 0 {
		prevSym := runes[index-1]
		isSlashedSlash = !isSlash(prevSym) && isSlash(curSym) && isSlash(nextSym)
		isCounter = !isSlash(prevSym) && unicode.IsDigit(curSym)
	}
	if isCounter {
		return true
	}

	if !unicode.IsDigit(nextSym) && !isSlashedSlash {
		return true
	}

	return false
}
