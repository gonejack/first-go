package main

import "fmt"

func main() {
	f := Fibonacci()
	for i := 0; i < 10; i++ {
		fmt.Println(f())
	}
}

func Fibonacci() func() int {
	a, b := 0, 1
	return func() (i int) {
		i, a, b = a, b, a+b
		return
	}
}
