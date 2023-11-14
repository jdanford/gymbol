package gymbol

import (
	"fmt"
)

type Ref struct {
	Index uint32
}

func (ref Ref) GoString() string {
	return fmt.Sprintf("Ref{%d}", ref.Index)
}

func AsRef(value Value) (Ref, error) {
	ref, ok := value.(Ref)
	if !ok {
		return Ref{}, fmt.Errorf("expected ref, got `%v`", value)
	}

	return ref, nil
}
