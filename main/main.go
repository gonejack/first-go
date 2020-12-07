package main

import (
	"fmt"
	"strconv"
)

func main() {
	fmt.Printf("%8s\n", strconv.FormatInt(1<<2>>2, 2))
	fmt.Printf("%8s\n", strconv.FormatInt(1<<4>>2, 2))

}
