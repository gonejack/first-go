package main

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	// Create .gz file to write to
	outputFile, err := os.Create("test.txt.gz")
	if err != nil {
		log.Fatal(err)
	}

	gzipw := gzip.NewWriter(outputFile)
	defer gzipw.Close()

	tarw := tar.NewWriter(gzipw)
	defer tarw.Close()

	tarw.WriteHeader(&tar.Header{
		Name: "abc.text",
		Size: 7,
	})

	io.Copy(tarw, strings.NewReader("abcdef."))

	log.Println("Compressed data written to file.")
}
