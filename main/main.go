package main

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/mitchellh/mapstructure"
)

func main() {
	type abc struct {
		Abc string
		Def int
	}

	var output abc

	var input = map[string]interface{}{
		"abc": "text",
		"def": "123",
	}
	spew.Dump(mapstructure.WeakDecode(input, &output))

	spew.Dump(output)
}
