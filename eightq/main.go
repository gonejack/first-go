package main

import (
	"fmt"
	"strings"
)

const N = 8

var M = [N][N]bool{}
var count = 0

func main() {
	findQ(0)

	fmt.Println(count)
}

func findQ(row int) {
	if row >= N {
		count++
		printQ()
		return
	}

	for col := range ran(N) {
		if canPutQ(row, col) {
			M[row][col] = true
			findQ(row + 1)
			M[row][col] = false
		}
	}
}

func canPutQ(row int, col int) bool {
	// 检查是否同列冲突
	for r := range ran(row) {
		if M[r][col] {
			return false
		}
	}
	// 检查正向斜线冲突/
	for r, c := range forwardCross(row, col) {
		if M[r][c] {
			return false
		}
	}
	// 检查反向斜线冲突\
	for r, c := range backwardCross(row, col) {
		if M[r][c] {
			return false
		}
	}
	return true
}

func printQ() {
	for r := range ran(N) {
		for c := range ran(N) {
			if M[r][c] {
				fmt.Println(strings.Repeat("✧ ", c) + "✦" + strings.Repeat(" ✧", N-1-c))
			}
		}
	}
	fmt.Println(strings.Repeat(" ", N))
}

func ran(n int) []struct{} {
	return make([]struct{}, n)
}

func forwardCross(row, col int) (m map[int]int) {
	m = make(map[int]int)
	for {
		row = row - 1
		col = col + 1
		if row < 0 || col >= N {
			break
		}
		m[row] = col
	}
	return
}

func backwardCross(row, col int) (m map[int]int) {
	m = make(map[int]int)
	for {
		row = row - 1
		col = col - 1
		if row < 0 || col < 0 {
			break
		}
		m[row] = col
	}
	return
}
