package snake

type direction int

const (
	_LEFT direction = iota + 1
	_UP
	_RIGHT
	_DOWN
)

func (d direction) opposite() direction {
	switch d {
	case _LEFT:
		return _RIGHT
	case _RIGHT:
		return _LEFT
	case _UP:
		return _DOWN
	case _DOWN:
		return _UP
	default:
		return 0
	}
}
