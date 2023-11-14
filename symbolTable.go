package gymbol

type SymbolTable struct {
	strings []string
	indices map[string]uint32
}

func NewSymbolTable() *SymbolTable {
	strings := []string{}
	indices := map[string]uint32{}
	table := &SymbolTable{strings, indices}

	for i := range CORE_SYMBOL_NAMES {
		str := CORE_SYMBOL_NAMES[i]
		table.Intern(str)
	}

	return table
}

func (table *SymbolTable) set(str string, index uint32) {
	table.strings = append(table.strings, str)
	table.indices[str] = index
}

func (table *SymbolTable) Intern(str string) Symbol {
	existingIndex, ok := table.indices[str]
	if ok {
		return Symbol{Index: existingIndex}
	}

	newIndex := uint32(len(table.strings))
	table.set(str, newIndex)
	return Symbol{Index: newIndex}
}

func (table *SymbolTable) Resolve(sym Symbol) string {
	return table.strings[sym.Index]
}
