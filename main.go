//go:build js && wasm

package main

import "github.com/2asm/snake_game/snake"

func main() {
	snake.NewGame(15, 20).Start()
}
