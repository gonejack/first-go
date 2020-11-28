package main

import (
	"fmt"
	"github.com/Jeffail/gabs/v2"
)

var json = []byte(`{
    "info": {
      "name": {
        "first": "lee",
        "last": "darjun"
      },
      "age": 18,
      "hobbies": [
        "game",
        "programming"
      ]
    }
    }`,
)

func main() {
	obj, _ := gabs.ParseJSON(json)

	fmt.Println("first name: ", obj.Search("info", "name", "first").Data().(string))
	fmt.Println("second name: ", obj.Path("info.name.last").Data().(string))

	container, _ := obj.JSONPointer("/info/age")
	fmt.Println("age: ", container.Data().(float64))
	fmt.Println("one hobby: ", obj.Path("info.hobbies.1").Data().(string))
}
