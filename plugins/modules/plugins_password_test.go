package modules

import (
	"fmt"
	"testing"
)

func TestGenHash(*testing.T) {
	fmt.Printf("%s", genHash("password"))
}
