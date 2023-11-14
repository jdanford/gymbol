package gymbol

import (
	"strings"
	"unicode"
)

const L_PAREN = '('
const R_PAREN = ')'

const S_QUOTE = '\''
const D_QUOTE = '"'
const BACKSLASH = '\\'

const DASH = '-'

const BACKTICK = '`'
const COMMA = ','
const AT = '@'

const nonSymbolChars = string(L_PAREN) + string(R_PAREN) + string(D_QUOTE) + string(S_QUOTE) + string(BACKTICK) + string(COMMA)

func IsSymbol(r rune) bool {
	return unicode.IsPrint(r) && !unicode.IsSpace(r) && !strings.ContainsRune(nonSymbolChars, r)
}

func IsQuote(r rune) bool {
	return r == S_QUOTE || r == BACKTICK || r == COMMA
}
