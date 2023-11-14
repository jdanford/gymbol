package gymbol

import (
	"fmt"
)

type Symbol struct {
	Index uint32
}

func (sym Symbol) String() string {
	return fmt.Sprintf("%d", sym.Index)
}

func (sym Symbol) GoString() string {
	return fmt.Sprintf("Symbol{%d}", sym.Index)
}

func AsSymbol(value Value) (Symbol, error) {
	return Symbol{}, fmt.Errorf("expected symbol, got `%v`", value)
}

var NIL = Symbol{Index: 0}
var CONS = Symbol{Index: 1}
var TRUE = Symbol{Index: 2}
var FALSE = Symbol{Index: 3}
var SYMBOL = Symbol{Index: 4}
var NUMBER = Symbol{Index: 5}
var STRING = Symbol{Index: 6}
var REF = Symbol{Index: 7}
var FN = Symbol{Index: 8}
var NATIVE_FN = Symbol{Index: 9}
var QUOTE = Symbol{Index: 10}
var QUASIQUOTE = Symbol{Index: 11}
var UNQUOTE = Symbol{Index: 12}
var UNQUOTE_SPLICING = Symbol{Index: 13}

var CORE_SYMBOL_NAMES = []string{
	"nil",
	"cons",
	"true",
	"false",
	"symbol",
	"number",
	"string",
	"ref",
	"fn",
	"native-fn",
	"quote",
	"quasiquote",
	"unquote",
	"unquote-splicing",
}
