package main

import (
	"fmt"
	"image/color"

	"github.com/Bredgren/geo"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

var square *ebiten.Image

func update(screen *ebiten.Image) error {
	screen.Fill(color.NRGBA{0xff, 0x00, 0x00, 0xff})
	ebitenutil.DebugPrint(screen, "Hello")

	if ebiten.IsKeyPressed(ebiten.KeyF) {
		ebitenutil.DebugPrint(screen, "\n\n\nF")
	}

	if square == nil {
		square, _ = ebiten.NewImage(32, 32, ebiten.FilterNearest)
	}

	square.Fill(color.White)

	squarePos := geo.VecXY(32, 32)
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(squarePos.XY())

	screen.DrawImage(square, opts)

	return nil
}

func main() {
	if err := ebiten.Run(update, 640, 480, 1, "Hello World!"); err != nil {
		panic(err)
	}
	fmt.Println("bye")
}
