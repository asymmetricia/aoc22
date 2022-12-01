package aoc

import (
	"bytes"
	_ "embed"
	"image"
	"image/color"
	"image/draw"
	"unicode/utf8"
)

//go:embed "font.txt"
var fontData []byte

const LineHeight = 12
const GlyphWidth = 8

type Glyph struct {
	Image         draw.Image
	Top, Left     int
	Width, Height int
}

var Glyphs = map[rune]Glyph{}

func init() {
	glyphdata := bytes.Split(bytes.TrimSpace(fontData), []byte("\n\n"))
	for _, gtxt := range glyphdata {
		g := Glyph{}

		rows := bytes.Split(gtxt, []byte("\n"))

	colsLeft:
		for x := 0; x < len(rows[1]); x++ {
			for _, row := range rows[1:] {
				if row[x] == '#' {
					break colsLeft
				}
			}
			g.Left++
		}

	colsRight:
		for x := len(rows[1]) - 1; x > 0; x-- {
			for _, row := range rows[1:] {
				if row[x] == '#' {
					break colsRight
				}
			}
			g.Width = len(rows[1]) - x - g.Left
		}

		r, _ := utf8.DecodeRune(rows[0])
		rows = rows[1:]
		g.Image = image.NewRGBA(image.Rect(0, 0, len(rows[0]), len(rows)))
		draw.Draw(g.Image, g.Image.Bounds(), image.Transparent, image.Point{}, draw.Src)
		for y, row := range rows {
			for x, pt := range row {
				switch pt {
				case '#':
					g.Image.Set(x, y, color.White)
				}
			}
		}
		Glyphs[r] = g
	}
}

type TypesetOpts struct {
	Scale int
	Kern  bool
}

// Typeset sets the given text on the image starting with the first glyph's (0,0)
// pixel at cursor. It returns the number of pixels wide the text is.
func Typeset(img draw.Image, cursor image.Point, line string, color color.Color, opts ...TypesetOpts) int {
	left := cursor.X
	right := cursor.X

	if len(opts) == 0 {
		opts = []TypesetOpts{{}}
	}

	scale := 1
	if opts[0].Scale != 0 {
		scale = opts[0].Scale
	}

	initX := cursor.X
	for _, g := range line {
		switch g {
		case '\n':
			cursor.Y += LineHeight * scale
			cursor.X = initX
		default:
			glyph, ok := Glyphs[g]
			if ok {
				for x := 0; x < glyph.Image.Bounds().Size().X*scale; x++ {
					for y := 0; y < glyph.Image.Bounds().Size().Y*scale; y++ {
						c := glyph.Image.At(x/scale, y/scale)
						_, _, _, a := c.RGBA()
						if a > 0 {
							img.Set(cursor.X+x, cursor.Y+y, color)
						}
					}
				}
			}
			if ok && opts[0].Kern {
				cursor.X += glyph.Width
			} else {
				cursor.X += 8 * scale
			}
			if cursor.X > right {
				right = cursor.X
			}
		}
	}

	return right - left
}
