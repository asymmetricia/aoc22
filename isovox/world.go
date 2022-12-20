package isovox

import (
	"image"
	"image/color"
	"image/draw"
	"math"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"

	"github.com/asymmetricia/aoc22/aoc"
)

type Coord struct {
	X, Y, Z int
}

func (v *Voxel) colors() (top, left, right, edge color.Color) {
	r, g, b, a := v.Color.RGBA()

	cl := color.RGBA64{
		R: uint16(r * 90 / 100),
		G: uint16(g * 90 / 100),
		B: uint16(b * 90 / 100),
		A: uint16(a),
	}

	cr := color.RGBA64{
		R: uint16(r * 75 / 100),
		G: uint16(g * 75 / 100),
		B: uint16(b * 75 / 100),
		A: uint16(a),
	}

	ce := color.NRGBA64Model.Convert(v.Color).(color.NRGBA64)
	for _, v := range []*uint16{&ce.R, &ce.G, &ce.B} {
		if *v >= 0xFFFF/105*100 {
			*v = 0xFFFF
		} else {
			*v = uint16(uint32(*v) * 105 / 100)
		}
	}

	return color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)}, cl, cr, ce
}

type World struct {
	Voxels map[Coord]*Voxel
}

func (w *World) Bounds(size int) (screen image.Rectangle, min Coord, max Coord) {
	dy := dy(size)
	dx := dx(size)

	min = Coord{X: math.MaxInt, Y: math.MaxInt, Z: math.MaxInt}
	max = Coord{X: math.MinInt, Y: math.MinInt, Z: math.MinInt}
	screen = image.Rectangle{
		Min: image.Point{math.MaxInt, math.MaxInt},
		Max: image.Point{math.MinInt, math.MinInt},
	}
	for c := range w.Voxels {
		min.X = aoc.Min(min.X, c.X)
		min.Y = aoc.Min(min.Y, c.Y)
		min.Z = aoc.Min(min.Z, c.Z)
		max.X = aoc.Min(max.X, c.X)
		max.Y = aoc.Min(max.Y, c.Y)
		max.Z = aoc.Min(max.Z, c.Z)

		x := c.X*dx + c.Y*dx
		y := c.X*dy - c.Y*dy - c.Z*size
		if x < screen.Min.X {
			screen.Min.X = x
		}
		if x+2*dx+1 > screen.Max.X {
			screen.Max.X = x + 2*dx + 1
		}
		if y-size-dy < screen.Min.Y {
			screen.Min.Y = y - size - dy
		}
		if y+dy+1 > screen.Max.Y {
			screen.Max.Y = y + dy + 1
		}
	}

	return screen, min, max
}

func (w *World) Render(size int) image.Image {
	dy := dy(size)
	dx := dx(size)

	r, _, _ := w.Bounds(size)

	voxelCoords := maps.Keys(w.Voxels)
	slices.SortFunc(voxelCoords, func(a, b Coord) bool {
		// front-most is highest dx, lowest dy, highest dz
		depthA := a.X - a.Y + a.Z
		depthB := b.X - b.Y + b.Z
		return depthA < depthB
	})

	imgRect := r.Sub(r.Min)
	if imgRect.Dx()%2 == 1 {
		imgRect.Max.X++
	}
	if imgRect.Dy()%2 == 1 {
		imgRect.Max.Y++
	}
	ret := image.NewRGBA64(imgRect)
	draw.Draw(ret, ret.Bounds(), image.Black, image.Pt(0, 0), draw.Over)
	for _, coord := range voxelCoords {
		vx := coord.X*dx + coord.Y*dx
		vy := coord.X*dy - coord.Y*dy - coord.Z*size - size - dy
		sprite := w.Voxels[coord].Sprite(size)
		draw.Draw(ret, sprite.Bounds().Add(image.Pt(vx, vy).Sub(r.Min)), sprite, image.Pt(0, 0), draw.Over)
	}
	return ret
}

func dx(size int) int {
	return int(math.Cos(math.Pi*30/180) * float64(size))
}

func dy(size int) int {
	return int(math.Sin(math.Pi*30/180) * float64(size))
}
