package gymbol

import "io"

type VM struct {
	Symbols *SymbolTable
	Heap    *Heap
}

func NewVM() *VM {
	table := NewSymbolTable()
	heap := NewHeap()
	return &VM{Symbols: table, Heap: heap}
}

func (vm *VM) Read(input io.Reader) (Value, error) {
	return Read(vm, input)
}
