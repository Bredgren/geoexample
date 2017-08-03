# geoexample

This projects purpose is to be an example/demo and to test some things for the [geo](https://github.com/Bredgren/geo) package.

## Here are the examples as gifs

### Ease Functions
```go
dt := time.Now().Sub(easeStartTime)
duration := 1 * time.Second
t := geo.Clamp(dt.Seconds()/duration.Seconds(), 0, 1)
startPos := geo.VecXY(0, 0)
endPos := geo.VecXY(10, 10)
pos := geo.EaseVec(startPos, endPos, t, geo.EaseOutQuad)
```
![alt text](https://github.com/Bredgren/geoexample/blob/master/gif/ease.gif "Ease")

### Perlin Noise
```go
for i := 0; i < w*h; i++ {
	x, y := float64(i%w), float64(i/w)
	val := geo.PerlinOctave(x*0.01, y*0.01, 0, 5, 0.5)
	perlinImg.Pix[4*i] = uint8(0xff * val)
	perlinImg.Pix[4*i+1] = uint8(0xff * val)
	perlinImg.Pix[4*i+2] = uint8(0xff * val)
	perlinImg.Pix[4*i+3] = 0xff
}
dst.ReplacePixels(perlinImg.Pix)
```
![alt text](https://github.com/Bredgren/geoexample/blob/master/gif/perlin.gif "Perlin")

### Shake
```go
shaker := geo.Shaker{
	Amplitude: 20,
	Frequency: 20,
	StartTime: time.Now(),
	Duration:  2 * time.Second,
	Falloff:   geo.EaseOutQuad,
}
actualPos := geo.VecXY(0, 0)
offset := shaker.Shake(time.Now())
pos := actualPos.Plus(offset)
...
// Start again later
if now.After(shaker.EndTime()) {
  shaker.StartTime = now
}
```
![alt text](https://github.com/Bredgren/geoexample/blob/master/gif/shake.gif "Shake")

### VecGen
```go
arcGen := geo.OffsetVec(geo.RandVecArc(30, 50, -math.Pi/2, math.Pi/4), geo.StaticVec(geo.VecXY(220, 70)))
points := [100]geo.Vec{}
for i := range points {
		points[i] = arcGen()
}
```
![alt text](https://github.com/Bredgren/geoexample/blob/master/gif/vecgen.gif "VecGen")

### VecMod
```go
func (m *movingBlock) update(t time.Time, bounds geo.Rect) {
	dt := t.Sub(m.lastUpdate)
	m.lastUpdate = t

	m.pos.Add(m.vel.Times(dt.Seconds()))
	m.pos.Mod(bounds)
}
```
![alt text](https://github.com/Bredgren/geoexample/blob/master/gif/vecmod.gif "VecMod")
