package main

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"math/rand"
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

type DisplayFunc func(dst *ebiten.Image)

type example struct {
	key  ebiten.Key
	name string
	fn   DisplayFunc
}

var options = []example{
	{ebiten.Key1, "Ease", easeFunctions},
	{ebiten.Key2, "Perlin", perlin},
	{ebiten.Key3, "Shake", shake},
	{ebiten.Key4, "VecGen", vecGen},
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

var (
	shakeStartTime = time.Now()
	shakyRect1     = shakyRect{
		rect: geo.RectXYWH(100, 100, 20, 20),
		shaker: geo.Shaker{
			Amplitude: 10,
			Frequency: 10,
		},
	}
	shakyRect2 = shakyRect{
		rect: geo.RectXYWH(150, 100, 20, 20),
		shaker: geo.Shaker{
			Amplitude: 20,
			Frequency: 20,
			Duration:  2 * time.Second,
			Falloff:   geo.EaseOutQuad,
		},
	}
	shakyRect3 = shakyRect{
		rect: geo.RectXYWH(200, 100, 20, 20),
		shaker: geo.Shaker{
			Amplitude: math.Pi / 3,
			Frequency: 10,
			Duration:  500 * time.Millisecond,
			Falloff:   geo.EaseOutQuad,
		},
	}
)

type shakyRect struct {
	rect    geo.Rect
	shaker  geo.Shaker
	shaking bool
}

func (s *shakyRect) update(t time.Time) {
	cursor := geo.VecXYi(ebiten.CursorPosition())

	shakeOver := t.After(s.shaker.EndTime())
	if shakeOver && s.rect.CollidePoint(cursor.XY()) {
		s.shaker.StartTime = t
		shakeOver = false
	}
	s.shaking = !shakeOver
}

func (s *shakyRect) draw(dst *ebiten.Image, offset geo.Vec, angle float64) {
	if s.shaking {
		square.img.Fill(color.NRGBA{0xff, 0x00, 0x00, 0xff})
	} else {
		square.img.Fill(color.White)
	}
	square.opts.GeoM.Reset()
	square.opts.GeoM.Scale(s.rect.W, s.rect.H)
	square.opts.GeoM.Translate(-s.rect.W/2, -s.rect.H/2)
	square.opts.GeoM.Rotate(angle)
	square.opts.GeoM.Translate(s.rect.W/2, s.rect.H/2)
	square.opts.GeoM.Translate(geo.VecXY(s.rect.TopLeft()).Plus(offset).XY())
	dst.DrawImage(square.img, &square.opts)
}

func shake(dst *ebiten.Image) {
	now := time.Now()

	shakyRect1.update(now)
	if shakyRect1.shaking {
		offset := shakyRect1.shaker.ShakeConst(now)
		shakyRect1.draw(dst, offset, 0)
	} else {
		shakyRect1.draw(dst, geo.Vec0, 0)
	}

	shakyRect2.update(now)
	if shakyRect2.shaking {
		offset := shakyRect2.shaker.Shake(now)
		shakyRect2.draw(dst, offset, 0)
	} else {
		shakyRect2.draw(dst, geo.Vec0, 0)
	}

	shakyRect3.update(now)
	if shakyRect3.shaking {
		offset := shakyRect3.shaker.Shake1(now)
		shakyRect3.draw(dst, geo.Vec0, offset)
	} else {
		shakyRect3.draw(dst, geo.Vec0, 0)
	}
}

const (
	pointSize = 2
)

var (
	vecGenInit = false
	points1    [250]geo.Vec
	vecGen1    = geo.OffsetVec(geo.RandVecCircle(0, 40), geo.StaticVec(geo.VecXY(120, 80)))
	points2    [150]geo.Vec
	vecGen2    = geo.OffsetVec(geo.RandVecArc(30, 50, -math.Pi/2, math.Pi/4), geo.StaticVec(geo.VecXY(220, 70)))
	points3    [300]geo.Vec
	vecGen3    = geo.RandVecRects([]geo.Rect{
		geo.RectXYWH(80, 130, 200, 10), // Top
		geo.RectXYWH(80, 190, 200, 10), // Bottom
		geo.RectXYWH(80, 140, 20, 50),  // Left
		geo.RectXYWH(260, 140, 20, 50), // Right
	})
)

func vecGen(dst *ebiten.Image) {
	if !vecGenInit {
		for i := range points1 {
			points1[i] = vecGen1()
		}
		for i := range points2 {
			points2[i] = vecGen2()
		}
		for i := range points3 {
			points3[i] = vecGen3()
		}
		vecGenInit = true
	}

	points1[rand.Intn(len(points1))] = vecGen1()
	points2[rand.Intn(len(points2))] = vecGen2()
	points3[rand.Intn(len(points3))] = vecGen3()

	square.img.Fill(color.NRGBA{0x88, 0x88, 0xff, 0x55})
	square.opts.GeoM.Reset()
	square.opts.GeoM.Scale(200, 70)
	square.opts.GeoM.Translate(80, 130)
	dst.DrawImage(square.img, &square.opts)

	square.img.Fill(color.Black)
	square.opts.GeoM.Reset()
	square.opts.GeoM.Scale(160, 50)
	square.opts.GeoM.Translate(100, 140)
	dst.DrawImage(square.img, &square.opts)

	square.img.Fill(color.White)
	for _, points := range [][]geo.Vec{points1[:], points2[:], points3[:]} {
		for _, p := range points {
			square.opts.GeoM.Reset()
			square.opts.GeoM.Scale(pointSize, pointSize)
			square.opts.GeoM.Translate(p.Minus(geo.VecXY(pointSize/2, pointSize/2)).XY())
			dst.DrawImage(square.img, &square.opts)
		}
	}
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
