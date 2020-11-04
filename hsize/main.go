package main

import (
	"bufio"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"strings"
)

const help = `Examples:
command:
  {exec} 123 45678
print:
  123 => 123B
  45678 => 44.61KB

command: 
  echo 123 | {exec}
print:
  123 => 123B
`

func main() {
	if len(os.Args) > 1 {
		for _, arg := range os.Args[1:] {
			if arg == "-h" || arg == "--help" {
				fmt.Print(strings.ReplaceAll(help, "{exec}", filepath.Base(os.Args[0])))
				return
			}
		}
		for _, arg := range os.Args[1:] {
			parse(arg)
		}
	} else {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			parse(scanner.Text())
		}
		if scanner.Err() != nil {
			fmt.Printf("error reading stdin: %s", scanner.Err())
		}
	}
}

var units = []string{"B", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"}
var scale = new(big.Rat).SetInt64(1 << 10)

func parse(text string) {
	size, ok := new(big.Rat).SetString(strings.TrimSpace(text))
	if !ok {
		fmt.Printf("can not parse \"%s\"\n", text)
		return
	}

	var unit, val string
	for _, unit = range units {
		if size.Cmp(scale) > -1 {
			size = size.Quo(size, scale)
		} else {
			break
		}
	}
	if size.IsInt() {
		val = size.RatString()
	} else {
		val = size.FloatString(2)
	}
	fmt.Printf("%s => %s%s\n", text, val, unit)
}
