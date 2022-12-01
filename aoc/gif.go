package aoc

import (
	"image"
	"image/color"
	"image/gif"
	"os"
)

func Optimize(imgs []*image.Paletted) {
	if len(imgs) < 2 {
		return
	}
	accum := image.NewPaletted(imgs[0].Rect, imgs[0].Palette)
	copy(accum.Pix, imgs[0].Pix)

	tr := imgs[0].Palette.Index(color.Transparent)
	for _, img := range imgs[1:] {
		for i, v := range img.Pix {
			if v == accum.Pix[i] {
				img.Pix[i] = uint8(tr)
			} else {
				accum.Pix[i] = img.Pix[i]
			}
		}
	}
}

func SaveGIF(g *gif.GIF, filename string) {
	Optimize(g.Image)
	f, err := os.OpenFile(filename, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
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
