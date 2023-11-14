package gymbol

import (
	"encoding/binary"
	"hash"
	"math"

	"github.com/dchest/siphash"
)

var SIPHASH_KEY = []byte{
	0x5a,
	0x52,
	0xff,
	0x2e,
	0x4e,
	0x09,
	0x8e,
	0x62,
	0x0b,
	0x61,
	0x46,
	0x6e,
	0x7a,
	0xf9,
	0x64,
	0x93,
}

type Hasher struct {
	hash.Hash64
	vm *VM
}

func NewHasher(vm *VM) *Hasher {
	sip := siphash.New(SIPHASH_KEY)
	return &Hasher{sip, vm}
}

func (hasher *Hasher) Sum() uint32 {
	sum64 := hasher.Hash64.Sum64()
	return uint32(sum64)
}

func (hasher *Hasher) WriteBytes(b []byte) {
	hasher.Write(b)
}

func (hasher *Hasher) WriteUint32(n uint32) {
	var b [4]byte
	binary.BigEndian.PutUint32(b[:], n)
	hasher.Write(b[:])
}

func (hasher *Hasher) WriteUint64(n uint64) {
	var b [8]byte
	binary.BigEndian.PutUint64(b[:], n)
	hasher.Write(b[:])
}

func (hasher *Hasher) WriteFloat64(f float64) {
	n := math.Float64bits(f)
	hasher.WriteUint64(n)
}

func (hasher *Hasher) WriteType(sym Symbol) {
	hasher.WriteUint32(sym.Index)
}

func (hasher *Hasher) WriteNumber(num Number) {
	f := float64(num)
	hasher.WriteType(NUMBER)
	hasher.WriteFloat64(f)
}

func (hasher *Hasher) WriteSymbol(sym Symbol) {
	hasher.WriteType(SYMBOL)
	hasher.WriteUint32(sym.Index)
}

func (hasher *Hasher) WriteString(str String) {
	bytes := []byte(str)
	hasher.WriteBytes(bytes)
}

func (hasher *Hasher) WriteRef(ref Ref) {
	type_, values := hasher.vm.Heap.Load(ref)
	hasher.WriteType(type_)

	for i := range values {
		value := values[i]
		hasher.WriteValue(value)
	}
}

func (hasher *Hasher) WriteQuoted(quoted Quoted) {
	hasher.WriteUint32(quoted.Type.Index)
	hasher.WriteValue(quoted.Value)
}

func (hasher *Hasher) WriteValue(val Value) {
	switch val := val.(type) {
	case Number:
		hasher.WriteNumber(val)
	case Symbol:
		hasher.WriteSymbol(val)
	case String:
		hasher.WriteString(val)
	case Ref:
		hasher.WriteRef(val)
	case Quoted:
		hasher.WriteQuoted(val)
	default:
		panic(val)
	}
}
