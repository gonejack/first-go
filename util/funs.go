package util

import (
	"log"
	"os"
)

func IdxRange(n int) []int {
	return make([]int, n)
}

func EmptyPanic(str string) string {
	if str == "" {
		panic("Illegal empty config")
	}

	return str
}
func EmptyFB(str string, fb string) string {
	if str == "" {
		return fb
	}

	return str
}
func LessFB(i int, lessThan int, fb int) int {
	if i < lessThan {
		return fb
	}

	return i
}

func CreativeDir(dir string) {
	if _, e := os.Stat(dir); os.IsNotExist(e) {
		if e := os.MkdirAll(dir, 0755); e == nil {
			log.Printf("Created directory[%s]", dir)
		} else {
			log.Fatalf("Error creating directory[%s]: %s", dir, e)
		}
	}
}
