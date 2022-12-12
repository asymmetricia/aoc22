package coord

import (
	"fmt"
	"math"
	"strings"
)

type SparseWorld map[Coord]rune

func (w SparseWorld) Find(r rune) []Coord {
	//TODO implement me
	panic("implement me")
}

func (w SparseWorld) Copy() World {
	r := make(SparseWorld, len(w))
	for c, obj := range w {
		r[c] = obj
	}
	return r
}

func (w SparseWorld) Rect() (minX, minY, maxX, maxY int) {
	minX, maxX, minY, maxY = math.MaxInt, math.MinInt, math.MaxInt, math.MinInt
	for c := range w {
		if c.X < minX {
			minX = c.X
		}
		if c.X > maxX {
			maxX = c.X
		}
		if c.Y < minY {
			minY = c.Y
		}
		if c.Y > maxX {
			maxX = c.Y
		}
	}
	return minX, minY, maxX, maxY
}

func (w SparseWorld) Print() {
	minx, miny, maxx, maxy := w.Rect()

	for y := miny; y <= maxy; y++ {
		sb := strings.Builder{}
		for x := minx; x <= maxx; x++ {
			if ch, ok := w[C(x, y)]; ok {
				sb.WriteRune(ch)
			} else {
				sb.WriteRune(' ')
			}
		}
		fmt.Println(sb.String())
	}
}

func (w SparseWorld) At(coord Coord) rune {
	if r, ok := w[coord]; !ok {
		return -1
	} else {
		return r
	}
}

func (w SparseWorld) Set(coord Coord, r rune) {
	w[coord] = r
}

func (w SparseWorld) Each(f func(Coord) bool) {
	for c := range w {
		if f(c) {
			return
		}
	}
}

type DenseWorld [][]rune

func (d DenseWorld) Find(r rune) []Coord {
	var ret []Coord
	for y, row := range d {
		for x, rune := range row {
			if rune == r {
				ret = append(ret, C(x, y))
			}
		}
	}
	return ret
}

func (d DenseWorld) Copy() World {
	r := make(DenseWorld, len(d))
	for i, row := range d {
		newRow := make([]rune, len(row))
		copy(newRow, row)
		r[i] = newRow
	}
	return &r
}

func (d DenseWorld) Rect() (minX, minY, maxX, maxY int) {
	return 0, 0, len(d[0]) - 1, len(d) - 1
}

func (d DenseWorld) Print() {
	sb := &strings.Builder{}
	for _, row := range d {
		sb.Reset()
		for _, cell := range row {
			sb.WriteRune(cell)
		}
		fmt.Println(sb.String())
	}
}

func (d DenseWorld) At(coord Coord) rune {
	if len(d) <= coord.Y || coord.Y < 0 {
		return -1
	}
	if len(d[coord.Y]) <= coord.X || coord.X < 0 {
		return -1
	}
	return d[coord.Y][coord.X]
}

func (d *DenseWorld) Set(coord Coord, r rune) {
	height := len(*d)
	if height <= coord.Y {
		*d = append(*d, make([][]rune, height-coord.Y+1)...)
	}
	width := len((*d)[coord.Y])
	if width <= coord.X {
		(*d)[coord.Y] = append((*d)[coord.Y], make([]rune, coord.X-width+1)...)
	}
	(*d)[coord.Y][coord.X] = r
}

func (d *DenseWorld) Each(f func(Coord) (stop bool)) {
	for y, row := range *d {
		for x := range row {
			if f(C(x, y)) {
				return
			}
		}
	}
}

type World interface {
	Print()
	At(Coord) rune
	Set(Coord, rune)
	Each(func(Coord) (stop bool))
	Rect() (minX, minY, maxX, maxY int)
	Copy() World
	Find(rune) []Coord
}

var _ World = (*SparseWorld)(nil)
var _ World = (*DenseWorld)(nil)

// Load loads a world from the given list of lines and returns it. The world is
// dense (array-based) if `dense` is true, otherwise it's sparse (map-based).
func Load(lines []string, dense bool) World {
	var w World
	if dense {
		w = new(DenseWorld)
	} else {
		w = new(SparseWorld)
	}
	for y, row := range lines {
		for x, char := range row {
			w.Set(C(x, y), char)
		}
	}
	return w
}
