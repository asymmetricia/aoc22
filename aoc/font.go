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
	Raw           [][]bool
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

		var glyphKey []byte
		accum := make([]uint8, 0, 2)
		var state int
		for _, keyByte := range rows[0] {
			switch state {
			case 0:
				if keyByte == '\\' {
					state++
					continue
				}
				glyphKey = append(glyphKey, keyByte)
			case 1:
				if keyByte == 'x' {
					state++
					continue
				} else {
					glyphKey = append(glyphKey, '\\', keyByte)
					state = 0
				}
			case 2:
				if keyByte >= '0' && keyByte <= '9' {
					keyByte -= '0'
				} else if keyByte >= 'a' && keyByte <= 'f' {
					keyByte -= 'a' + 10
				} else if keyByte >= 'A' && keyByte <= 'F' {
					keyByte -= 'A' + 10
				} else {
					panic("expected hex digit, found " + string(keyByte))
				}
				accum = append(accum, keyByte)
				if len(accum) == 2 {
					glyphKey = append(glyphKey, accum[0]<<4|accum[1])
				}
			case 3:
				accum[1] = keyByte
			}
		}
		r, s := utf8.DecodeRune(glyphKey)
		if r == '\000' || s == 0 {
			r = ' '
		}
		rows = rows[1:]
		g.Image = image.NewRGBA(image.Rect(0, 0, len(rows[0]), len(rows)))
		draw.Draw(g.Image, g.Image.Bounds(), image.Transparent, image.Point{}, draw.Src)
		g.Raw = make([][]bool, len(rows))
		for y, row := range rows {
			g.Raw[y] = make([]bool, len(row))
			for x, pt := range row {
				if pt == '#' {
					g.Image.Set(x, y, color.White)
					g.Raw[y][x] = true
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

func TypesetBytes(line string, opts ...TypesetOpts) [][]byte {
	opt := TypesetOpts{1, false}
	if len(opts) > 0 {
		opt = opts[0]
	}

	var ret [][]byte
	cursorX := 0
	y := 0
	for _, g := range line {
		if g == '\n' {
			y++
			cursorX = 0
			continue
		}
		glyph, ok := Glyphs[g]
		if !ok {
			glyph = Glyphs['?']
		}
		for glyphY, glyphRow := range glyph.Raw {
			canvasY := y*LineHeight*opt.Scale + glyphY*opt.Scale
			for len(ret) < canvasY+opt.Scale {
				ret = append(ret, nil)
			}
			for glyphX, bit := range glyphRow {
				canvasX := cursorX*GlyphWidth*opt.Scale + glyphX*opt.Scale
				for dy := 0; dy < opt.Scale; dy++ {
					for len(ret[canvasY+dy]) < canvasX+opt.Scale {
						ret[canvasY+dy] = append(ret[canvasY+dy], ' ')
					}
					if !bit {
						continue
					}
					for dx := 0; dx < opt.Scale; dx++ {
						ret[canvasY+dy][canvasX+dx] = '#'
					}
				}
			}
		}
		cursorX++
	}
	return ret
}

func TypesetString(line string, opts ...TypesetOpts) string {
	return string(bytes.Join(TypesetBytes(line, opts...), []byte{'\n'}))
}

// Typeset sets the given text on the image starting with the first glyph's (0,0)
// pixel at cursor. It returns the number of pixel/s wide the text is.
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
