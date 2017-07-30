package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

func update(screen *ebiten.Image) error {
	ebitenutil.DebugPrint(screen, "Hello")
	return nil
}

func main() {
	if err := ebiten.Run(update, 320, 240, 2, "Hello World!"); err != nil {
		panic(err)
	}
	fmt.Println("bye")
}
