package gymbol

import (
	"fmt"
)

func (vm *VM) Render(val Value) string {
	switch val := val.(type) {
	case Number:
		return renderNumber(val)
	case Symbol:
		return vm.Symbols.Resolve(val)
	default:
		return fmt.Sprintf("`%v`", val)
	}
}

func renderNumber(num Number) string {
	f := float64(num)
	i := int32(f)

	if f == float64(i) {
		return fmt.Sprintf("%d", i)
	}

	return fmt.Sprintf("%f", f)
}
