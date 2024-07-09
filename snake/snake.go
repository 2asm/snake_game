//go:build js && wasm

package snake

import "fmt"

type snake struct {
	body      []coord
	direction direction
	len       int
}

func newSnake(body []coord, direction direction) *snake {
	return &snake{
		body:      body,
		direction: direction,
		len:       len(body),
	}
}

func initialSnake() *snake {
	return &snake{
		body:      []coord{{0, 0}},
		direction: _RIGHT,
		len:       1,
	}
}

func (s *snake) head() coord {
	return s.body[len(s.body)-1]
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

func (s *snake) move() error {
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
		return fmt.Errorf("died")
	}
	if s.len > len(s.body) {
		s.body = append(s.body, h)
	} else {
		tmp := s.body[0]
		grid[tmp.x][tmp.y].Set("style", "background:#ddd;") // fast inplace rendering
		s.body = append(s.body[1:], h)
	}
	return nil
}
