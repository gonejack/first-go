package main

import (
	"github.com/SlyMarbo/rss"
)

func main() {
	feed, err := rss.Fetch("http://www.ruanyifeng.com/blog/atom.xml")
	if err != nil {
		// handle error.
	}

	// ... Some time later ...

	err = feed.Update()
	if err != nil {
		// handle error.
	}
}
