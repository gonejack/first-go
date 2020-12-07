package main

import (
	"fmt"
	"strings"
)

type direction int

const (
	UP direction = iota
	RIGHT
	DOWN
	LEFT
)

func (d direction) String() string {
	switch d {
	case UP:
		return "上"
	case RIGHT:
		return "右"
	case DOWN:
		return "下"
	case LEFT:
		return "左"
	}
	return "x"
}

var want = nmap{{1, 2, 3}, {4, 5, 6}, {7, 8, 0}}
var directions = [...]direction{UP, RIGHT, DOWN, LEFT}

const N = 3

type nmap [N][N]int
type nmapSet map[nmap]struct{}

type eight struct {
	nmap
	moves      []direction
	researched nmapSet
	researches researchList
}

func (e *eight) resolve() bool {
	if e.nmap == want {
		return true
	} else {
		e.addResearched(e.nmap)
		e.addResearch(newResearch(e.nmap, nil, nil))
		return e.doResearch()
	}
}
func (e *eight) doResearch() bool {
	var r *research

	for len(e.researches) > 0 {
		r, e.researches = e.researches[0], e.researches[1:]

		if r.nmap == want {
			for r != nil && r.direction != nil {
				e.moves = append(e.moves, *r.direction)
				r = r.prev
			}
			printNmap(e.nmap)
			printMoves(e.moves)
			return true
		} else {
			for i := range directions {
				var dir = directions[i]
				if r.canMove(dir) {
					r.moveForward(dir)
					if !e.isResearched(r.nmap) {
						e.addResearched(r.nmap)
						e.addResearch(newResearch(r.nmap, &dir, r))
					}
					r.moveBack(dir)
				}
			}
		}
	}

	return false
}
func (e *eight) addResearch(search *research) {
	e.researches = append(e.researches, search)
}
func (e *eight) addResearched(nmap nmap) {
	e.researched[nmap] = struct{}{}
}
func (e *eight) isResearched(nmap nmap) (exist bool) {
	_, exist = e.researched[nmap]
	return
}

func printNmap(m nmap) {
	for _, r := range m {
		fmt.Println(r)
	}
	fmt.Println(strings.Repeat("-", 3))
}
func printMoves(moves []direction) {
	var chinese []string
	var length = len(moves)
	for i := length - 1; i >= 0; i-- {
		chinese = append(chinese, moves[i].String())
	}
	fmt.Println(strings.Join(chinese, " "))
}
func new8(m [3][3]int) *eight {
	return &eight{
		nmap:       m,
		researched: make(nmapSet),
		researches: make(researchList, 0),
	}
}

func main() {
	m := [N][N]int{
		{3, 4, 1},
		{5, 6, 0},
		{8, 2, 7},
	}
	new8(m).resolve()
}
