package gymbol

import (
	"fmt"
)

type Number float64

func (num Number) GoString() string {
	return fmt.Sprintf("Number(%f)", num)
}

func AsNumber(value Value) (Number, error) {
	num, ok := value.(Number)
	if !ok {
		return 0, fmt.Errorf("expected number, got `%v`", value)
	}

	return num, nil
}
