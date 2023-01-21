package tool

import (
	"fmt"
)

// Print print value
func Print(value interface{}) {
	fmt.Printf("%+v\n", value)
}
