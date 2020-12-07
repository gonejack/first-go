package main

var changes = [...][2]int{{0, -1}, {1, 0}, {0, 1}, {-1, 0}}

type (
	research struct {
		nmap      nmap
		direction *direction
		prev      *research
		spaceX    int
		spaceY    int
	}

	researchList []*research
)

func (r *research) moveForward(dir direction) {
	targetX, targetY := r.spaceX+changes[dir][0], r.spaceY+changes[dir][1]
	r.nmap[targetY][targetX], r.nmap[r.spaceY][r.spaceX] = 0, r.nmap[targetY][targetX]
	r.spaceX, r.spaceY = targetX, targetY
}
func (r *research) moveBack(dir direction) {
	r.moveForward((dir + 2) % 4)
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
	for y := range r.nmap {
		for x := range r.nmap[y] {
			if r.nmap[y][x] == 0 {
				r.spaceX, r.spaceY = x, y
				return
			}
		}
	}
}

func newResearch(nmap nmap, dir *direction, prevResearch *research) (r *research) {
	r = &research{
		nmap:      nmap,
		direction: dir,
		prev:      prevResearch,
	}
	r.setup()

	return
}
