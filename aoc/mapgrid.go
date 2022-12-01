package aoc

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math"
	"strconv"

	"github.com/asymmetricia/aoc21/coord"
)

type MapGrid map[int]map[int]int

func (m MapGrid) Inc(c coord.Coord) {
	m.Set(c, m.Get(c)+1)
}

func (m MapGrid) Set(c coord.Coord, n int) {
	if _, ok := m[c.Y]; !ok {
		m[c.Y] = map[int]int{}
	}
	m[c.Y][c.X] = n
}

func (m MapGrid) Get(c coord.Coord) int {
	if _, ok := m[c.Y]; !ok {
		return 0
	}
	return m[c.Y][c.X]
}

func (m MapGrid) Print() {
	strings := map[int]map[int]string{}
	minX, minY := math.MaxInt, math.MaxInt
	maxX, maxY := math.MinInt, math.MinInt
	width := 0
	for y, row := range m {
		if _, ok := strings[y]; !ok {
			strings[y] = map[int]string{}
		}
		for x, cell := range row {
			if x < minX {
				minX = x
			}
			if x > maxX {
				maxX = x
			}
			if y < minY {
				minY = y
			}
			if y > maxY {
				maxY = y
			}
			s := strconv.Itoa(cell)
			strings[y][x] = s
			if len(s) > width {
				width = len(s)
			}
		}
	}

	fmt.Println(minX, maxX, minY, maxY)
	space := ""
	for i := 0; i < width; i++ {
		space += " "
	}
	fmtI := "%" + strconv.Itoa(width) + "d"

	for y := minY; y < maxY; y++ {
		row, ok := m[y]
		if ok {
			for x := minX; x < maxX; x++ {
				cell, ok := row[x]
				if ok {
					fmt.Printf(fmtI, cell)
				} else {
					fmt.Print(space)
				}
			}
		}
		fmt.Print("\n")
	}
}

func (m MapGrid) Draw(r image.Rectangle, bg color.Color, cellToColorIndex map[int]uint8, palette color.Palette) *image.Paletted {
	target := image.NewPaletted(r.Sub(r.Min), palette)
	draw.Draw(target, target.Bounds(), image.NewUniform(bg), image.Point{}, draw.Over)
	for y, row := range m {
		for x, cell := range row {
			idx, ok := cellToColorIndex[cell]
			if !ok {
				continue
			}
			target.SetColorIndex(x-r.Min.X, y-r.Min.Y, idx)
		}
	}
	return target
}
