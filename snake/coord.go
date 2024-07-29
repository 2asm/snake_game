package snake

type coord struct {
	x, y int
}

func newCoord(x, y int) coord {
	return coord{x, y}
}
