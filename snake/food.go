package snake

import "math/rand"

type food struct {
	emoji  rune
	points int
	pos    coord
}

var foods = []rune{'🍪', '🍑', '🍗', '🍏'}

func newFood(c coord) food {
	points := rand.Intn(len(foods)) + 1 // random points
	return food{
		emoji:  foods[points-1],
		points: points,
		pos:    c,
	}
}
