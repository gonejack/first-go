package main

import (
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		return
	}

	fd, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	(&sorter{source: fd}).sort()
}
