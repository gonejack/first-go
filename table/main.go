package main

import (
	"fmt"
	"github.com/scylladb/termtables"
)

func main() {
	table := termtables.CreateTable()

	table.AddHeaders("Name", "Age")
	table.AddRow("John", "", "44", "55")
	table.AddRow("Sam", 18)
	table.AddRow("Julie", 20.14)
	table.UTF8Box()
	fmt.Println(table.Render())
}
