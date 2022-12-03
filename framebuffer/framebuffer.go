package framebuffer

import (
	"image"
	"image/color"
	"strings"

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

func (f *Framebuffer) TypeSet(opts ...aoc.TypesetOpts) *image.Paletted {
	return f.TypeSetRect(image.Rectangle{}, opts...)
}

func (f *Framebuffer) TypeSetRect(minRect image.Rectangle, opts ...aoc.TypesetOpts) *image.Paletted {
	opt := aoc.TypesetOpts{Scale: 1}
	if len(opts) > 0 {
		opt = opts[0]
	}

	max := aoc.MaxFn(*f, func(c []Cell) int { return len(c) })
	minWidth := minRect.Dx() * aoc.GlyphWidth * opt.Scale
	width := max * aoc.GlyphWidth * opt.Scale
	minHeight := minRect.Dy() * aoc.LineHeight * opt.Scale
	height := len(*f) * aoc.LineHeight * opt.Scale
	if minWidth > width {
		width = minWidth
	}
	if minHeight > height {
		height = minHeight
	}

	img := image.NewPaletted(image.Rect(0, 0, width, height), aoc.TolVibrant)
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
				aoc.Typeset(img, image.Pt(x*aoc.GlyphWidth*opt.Scale, y*aoc.LineHeight*opt.Scale), string(accum), c, opt)
				x += len(accum)
				accum = accum[0:0]
			}
			c = cell.Color
			accum = append(accum, cell.Value)
		}
		if len(accum) > 0 && c != nil {
			aoc.Typeset(img, image.Pt(x*aoc.GlyphWidth*opt.Scale, y*aoc.LineHeight*opt.Scale), string(accum), c, opt)
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
	f.PrintAt(x, y, aoc.TypesetString(s), c)
}

func (f *Framebuffer) BlockSet(x, y int, value Cell) {
	f.PrintAt(x, y, aoc.TypesetString(string(value.Value)), value.Color)
}

func (f *Framebuffer) Rect() image.Rectangle {
	x := aoc.MaxFn(*f, func(cs []Cell) int { return len(cs) })
	return image.Rect(0, 0, x, len(*f))
}

type TextBox struct {
	// If Middle is true, Top is ignored and the box is placed vertically in the
	// middle of the existing framebuffer.
	Top    int
	Middle bool

	// If Center is true, Left is ignored and the box is placed horizontally in the
	// center of the existing framebuffer.
	Left   int
	Center bool

	Title           []rune
	TitleRightAlign bool

	Body      []rune
	BodyBlock bool
	// if BodyPad is true, left and right side of body will be padded in from the
	// frame. Padding will be one space, or one block space if BodyBlock is true.
	BodyPad bool

	Footer          []rune
	FooterLeftAlign bool

	// Defaults to aoc.TolVibrantGrey
	BodyColor color.Color

	// Defaults to same as BodyColor
	TitleColor color.Color

	// Defaults to same as BodyColor
	FrameColor color.Color

	// Defaults to same as **TitleColor**
	FooterColor color.Color
}

func (t TextBox) On(f *Framebuffer) {
	if t.BodyBlock {
		var blockBody string
		for _, line := range strings.Split(string(t.Body), "\n") {
			if t.BodyPad {
				line = " " + line + " "
			}
			if blockBody != "" {
				blockBody += "\n"
			}
			blockBody += aoc.TypesetString(line)
		}
		t.Body = []rune(blockBody)
		t.BodyBlock = false
		t.BodyPad = false
	}

	// compute body size
	bodyWidth := 0
	bodyHeight := 0
	for _, line := range strings.Split(string(t.Body), "\n") {
		if !t.BodyPad && len(line) > bodyWidth {
			bodyWidth = len(line)
		} else if t.BodyPad && len(line)+2 > bodyWidth {
			bodyWidth = len(line) + 2
		}
		bodyHeight++
	}

	if len(t.Title) > bodyWidth {
		t.Title = t.Title[0:bodyWidth]
	}

	if len(t.Footer) > bodyWidth {
		t.Footer = t.Footer[0:bodyWidth]
	}

	// handle middle or center positioning
	fRect := f.Rect()
	if t.Middle {
		t.Top = fRect.Dy()/2 - (bodyHeight+2)/2
	}
	if t.Center {
		t.Left = fRect.Dx()/2 - (bodyWidth+4)/2
	}

	if t.BodyColor == nil {
		t.BodyColor = aoc.TolVibrantGrey
	}
	if t.TitleColor == nil {
		t.TitleColor = t.BodyColor
	}
	if t.FrameColor == nil {
		t.FrameColor = t.BodyColor
	}
	if t.FooterColor == nil {
		t.FooterColor = t.TitleColor
	}

	// Draw the title, aligned as per
	f.Set(t.Left, t.Top, Cell{t.FrameColor, '┏'})
	titleStart := bodyWidth - len(t.Title)
	titleEnd := bodyWidth
	if !t.TitleRightAlign {
		titleStart = 0
		titleEnd = len(t.Title)
	}
	for dy := 0; dy < titleStart; dy++ {
		f.Set(t.Left+dy+1, t.Top, Cell{t.FrameColor, '━'})
	}
	for dy := titleStart; dy < titleEnd; dy++ {
		f.Set(t.Left+dy+1, t.Top, Cell{t.TitleColor, t.Title[dy-titleStart]})
	}
	for dy := titleEnd; dy < bodyWidth; dy++ {
		f.Set(t.Left+dy+1, t.Top, Cell{t.FrameColor, '━'})
	}
	f.Set(t.Left+bodyWidth+1, t.Top, Cell{t.FrameColor, '┓'})
	t.Top++

	for _, line := range strings.Split(string(t.Body), "\n") {
		lineRunes := []rune(line)
		f.Set(t.Left, t.Top, Cell{t.FrameColor, '┃'})
		padX := 0
		if t.BodyPad {
			padX = 1
		}
		for bodyX := 0; bodyX < bodyWidth; bodyX++ {
			var r rune = ' '
			if bodyX < len(lineRunes) {
				r = lineRunes[bodyX]
			}
			f.Set(t.Left+1+bodyX+padX, t.Top, Cell{t.BodyColor, r})
		}
		f.Set(t.Left+1+bodyWidth, t.Top, Cell{t.FrameColor, '┃'})
		t.Top++
	}

	// Draw the footer, aligned as per
	f.Set(t.Left, t.Top, Cell{t.FrameColor, '┗'})
	footerStart := bodyWidth - len(t.Footer)
	footerEnd := bodyWidth
	if t.FooterLeftAlign {
		footerStart = 0
		footerEnd = len(t.Footer)
	}
	for dy := 0; dy < footerStart; dy++ {
		f.Set(t.Left+dy+1, t.Top, Cell{t.FrameColor, '━'})
	}
	for dy := footerStart; dy < footerEnd; dy++ {
		f.Set(t.Left+dy+1, t.Top, Cell{t.FooterColor, t.Footer[dy-footerStart]})
	}
	for dy := footerEnd; dy < bodyWidth; dy++ {
		f.Set(t.Left+dy+1, t.Top, Cell{t.FrameColor, '━'})
	}
	f.Set(t.Left+bodyWidth+1, t.Top, Cell{t.FrameColor, '┛'})
	t.Top++
}

func Rect(frames []*Framebuffer) image.Rectangle {
	x, y := 0, 0
	for _, frame := range frames {
		r := frame.Rect()
		if r.Dx() > x {
			x = r.Dx()
		}
		if r.Dy() > y {
			y = r.Dy()
		}
	}
	return image.Rect(0, 0, x, y)
}
