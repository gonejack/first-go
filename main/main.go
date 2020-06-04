package main

import (
	"github.com/davecgh/go-spew/spew"
	"net/http"
)

type Key struct {
	Date int
	Hour int
	//...
}
type Log struct {
	Key
	Bids   int
	Offers int
}

type Person struct {
	Name string
	Age  int
}

func main() {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("abc"))
		return
	}
	http.HandleFunc("/abc", fn)
	err := http.ListenAndServe(":1234", nil)

	spew.Dump(err)

	http.DefaultClient.Get("abc")
}
