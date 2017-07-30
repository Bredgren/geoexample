package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

var square *ebiten.Image

var (
	buttonDown     bool
	buttonJustDown bool
)

func update(screen *ebiten.Image) error {
	pressed := ebiten.IsKeyPressed(ebiten.KeyF)
	buttonJustDown = pressed && !buttonDown
	buttonDown = pressed

	if buttonJustDown {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}

	screen.Fill(color.NRGBA{0xff, 0x00, 0x00, 0xff})
	ebitenutil.DebugPrint(screen, "Hello")

	if buttonDown {
		ebitenutil.DebugPrint(screen, "\n\n\nF")
	}

	if square == nil {
		square, _ = ebiten.NewImage(1, 1, ebiten.FilterNearest)
	}

	square.Fill(color.White)

	opts := &ebiten.DrawImageOptions{}

	// squares := []geo.Rect{
	// 	geo.RectWH(32, 32), geo.RectWH(32, 32), geo.RectWH(32, 32), geo.RectWH(32, 32),
	// }
	// squares[0].SetTopLeft(0, 0)
	// squares[1].SetTopRight(640, 0)
	// squares[2].SetBottomLeft(0, 480)
	// squares[3].SetBottomRight(640, 480)

	// for _, s := range squares {
	// 	opts.GeoM.Translate(s.TopLeft())
	// 	screen.DrawImage(square, opts)
	// 	opts.GeoM.Reset()
	// }

	for y := 0.0; y < 480; y += 32 {
		for x := 0.0; x < 640; x += 32 {
			opts.GeoM.Scale(31, 31)
			opts.GeoM.Translate(x, y)
			screen.DrawImage(square, opts)
			opts.GeoM.Reset()
		}
	}

	return nil
}

func main() {
	// On X1 Yoga 2560x1440, at any scale it seems
	//  - 640x480 becomes 960x720, at 150% scale (not including title bar)
	//  - title bar ovlaps draw area, but draw area is properly scaled
	// On Desktop 1920x1080, 100% scale
	//  - draw area is proper size (640x480)
	//  - title bar does not overlap draw area
	if err := ebiten.Run(update, 640, 480, 1, "Hello World!"); err != nil {
		panic(err)
	}
	fmt.Println("bye")
}
