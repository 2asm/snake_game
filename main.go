//go:build js && wasm

package main

import "github.com/2asm/snake_game/snake"

func main() {
	snake.NewGame(20, 30, 30).Start()
}
