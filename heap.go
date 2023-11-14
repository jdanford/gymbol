package gymbol

import (
	"fmt"
)

type HeapRecord struct {
	Type   Symbol
	Size   int
	Offset int
}

type Heap struct {
	data    []Value
	records []HeapRecord
}

func NewHeap() *Heap {
	data := []Value{}
	records := []HeapRecord{}
	return &Heap{data, records}
}

func (heap *Heap) Alloc(type_ Symbol, values []Value) (Ref, error) {
	size := len(values)
	offset := len(heap.data)
	index := uint32(len(heap.records))

	record := HeapRecord{type_, size, offset}
	heap.data = append(heap.data, values...)
	heap.records = append(heap.records, record)

	return Ref{Index: index}, nil
}

func (heap *Heap) Load(ref Ref) (Symbol, []Value) {
	record := heap.records[ref.Index]
	start := record.Offset
	end := start + record.Size
	data := heap.data[start:end]
	return record.Type, data
}

func (heap *Heap) LoadChecked(ref Ref, expectedType Symbol, expectedSize int) ([]Value, error) {
	refType, data := heap.Load(ref)
	if refType != expectedType {
		return nil, fmt.Errorf("expected type %s, got %s", expectedType, refType)
	}

	size := len(data)
	if size != expectedSize {
		return nil, fmt.Errorf("expected size %d, got %d", expectedSize, size)
	}

	return data, nil
}

func (heap *Heap) AllocCons(head Value, tail Value) (Ref, error) {
	data := []Value{head, tail}
	return heap.Alloc(CONS, data)
}

func (heap *Heap) LoadConsRef(ref Ref) (Value, Value, error) {
	data, err := heap.LoadChecked(ref, CONS, 2)
	head, tail := data[0], data[1]
	return head, tail, err
}

func (heap *Heap) AllocList(values []Value) (Value, error) {
	var list Value = NIL
	for i := len(values) - 1; i >= 0; i-- {
		head := values[i]
		newList, err := heap.AllocCons(head, list)
		if err != nil {
			return nil, err
		}

		list = newList
	}

	return list, nil
}

func (heap *Heap) LoadList(list Value) ([]Value, error) {
	values := []Value{}

	for list != NIL {
		ref, err := AsRef(list)
		if err != nil {
			return nil, err
		}

		head, tail, err := heap.LoadConsRef(ref)
		if err != nil {
			return nil, err
		}

		values = append(values, head)
		list = tail
	}

	return values, nil
}
