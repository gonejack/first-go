package main

import (
	"github.com/davecgh/go-spew/spew"
	"path/filepath"
)

func main() {
	spew.Dump(filepath.ToSlash(`c:\windows\abc`))
}
