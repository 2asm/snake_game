//go:build js && wasm

package snake

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"syscall/js"
	"time"
)

type game struct {
	snake         *snake
	food          food
	score         int
	height, width int
	isOver        bool
	phaseThrough  bool
}

func NewGame(height, width int) *game {
	new_game := &game{
		snake:  initialSnake(),
		score:  0,
		height: height,
		width:  width,
	}
	grid[0][0].Set("style", "background:grey;") // fast inplace rendering
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
	grid[g.food.pos.x][g.food.pos.y].Set("innerText", "") // fast inplace rendering
	g.food = newFood(c)
	grid[g.food.pos.x][g.food.pos.y].Set("innerText", string(g.food.emoji)) // fast inplace rendering
}

func (g *game) moveSnake() error {
	if err := g.snake.move(); err != nil {
		return err
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
	grid[h.x][h.y].Set("style", "background:grey;") // fast inplace rendering
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
		case s := <-restartChan:
			g.cleanUp()
			g = NewGame(g.height, g.width)
			if s == "p" || s == "P" {
				g.phaseThrough = true
			}
			g.setMode()
			g.placeFood()
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
	result       js.Value
	r, p         js.Value
	grid         = make([][]js.Value, 0)
	moveChan     = make(chan direction)
	restartChan  = make(chan string)
	idToCoordMap = map[string]coord{}
)

func coordToId(c coord) string {
	return fmt.Sprintf("%v-%v", c.x, c.y)
}

func idToCoord(id string) coord {
	if c, ok := idToCoordMap[id]; ok {
		return c
	}
	parts := strings.Split(id, "-")
	x, err := strconv.Atoi(parts[0])
	if err != nil {
		panic("atoi")
	}
	y, err := strconv.Atoi(parts[0])
	if err != nil {
		panic("atoi")
	}
	out := coord{x, y}
	idToCoordMap[id] = out
	return out
}

func init() {
	for i := range 15 { // fix
		gi := []js.Value{}
		for j := range 20 {
			e := js.Global().Get("document").Call("getElementById", coordToId(coord{i, j}))
			gi = append(gi, e)
		}
		grid = append(grid, gi)
	}

	result = js.Global().Get("document").Call("getElementById", "result")
	r = js.Global().Get("document").Call("getElementById", "r")
	p = js.Global().Get("document").Call("getElementById", "p")

	js.Global().Get("document").
		Call("addEventListener", "keydown", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			key := args[0].Get("key").String()
			switch key {
			case "ArrowUp":
				moveChan <- _UP
			case "ArrowDown":
				moveChan <- _DOWN
			case "ArrowLeft":
				moveChan <- _LEFT
			case "ArrowRight":
				moveChan <- _RIGHT
			case "r", "R", "p", "P":
				restartChan <- key
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

func (g *game) cleanUp() {
	for _, c := range g.snake.body {
		if !g.outOfArena(c) {
			grid[c.x][c.y].Set("style", "background:#ddd;")
		}
	}
	grid[g.food.pos.x][g.food.pos.y].Set("innerText", "")
}
