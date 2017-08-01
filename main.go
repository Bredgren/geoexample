package main

import (
	"fmt"
	"image"
	"image/color"
	"math"
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

func checkOptions() {
	for i, option := range options {
		if ebiten.IsKeyPressed(option.key) {
			currentOption = i
		}
	}
}

func drawOptions(target *ebiten.Image) {
	ebitenutil.DebugPrint(target,
		fmt.Sprintf("\nPress a number. Current: %s", options[currentOption].name))
	for i, ex := range options {
		newLines := "\n\n"
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

	checkOptions()

	if ebiten.IsRunningSlowly() {
		return nil
	}

	options[currentOption].fn(screen)

	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f", ebiten.CurrentFPS()))
	drawOptions(screen)

	return nil
}

const (
	easeTime = 4 * time.Second
	easeWait = 1 * time.Second
	easeSize = 5.0
)

var (
	easeStart = time.Now()
	easeFns   = []geo.EaseFn{
		geo.EaseLinear,

		geo.EaseInQuad,
		geo.EaseInCubic,
		geo.EaseInQuart,
		geo.EaseInQuint,
		geo.EaseInSine,
		geo.EaseInCirc,
		geo.EaseInExpo,
		geo.EaseInElastic,
		geo.EaseInBack,
		geo.EaseInBounce,

		geo.EaseOutQuad,
		geo.EaseOutCubic,
		geo.EaseOutQuart,
		geo.EaseOutQuint,
		geo.EaseOutSine,
		geo.EaseOutCirc,
		geo.EaseOutExpo,
		geo.EaseOutElastic,
		geo.EaseOutBack,
		geo.EaseOutBounce,

		geo.EaseInOutQuad,
		geo.EaseInOutCubic,
		geo.EaseInOutQuart,
		geo.EaseInOutQuint,
		geo.EaseInOutSine,
		geo.EaseInOutCirc,
		geo.EaseInOutExpo,
		geo.EaseInOutElastic,
		geo.EaseInOutBack,
		geo.EaseInOutBounce,
	}
)

func easeFunctions(dst *ebiten.Image) {
	square.img.Fill(color.White)

	now := time.Now()
	dt := now.Sub(easeStart)
	if dt > easeTime+easeWait {
		easeStart = now.Add(easeWait)
		dt = 0
	}

	t := geo.Clamp(dt.Seconds()/easeTime.Seconds(), 0, 1)

	startY := 30.0
	start, end := geo.VecXY(100, startY), geo.VecXY(Width-20, startY)
	offset := geo.VecXY(0, easeSize*1.2)

	for i, fn := range easeFns {
		pos := geo.EaseVec(start, end, t, fn)
		square.opts.GeoM.Reset()
		square.opts.GeoM.Scale(easeSize, easeSize)
		square.opts.GeoM.Translate(pos.XY())
		dst.DrawImage(square.img, &square.opts)
		start.Add(offset)
		end.Add(offset)
		if i%10 == 0 {
			start.Y += easeSize
			end.Y += easeSize
		}
		if i == 0 {
			square.img.Fill(color.NRGBA{0x0, 0x0, 0xff, 0xff})
		}
		square.opts.ColorM.Reset()
		square.opts.ColorM.RotateHue(float64(i%10) / 10 * 2 * math.Pi)
	}
}

var (
	perlinImg      *image.RGBA
	perlinZ        = 0.0
	perlinRate     = 0.3
	perlinPrevTime = time.Now()
)

func perlin(dst *ebiten.Image) {
	w, h := dst.Size()
	if perlinImg == nil {
		perlinImg = image.NewRGBA(image.Rect(0, 0, w, h))

	}
	for i := 0; i < w*h; i++ {
		x, y := float64(i%w), float64(i/w)
		val := geo.PerlinOctave(x*0.01, y*0.01, perlinZ, 5, 0.5)
		perlinImg.Pix[4*i] = uint8(0xff * val)
		perlinImg.Pix[4*i+1] = uint8(0xff * val)
		perlinImg.Pix[4*i+2] = uint8(0xff * val)
		perlinImg.Pix[4*i+3] = 0xff
	}
	dst.ReplacePixels(perlinImg.Pix)
	t := time.Now()
	dt := perlinPrevTime.Sub(t).Seconds()
	perlinZ += dt * perlinRate
	perlinPrevTime = t
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
