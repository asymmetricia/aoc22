package isovox

import (
	"image"
	"image/color"
	"sync"

	"github.com/asymmetricia/pencil"
)

type Voxel struct {
	Color  color.Color
	sprite image.Image
	size   int
}

func (v *Voxel) Sprite(size int) image.Image {
	r, g, b, a := v.Color.RGBA()
	colorCK := color.RGBA64{
		R: uint16(r),
		G: uint16(g),
		B: uint16(b),
		A: uint16(a),
	}

	spriteCacheMu.RLock()
	sc, ok := spriteCache[size]
	spriteCacheMu.RUnlock()
	if !ok {
		spriteCacheMu.Lock()
		spriteCache[size] = map[color.RGBA64]image.Image{}
		sc = spriteCache[size]
		spriteCacheMu.Unlock()
	}

	if img, ok := sc[colorCK]; ok {
		return img
	}

	dy := dy(size)
	dx := dx(size)
	h := size + dy*2
	w := dx * 2

	sprite := image.NewRGBA64(image.Rect(0, 0, w+1, h+1))
	top, left, right, edge := v.colors()

	center := image.Pt(dx, size)
	twelve := image.Pt(dx, 0)
	two := image.Pt(2*dx, dy)
	four := image.Pt(2*dx, size+dy)
	six := image.Pt(dx, size*2)
	eight := image.Pt(0, size+dy)
	ten := image.Pt(0, dy)

	/*      1 2
	        .
	10    /   \     2
	    /       \
	   |\       /|
	   |  \   /  |
	8  |    c    |  4
	    \   |   /
	      \ | /
	        |
	        6
	*/

	for tri, col := range map[[3]image.Point]color.Color{
		{six, center, four}:   right,
		{four, center, two}:   right,
		{center, twelve, two}: top,
		{center, ten, twelve}: top,
		{six, eight, center}:  left,
		{center, eight, ten}:  left,
	} {
		pencil.FillTriangle(tri[0], tri[1], tri[2], col, sprite)
	}

	for _, edgePt := range [][2]image.Point{
		{six, eight},
		{eight, ten},
		{ten, twelve},
		{twelve, two},
		{two, four},
		{four, six},
		{six, center},
		{ten, center},
		{two, center},
	} {
		pencil.Line(sprite, edgePt[0], edgePt[1], edge)
	}

	spriteCacheMu.Lock()
	spriteCache[size][colorCK] = sprite
	spriteCacheMu.Unlock()
	return sprite
}

var spriteCacheMu = &sync.RWMutex{}
var spriteCache = map[int]map[color.RGBA64]image.Image{}
