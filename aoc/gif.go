package aoc

import (
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"os"

	"github.com/sirupsen/logrus"
)

func Optimize(imgs []*image.Paletted, log ...logrus.FieldLogger) {
	if len(imgs) < 2 {
		return
	}

	var x, y int
	for _, img := range imgs {
		ix := img.Rect.Dx()
		iy := img.Rect.Dy()
		if ix > x {
			x = ix
		}
		if iy > y {
			y = iy
		}
	}

	lp := 0
	for n, img := range imgs {
		ix := img.Rect.Dx()
		iy := img.Rect.Dy()
		if ix != x || iy != y {
			repl := image.NewPaletted(image.Rect(0, 0, x, y), img.Palette)
			draw.Draw(repl, img.Rect, img, image.Point{}, draw.Over)
			*img = *repl
		}
		if len(log) > 0 {
			perc := n * 100 / len(imgs)
			if perc%10 == 0 && perc > lp {
				lp = perc
				log[0].Printf("resizing %d%%...", perc)
			}
		}
	}

	accum := image.NewPaletted(image.Rect(0, 0, x, y), imgs[0].Palette)
	copy(accum.Pix, imgs[0].Pix)

	tr := imgs[0].Palette.Index(color.Transparent)

	lp = 0
	for n, img := range imgs[1:] {
		for i, v := range img.Pix {
			if v == accum.Pix[i] {
				img.Pix[i] = uint8(tr)
			} else {
				accum.Pix[i] = img.Pix[i]
			}
		}

		if len(log) > 0 {
			perc := n * 100 / len(imgs)
			if perc%10 == 0 && perc > lp {
				lp = perc
				log[0].Printf("optimizing %d%%...", perc)
			}
		}
	}
}

func SaveGIF(g *gif.GIF, filename string, log ...logrus.FieldLogger) {
	Optimize(g.Image, log...)
	f, err := os.OpenFile(filename, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if len(log) > 0 {
		log[0].Print("encoding GIF...")
	}
	if err == nil {
		err = gif.EncodeAll(f, g)
	}
	if len(log) > 0 {
		log[0].Print("done!")
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
