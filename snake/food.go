package snake

import "math/rand"

type food struct {
	emoji  rune
	points int
	pos    coord
}

var foods = []rune{'🍒', '🍌', '🍇', '🍑', '🍍'}
var points = []int{1, 2, 3, 10, 100}
var prob = []int{500, 300, 189, 10, 1} // total 1000

func getRandIdxWithProbability() int {
	sm := 0
	for i := range prob {
		sm += prob[i]
	}
	rnd := rand.Intn(sm) + 1
	sm = 0
	for i := range prob {
		sm += prob[i]
		if sm >= rnd {
			return i
		}
	}
	return -1
}

func newFood(c coord) food {
	idx := getRandIdxWithProbability()
	return food{
		emoji:  foods[idx],
		points: points[idx],
		pos:    c,
	}
}
