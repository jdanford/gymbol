package gymbol

import (
	"github.com/benbjohnson/immutable"
)

type Env struct {
	values *immutable.Map[Symbol, Value]
	vm     *VM
}

type SymbolHasher struct {
	vm *VM
}

func (s SymbolHasher) Hash(sym Symbol) uint32 {
	hasher := NewHasher(s.vm)
	hasher.WriteUint32(sym.Index)
	return hasher.Sum()
}

func (s SymbolHasher) Equal(a Symbol, b Symbol) bool {
	return a == b
}

func NewEnv(vm *VM) *Env {
	hasher := SymbolHasher{vm}
	values := immutable.NewMap[Symbol, Value](hasher)
	return &Env{values, vm}
}

func (env *Env) Get(key Symbol) (val Value, ok bool) {
	return env.values.Get(key)
}

func (env *Env) Set(key Symbol, val Value) *Env {
	newValues := env.values.Set(key, val)
	return &Env{newValues, env.vm}
}
