package main

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/go-getter"
)

func main() {
	err := getter.GetFile("out.jpg", "https://img.solidot.org//0/446/liiLIZF8Uh6yM.jpg")

	spew.Dump(err)
}
