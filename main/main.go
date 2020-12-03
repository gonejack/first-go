package main

import (
	"fmt"
	"strconv"
)

func main() {

	fmt.Printf("%8s\n", strconv.FormatInt(0x20, 2))
	fmt.Printf("%8s\n", strconv.FormatInt(int64('A'), 2))
	fmt.Printf("%8s\n", strconv.FormatInt(int64('a'), 2))
	fmt.Printf("%8s\n", strconv.FormatInt(int64('Z'), 2))
	fmt.Printf("%8s\n", strconv.FormatInt(int64('z'), 2))
}
