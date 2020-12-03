package main

import (
	"fmt"
	"math"
)

func Sqrt(x float64) float64 {
	low := 0.0
	up := x
	prec := 1e-10

	for math.Abs(up-low) > prec {
		m := (low + up) / 2
		sqrt := m * m
		if sqrt == x {
			return m
		} else if sqrt > x {
			up = m
		} else {
			low = m
		}
	}

	return up
}

func main() {
	fmt.Println(Sqrt(2))
	fmt.Println(math.Sqrt(2))
}
