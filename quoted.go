package gymbol

import (
	"fmt"
)

type Quoted struct {
	Type  Symbol
	Value Value
}

func NewQuoted(quoteType Symbol, value Value) Quoted {
	return Quoted{Type: quoteType, Value: value}
}

func AsQuoted(value Value) (Quoted, error) {
	quoted, ok := value.(Quoted)
	if !ok {
		return Quoted{}, fmt.Errorf("expected quoted value, got `%v`", value)
	}

	return quoted, nil
}

func (quoted Quoted) GoString() string {
	return fmt.Sprintf("Quoted{%v, %v}", quoted.Type, quoted.Value)
}

var quoteTypes = map[rune]Symbol{
	S_QUOTE:  QUOTE,
	BACKTICK: QUASIQUOTE,
	COMMA:    UNQUOTE,
}

var quoteNames = map[Symbol]string{
	QUOTE:            string(S_QUOTE),
	QUASIQUOTE:       string(BACKTICK),
	UNQUOTE:          string(COMMA),
	UNQUOTE_SPLICING: string(COMMA) + string(AT),
}
