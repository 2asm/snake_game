//go:build js && wasm

package snake

import (
	"fmt"
	"math/rand"
	"syscall/js"
	"time"
)

type game struct {
	snake         *snake
	food          food
	score         int
	height, width int
	scale         int // pixel size
	isOver        bool
	phaseThrough  bool
}

func NewGame(height, width, scale int) *game {
	new_game := &game{
		snake:  initialSnake(),
		score:  0,
		height: height,
		width:  width,
		scale:  scale,
	}
	new_game.init()
	new_game.fillCell(0, 0) // fast inplace rendering
	new_game.placeFood()
	return new_game
}

func (g *game) outOfArena(c coord) bool {
	return c.x < 0 || c.x >= g.height || c.y < 0 || c.y >= g.width
}

func (g *game) hasFood(c coord) bool {
	return g.food.pos == c
}

func (g *game) placeFood() {
	c := coord{-1, -1}
	for {
		c.x = rand.Intn(g.height)
		c.y = rand.Intn(g.width)
		if !g.snake.onBody(c) {
			break
		}
	}
	g.food = newFood(c)
	g.fillTextCell(c.x, c.y, g.food.emoji) // fast inplace rendering
}

func (g *game) moveSnake() error {
	err, c := g.snake.move()
	if err != nil {
		return err
	}
	if c != nil {
		g.clearCell(c.x, c.y)
	}
	h := g.snake.head()
	if g.outOfArena(h) {
		if !g.phaseThrough {
			return fmt.Errorf("died")
		}
		h.x = (h.x + g.height) % g.height
		h.y = (h.y + g.width) % g.width
		g.snake.UpdateHead(h)
	}
	g.fillCell(h.x, h.y) // fast inplace rendering
	if g.hasFood(h) {
		g.score += g.food.points
		g.renderResult()
		g.snake.len += 1
		g.placeFood()
	}
	return nil
}

func (g *game) moveInterval() time.Duration {
	ms := max(50, 100-g.snake.len/10) // milliseconds
	return time.Duration(time.Millisecond * time.Duration(ms))
}

func (g *game) Start() {
	g.setMode()
	for {
		select {
		case d := <-moveChan:
			g.snake.changeDirection(d)
		case keyCode := <-restartChan:
			g.cleanUpSnake()
            g.clearCell(g.food.pos.x, g.food.pos.y)
			// g.clearAll()
			g = NewGame(g.height, g.width, g.scale)
			if keyCode == 80 { // p
				g.phaseThrough = true
			}
			g.setMode()
			g.renderResult()
		default:
			if !g.isOver {
				if err := g.moveSnake(); err != nil {
					g.isOver = true
				}
			}
			time.Sleep(g.moveInterval())
		}
	}
}

var (
	_INIT       bool
	result      js.Value
	r, p        js.Value
	gameCanvas  js.Value
	moveChan    = make(chan direction)
	restartChan = make(chan int)
)

func (g *game) init() {
	if _INIT {
		return
	}
	_INIT = true
	// todo: change scale accoding to resolution
	c := js.Global().Get("document").Call("getElementById", "gameCanvas")
	c.Set("height", g.scale*g.height)
	c.Set("width", g.scale*g.width)
	gameCanvas = c.Call("getContext", "2d")

	result = js.Global().Get("document").Call("getElementById", "result")
	r = js.Global().Get("document").Call("getElementById", "r")
	p = js.Global().Get("document").Call("getElementById", "p")

	js.Global().Get("document").
		Call("addEventListener", "keydown", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			keyCode := args[0].Get("keyCode").Int()
			switch keyCode {
			case 37, 38, 39, 40:
				moveChan <- direction(keyCode - 37 + 1)
			case 82, 80: // r, p
				restartChan <- keyCode
			}
			return nil
		}))
}

func (g *game) setMode() {
	if g.phaseThrough {
		r.Call("removeAttribute", "style")
		p.Call("setAttribute", "style", "color:orange;")
	} else {
		p.Call("removeAttribute", "style")
		r.Call("setAttribute", "style", "color:orange;")
	}
}

func (g *game) renderResult() {
	result.Set("innerText", fmt.Sprintf("Score: %v", g.score))
}

func (g *game) cleanUpSnake() {
	for _, c := range g.snake.body {
		if !g.outOfArena(c) {
			g.clearCell(c.x, c.y)
		}
	}
}

func (g *game) clearAll() {
	for x := range g.height {
		for y := range g.width {
			gameCanvas.Call("clearRect", y*g.scale, x*g.scale, g.scale, g.scale)
		}
	}
}

func (g *game) fillCell(x, y int) {
	gameCanvas.Set("fillStyle", g.snake.color)
	gameCanvas.Call("fillRect", y*g.scale, x*g.scale, g.scale, g.scale)
	gameCanvas.Call("fill")
}

func (g *game) fillTextCell(x, y int, ch rune) {
	// gameCanvas.Set("fillStyle", "green")
	gameCanvas.Set("font", fmt.Sprintf("%vpx arial", g.scale*9/10))
	gameCanvas.Set("textAlign", "left")
	gameCanvas.Set("textBaseline", "top")
	gameCanvas.Call("fillText", string(ch), y*g.scale, x*g.scale+g.scale/10)
	gameCanvas.Call("fill")
}

func (g *game) clearCell(x, y int) {
	gameCanvas.Call("clearRect", y*g.scale, x*g.scale, g.scale, g.scale)
	gameCanvas.Call("fill")
}
