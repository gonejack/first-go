package main

var changes = [...][2]int{
	{0, -1},
	{1, 0},
	{0, 1},
	{-1, 0},
}

type research struct {
	t9

	moveDirection  *direction
	parentResearch *research

	spaceX int
	spaceY int
}
type researchList []*research

func (r *research) move(dir direction) {
	targetX, targetY := r.spaceX+changes[dir][0], r.spaceY+changes[dir][1]
	r.t9[targetY][targetX], r.t9[r.spaceY][r.spaceX] = 0, r.t9[targetY][targetX]
	r.spaceX, r.spaceY = targetX, targetY
}
func (r *research) revert(dir direction) {
	r.move((dir + 2) % 4)
}
func (r *research) canMove(d direction) (can bool) {
	switch d {
	case UP:
		return r.spaceY > 0
	case RIGHT:
		return r.spaceX < 2
	case DOWN:
		return r.spaceY < 2
	case LEFT:
		return r.spaceX > 0
	}
	return
}
func (r *research) setup() {
	for y := range r.t9 {
		for x := range r.t9[y] {
			if r.t9[y][x] == 0 {
				r.spaceX, r.spaceY = x, y
				return
			}
		}
	}
}

func newResearch(t9 t9, direction *direction, parentResearch *research) (r *research) {
	r = &research{
		t9:             t9,
		moveDirection:  direction,
		parentResearch: parentResearch,
	}
	r.setup()

	return
}
