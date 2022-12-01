package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"math/rand"
	"os"
	"sort"
	"strings"

	"github.com/asymmetricia/aoc21/aoc"
)

func main() {
	var glyphs []rune
	for r := range aoc.Glyphs {
		glyphs = append(glyphs, r)
	}
	sort.Slice(glyphs, func(i, j int) bool {
		return glyphs[i] < glyphs[j]
	})

	const message = "the quick brown fox jumps over the lazy dog at (1,1)."
	xdim := len(message)
	ydim := len(glyphs)/xdim + 5
	const scale = 4

	g := &gif.GIF{}

	for i := 0; i < 10; i++ {
		img := image.NewPaletted(image.Rect(0, 0, xdim*8*scale, ydim*aoc.LineHeight*scale), aoc.TolVibrant)
		draw.Draw(img, img.Bounds(), image.Black, image.Point{}, draw.Src)
		for i, r := range glyphs {
			aoc.Typeset(img, image.Pt(
				(i%xdim)*8*scale,
				(i/xdim)*aoc.LineHeight*scale,
			), string(r), aoc.TolVibrant[rand.Intn(len(aoc.TolVibrant)-4)+3], aoc.TypesetOpts{Scale: scale})
		}
		aoc.Typeset(img, image.Pt(0, (len(glyphs)/xdim+2)*aoc.LineHeight*scale), message, color.White, aoc.TypesetOpts{Scale: scale})
		aoc.Typeset(img, image.Pt(0, (len(glyphs)/xdim+3)*aoc.LineHeight*scale), strings.ToUpper(message), color.White, aoc.TypesetOpts{Scale: scale})

		aoc.Typeset(
			img,
			image.Pt(0, (len(glyphs)/xdim+4)*aoc.LineHeight*scale),
			strings.ToUpper(message),
			color.White,
			aoc.TypesetOpts{scale, true},
		)

		g.Image = append(g.Image, img)
		g.Delay = append(g.Delay, 50)
		g.Disposal = append(g.Disposal, gif.DisposalNone)
	}

	aoc.Optimize(g.Image)

	f, err := os.OpenFile("out.gif", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err == nil {
		err = gif.EncodeAll(f, g)
	}
	if err == nil {
		err = f.Sync()
	}
	if err == nil {
		err = f.Close()
	}
	if err != nil {
		panic(err)
	}
}
