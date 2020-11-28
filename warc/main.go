package main

import (
	"github.com/slyrz/warc"
	"os"
	"strings"
)

func main() {
	writer := warc.NewWriter(os.Stdout)
	record := warc.NewRecord()
	record.Header.Set("warc-type", "resource")
	record.Header.Set("content-type", "plain/text")
	record.Content = strings.NewReader("Hello, World!")
	if _, err := writer.WriteRecord(record); err != nil {
		panic(err)
	}
}
