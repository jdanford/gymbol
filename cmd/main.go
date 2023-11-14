package main

import (
	"fmt"
	"strings"

	"gymbol"
)

func main() {
	vm := gymbol.NewVM()
	input := strings.NewReader("(1 2 3)")
	value, err := vm.Read(input)
	if err != nil {
		fmt.Println(err)
		return
	}

	str := vm.Render(value)
	fmt.Println(str)
}
