package framebuffer

import (
	"image"
	"image/color"

	"github.com/asymmetricia/aoc22/aoc"
	"github.com/asymmetricia/aoc22/term"
)

type Cell struct {
	Color color.Color
	Value rune
}

type Framebuffer [][]Cell

func (f *Framebuffer) Set(x, y int, value Cell) {
	for y >= len(*f) {
		*f = append(*f, nil)
	}
	if x >= len((*f)[y]) {
		row := make([]Cell, x+1)
		copy(row, (*f)[y])
		(*f)[y] = row
	}
	(*f)[y][x] = value
}

func (f *Framebuffer) PrintAt(x, y int, s string, c color.Color) {
	i := 0
	for _, char := range s {
		if char == '\n' {
			i = 0
			y++
			continue
		}
		f.Set(x+i, y, Cell{c, char})
		i++
	}
}

func (f *Framebuffer) TypeSet() *image.Paletted {
	max := aoc.MaxFn(*f, func(c []Cell) int { return len(c) })
	img := image.NewPaletted(image.Rect(0, 0, max*aoc.GlyphWidth, len(*f)*aoc.LineHeight), aoc.TolVibrant)
	for y, row := range *f {
		var c color.Color
		var accum []rune
		var x int
		for _, cell := range row {
			if cell.Color == nil {
				cell.Color = c
				cell.Value = ' '
			}
			if c != nil && cell.Color != c && len(accum) > 0 {
				aoc.Typeset(img, image.Pt(x*aoc.GlyphWidth, y*aoc.LineHeight), string(accum), c)
				x += len(accum)
				accum = accum[0:0]
			}
			c = cell.Color
			accum = append(accum, cell.Value)
		}
		if len(accum) > 0 && c != nil {
			aoc.Typeset(img, image.Pt(x*aoc.GlyphWidth, y*aoc.LineHeight), string(accum), c)
		}
	}
	return img
}

func (f *Framebuffer) String() string {
	var ret string
	var c color.Color
	var accum []rune
	for _, row := range *f {
		for _, cell := range row {
			if cell.Color == nil {
				cell.Color = c
				cell.Value = ' '
			}
			if c != nil && cell.Color != c && len(accum) > 0 {
				ret += term.ScolorC(c) + string(accum)
				accum = accum[0:0]
			}
			c = cell.Color
			accum = append(accum, cell.Value)
		}
		accum = append(accum, '\n')
	}
	if len(accum) > 0 {
		ret += term.ScolorC(c) + string(accum)
	}
	return ret
}

func (f *Framebuffer) Copy() Framebuffer {
	var ret Framebuffer
	ret = make([][]Cell, len(*f))
	for i, row := range *f {
		(ret)[i] = make([]Cell, len(row))
		copy((ret)[i], (*f)[i])
	}
	return ret
}

func (f *Framebuffer) BlockPrintAt(x, y int, s string, c color.Color) {
	xx := 0
	for _, char := range s {
		f.BlockSet(x+xx*aoc.GlyphWidth, y, Cell{c, char})
		xx++
	}
}

func (f *Framebuffer) BlockSet(x, y int, value Cell) {
	for yy := 0; yy < aoc.LineHeight; yy++ {
		for xx := 0; xx < aoc.GlyphWidth; xx++ {
			f.Set(x+xx, y+yy, Cell{value.Color, ' '})
		}
	}

	if value.Value == ' ' {
		return
	}

	glyph, ok := aoc.Glyphs[value.Value]
	if !ok {
		glyph = aoc.Glyphs['?']
	}
	for yy, row := range glyph.Raw {
		for xx, set := range row {
			var v = ' '
			if set {
				v = '#'
			}
			f.Set(x+xx, y+yy, Cell{value.Color, v})
		}
	}
}

func (f *Framebuffer) Rect() image.Rectangle {
	x := aoc.MaxFn(*f, func(cs []Cell) int { return len(cs) })
	return image.Rect(0, 0, x, len(*f))
}
