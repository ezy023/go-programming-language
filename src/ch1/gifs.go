package main

import (
	"image"
	"image/color"
	"image/gif"
	"io"
	"math"
	"math/rand"
)

var green color.Color = color.RGBA{0, 184, 0, 0xff}
var blue color.Color = color.RGBA{0, 0, 255, 0xff}
var red color.Color = color.RGBA{255, 0, 0, 0xff}
var palette = []color.Color{color.Black, green, blue, red}

const (
	background = 0 // first color in palette
	primary    = 1 // next color in palette
	secondary  = 2
	tertiary   = 3
)

// func main() {
// 	lissajous(os.Stdout)
// }

func lissajous(out io.Writer, cycles float64) {
	const (
		// cycles = 5     // number of complet x oscillator revolutions
		res     = 0.001 // angular resolution
		size    = 100   // image canvas covers [-size..+size]
		nframes = 64    // number of animation frames
		delay   = 8     // delay between frams in 10ms units
	)

	freq := rand.Float64() * 3.0 // relative frequency of y oscillator
	var anim gif.GIF = gif.GIF{LoopCount: nframes}
	phase := 0.0 // phase difference
	for i := 0; i < nframes; i++ {
		rect := image.Rect(0, 0, 2*size+1, 2*size+1)
		img := image.NewPaletted(rect, palette)
		for t := 0.0; t < cycles*2*math.Pi; t += res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)
			xin := size + int(x*size+0.5)
			yin := size + int(y*size+0.5)
			var colorIdx uint8
			switch {
			case xin < 70:
				colorIdx = 1
			case xin >= 70 && xin < 140:
				colorIdx = 2
			case xin >= 140:
				colorIdx = 3
			}
			img.SetColorIndex(xin, yin, colorIdx)
		}
		phase += 0.1
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}
	gif.EncodeAll(out, &anim)
}
