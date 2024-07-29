//go:build js && wasm

package snake

import "fmt"

type snake struct {
	body      []coord
	direction direction
	len       int
	color     string
}

func newSnake(body []coord, direction direction, color string) *snake {
	return &snake{
		body:      body,
		direction: direction,
		len:       len(body),
		color:     color,
	}
}

func initialSnake() *snake {
	return &snake{
		body:      []coord{{0, 0}, {0, 1}},
		direction: _RIGHT,
		color:     "grey",
		len:       2,
	}
}

func (s *snake) head() coord {
	return s.body[len(s.body)-1]
}

func (s *snake) neck() coord {
	return s.body[len(s.body)-2]
}

func (s *snake) UpdateHead(newHead coord) {
	s.body[len(s.body)-1].x = newHead.x
	s.body[len(s.body)-1].y = newHead.y
}

func (s *snake) onBody(c coord) bool {
	for _, p := range s.body {
		if p == c {
			return true
		}
	}
	return false
}

func (s *snake) changeDirection(d direction) {
	if d != s.direction.opposite() {
		s.direction = d
	}
}

func (s *snake) move() (err error, retCell *coord) {
	h := s.head()
	switch s.direction {
	case _LEFT:
		h.y--
	case _RIGHT:
		h.y++
	case _UP:
		h.x--
	case _DOWN:
		h.x++
	}
	if s.onBody(h) {
		err = fmt.Errorf("died")
	}
	if s.len > len(s.body) {
		s.body = append(s.body, h)
	} else {
		tmp := s.body[0]
		retCell = &tmp
		s.body = append(s.body[1:], h)
	}
	return
}
