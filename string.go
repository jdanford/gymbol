package gymbol

import (
	"fmt"
)

type String string

func (str String) GoString() string {
	return fmt.Sprintf("String(%v)", str)
}

func AsString(value Value) (String, error) {
	str, ok := value.(String)
	if !ok {
		return "", fmt.Errorf("expected string, got `%v`", value)
	}

	return str, nil
}
