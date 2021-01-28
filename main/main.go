package main

import (
	"github.com/davecgh/go-spew/spew"
	"net/url"
	"path/filepath"
)

func main() {
	s := "https://cdn.sspai.com/article/d3d0016c-6dbc-72c5-9ed2-bfa7fbb5b323.png?imageMogr2/auto-orient/quality/95/thumbnail/!690x690r/gravity/Center/crop/690x690/interlace/1"

	spew.Dump(filepath.Base(s))

	u, _ := url.Parse(s)
	spew.Dump(filepath.Base(u.Path))
}
