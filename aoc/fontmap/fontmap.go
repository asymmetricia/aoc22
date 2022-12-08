package main

import (
	"image/color"
	"image/gif"
	"math/rand"
	"os"
	"sort"
	"strings"

	"github.com/asymmetricia/aoc22/aoc"
	"github.com/asymmetricia/aoc22/canvas"
)

func main() {
	var glyphs []rune
	for r := range aoc.Glyphs[aoc.Pixl] {
		glyphs = append(glyphs, r)
	}
	sort.Slice(glyphs, func(i, j int) bool {
		return glyphs[i] < glyphs[j]
	})

	const message = "the quick brown fox jumps over the lazy dog"
	xdim := len(message)
	ydim := len(glyphs)/xdim + 4
	const scale = 3

	g := &gif.GIF{}

	for i := 0; i < 10; i++ {
		frame := &canvas.Canvas{}

		canvas.TextBox{
			Title:  []rune("'Pixl' Font Demo"),
			Width:  xdim,
			Height: ydim,
		}.On(frame)
		for i, r := range glyphs {
			frame.PrintAt(
				1+i%xdim,
				1+i/xdim,
				string(r),
				aoc.TolVibrant[rand.Intn(len(aoc.TolVibrant)-4)+3],
			)
		}
		frame.PrintAt(
			1, ydim-1,
			message,
			color.White)
		frame.PrintAt(
			1, ydim,
			strings.ToUpper(message),
			color.White)

		g.Image = append(g.Image, frame.Render(aoc.TypesetOpts{Scale: scale}))
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
