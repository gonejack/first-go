package main

import (
	"fmt"
	"strings"
)

type direction int

func (d direction) String() string {
	return directionTexts[d]
}

const (
	UP direction = iota
	RIGHT
	DOWN
	LEFT
)

var (
	directions     = [...]direction{UP, RIGHT, DOWN, LEFT}
	directionTexts = [...]string{"上", "右", "下", "左"}
)

type (
	t9    [3][3]int
	t9set map[t9]struct{}
)

var want = t9{
	{1, 2, 3},
	{4, 5, 6},
	{7, 8, 0},
}

type eight struct {
	t9
	moves      []direction
	researched t9set
	researches researchList
}

func (e *eight) resolve() bool {
	if e.t9 == want {
		return true
	} else {
		e.addResearched(e.t9)
		e.enqueueResearch(newResearch(e.t9, nil, nil))
		return e.doResearch()
	}
}
func (e *eight) doResearch() bool {
	var rs *research

	for len(e.researches) > 0 {
		rs, e.researches = e.researches[0], e.researches[1:]

		if rs.t9 == want {
			for rs != nil && rs.moveDirection != nil {
				e.moves = append(e.moves, *rs.moveDirection)
				rs = rs.parentResearch
			}
			printNmap(e.t9)
			printMoves(e.moves)
			return true
		} else {
			for i := range directions {
				var dir = directions[i]
				if rs.canMove(dir) {
					rs.move(dir)
					if !e.isResearched(rs.t9) {
						e.addResearched(rs.t9)
						e.enqueueResearch(newResearch(rs.t9, &dir, rs))
					}
					rs.revert(dir)
				}
			}
		}
	}

	return false
}
func (e *eight) enqueueResearch(search *research) {
	e.researches = append(e.researches, search)
}
func (e *eight) addResearched(t9 t9) {
	e.researched[t9] = struct{}{}
}
func (e *eight) isResearched(t9 t9) (exist bool) {
	_, exist = e.researched[t9]
	return
}

func newResolver(data t9) (e *eight) {
	e = &eight{
		t9: data,

		researched: make(t9set),
		researches: make(researchList, 0),
	}
	return
}

func printNmap(m t9) {
	for _, r := range m {
		fmt.Println(r)
	}
	fmt.Println(strings.Repeat("-", 3))
}
func printMoves(moves []direction) {
	var chars []string
	var length = len(moves)
	for i := length - 1; i >= 0; i-- {
		chars = append(chars, moves[i].String())
	}
	fmt.Println(strings.Join(chars, " "))
}

func main() {
	data := t9{
		{3, 4, 1},
		{5, 6, 0},
		{8, 2, 7},
	}

	newResolver(data).resolve()
}
