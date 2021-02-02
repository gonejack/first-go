package main

import "github.com/bmaupin/go-epub"

func main() {
	// Create a new EPUB
	e := epub.NewEpub("My title")

	// Set the author
	e.SetAuthor("Hingle McCringleberry")

	// Add a section
	section1Body := `<h1>Section 1</h1>
	<p>This is a paragraph.</p>`
	e.AddSection(section1Body, "Section 1", "", "")

	// Write the EPUB
	err := e.Write("My EPUB.epub")
	if err != nil {
		// handle error
	}
}
