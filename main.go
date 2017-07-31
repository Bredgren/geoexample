package main

import (
	"fmt"
	"image/color"
	"time"

	"github.com/Bredgren/geo"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

const (
	Width  = 320
	Height = 240
)

type Square struct {
	img  *ebiten.Image
	opts ebiten.DrawImageOptions
}

var (
	square Square
)

var (
	buttonDown     bool
	buttonJustDown bool
)

var (
	currentOption int = 0
)

type DisplayFunc func(target *ebiten.Image)

type example struct {
	key  ebiten.Key
	name string
	fn   DisplayFunc
}

var options = []example{
	{ebiten.Key1, "Ease Functions", easeFunctions},
	{ebiten.Key2, "Perlin", perlin},
}

func drawOptions(target *ebiten.Image) {
	ebitenutil.DebugPrint(target,
		fmt.Sprintf("Press a number. Current: %s", options[currentOption].name))
	for i, ex := range options {
		newLines := "\n"
		for l := 0; l < i; l++ {
			newLines += "\n"
		}
		ebitenutil.DebugPrint(target, fmt.Sprintf("%s%d - %s", newLines, ex.key, ex.name))
	}
}

func update(screen *ebiten.Image) error {
	pressed := ebiten.IsKeyPressed(ebiten.KeyF)
	buttonJustDown = pressed && !buttonDown
	buttonDown = pressed

	if buttonJustDown {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}

	screen.Fill(color.NRGBA{0xff, 0x00, 0x00, 0xff})

	// square.img.Fill(color.White)
	//
	// for y := 0.0; y < 480; y += 32 {
	// 	for x := 0.0; x < 640; x += 32 {
	// 		square.opts.GeoM.Scale(31, 31)
	// 		square.opts.GeoM.Translate(x, y)
	// 		screen.DrawImage(square.img, &square.opts)
	// 		square.opts.GeoM.Reset()
	// 	}
	// }

	options[currentOption].fn(screen)

	drawOptions(screen)

	return nil
}

const (
	easeTime = 4 * time.Second
	easeWait = 1 * time.Second
	easeSize = 10.0
)

var (
	easeStart = time.Now()
	easeFns   = []geo.EaseFn{
		geo.EaseLinear,
		geo.EaseInQuad,
		geo.EaseOutQuad,
		geo.EaseInOutQuad,

		geo.EaseInElastic,
		geo.EaseOutElastic,
		geo.EaseInOutElastic,
		geo.EaseInBack,
		geo.EaseOutBack,
		geo.EaseInOutBack,
		geo.EaseInBounce,
		geo.EaseOutBounce,
		geo.EaseInOutBounce,
	}
)

func easeFunctions(target *ebiten.Image) {
	square.img.Fill(color.White)

	now := time.Now()
	dt := now.Sub(easeStart)
	if dt > easeTime+easeWait {
		easeStart = now.Add(easeWait)
		dt = 0
	}

	t := geo.Clamp(dt.Seconds()/easeTime.Seconds(), 0, 1)

	startY := 50.0
	start, end := geo.VecXY(100, startY), geo.VecXY(Width-20, startY)
	offset := geo.VecXY(0, easeSize*1.2)

	for _, fn := range easeFns {
		pos := geo.EaseVec(start, end, t, fn)
		square.opts.GeoM.Reset()
		square.opts.GeoM.Scale(easeSize, easeSize)
		square.opts.GeoM.Translate(pos.XY())
		target.DrawImage(square.img, &square.opts)
		start.Add(offset)
		end.Add(offset)
	}
}

func perlin(target *ebiten.Image) {
}

func main() {
	// On X1 Yoga 2560x1440, at any scale it seems
	//  - 640x480 becomes 960x720, at 150% scale (not including title bar)
	//  - title bar ovlaps draw area, but draw area is properly scaled
	// On Desktop 1920x1080, 100% scale
	//  - draw area is proper size (640x480)
	//  - title bar does not overlap draw area
	square.img, _ = ebiten.NewImage(1, 1, ebiten.FilterNearest)

	if err := ebiten.Run(update, Width, Height, 2, "Geo Examples"); err != nil {
		panic(err)
	}
	fmt.Println("bye")
}
