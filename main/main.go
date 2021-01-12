package main

import (
	"github.com/davecgh/go-spew/spew"
	"net/url"
)

func main() {
	u := "http%3A%2F%2Flocalhost%3A9001%2Fios"
	spew.Dump(u)

	u, _ = url.QueryUnescape(u)
	spew.Dump(u)

	u, _ = url.QueryUnescape(u)
	spew.Dump(u)
}
