package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/png"
	"log"
	"math/rand"
	"os"
	"time"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"

	"github.com/asymmetricia/aoc22/isovox"
)

func main() {
	w := isovox.World{Voxels: map[isovox.Coord]*isovox.Voxel{}}

	colors := []color.Color{
		color.NRGBA{245, 169, 184, 0x75},
		color.RGBA{91, 206, 250, 255},
		color.White,
	}

	for x := -1; x < 6; x++ {
		w.Voxels[isovox.Coord{x, -1, -1}] = &isovox.Voxel{Color: colors[0]}
		w.Voxels[isovox.Coord{x, 5, -1}] = &isovox.Voxel{Color: colors[0]}
		w.Voxels[isovox.Coord{x, -1, 5}] = &isovox.Voxel{Color: colors[0]}
		w.Voxels[isovox.Coord{x, 5, 5}] = &isovox.Voxel{Color: colors[0]}
	}
	for y := -1; y < 6; y++ {
		w.Voxels[isovox.Coord{-1, y, -1}] = &isovox.Voxel{Color: colors[0]}
		w.Voxels[isovox.Coord{5, y, -1}] = &isovox.Voxel{Color: colors[0]}
		w.Voxels[isovox.Coord{-1, y, 5}] = &isovox.Voxel{Color: colors[0]}
		w.Voxels[isovox.Coord{5, y, 5}] = &isovox.Voxel{Color: colors[0]}
	}
	for z := -1; z < 6; z++ {
		w.Voxels[isovox.Coord{-1, -1, z}] = &isovox.Voxel{Color: colors[0]}
		w.Voxels[isovox.Coord{5, -1, z}] = &isovox.Voxel{Color: colors[0]}
		w.Voxels[isovox.Coord{-1, 5, z}] = &isovox.Voxel{Color: colors[0]}
		w.Voxels[isovox.Coord{5, 5, z}] = &isovox.Voxel{Color: colors[0]}
	}

	images := []image.Image{}
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 10; i++ {
		for x := 0; x < 5; x++ {
			for y := 0; y < 5; y++ {
				for z := 0; z < 5; z++ {
					w.Voxels[isovox.Coord{x, y, z}] = &isovox.Voxel{Color: colors[rand.Intn(len(colors))]}
				}
			}
		}
		img := w.Render(100)
		log.Printf("%d x %d", img.Bounds().Dx(), img.Bounds().Dy())
		images = append(images, img)
	}

	p := PaletteFrom(images)
	anim := &gif.GIF{
		Image:    make([]*image.Paletted, len(images)),
		Delay:    make([]int, len(images)),
		Disposal: make([]byte, len(images)),
	}

	for i, img := range images {
		pi := image.NewPaletted(img.Bounds(), p)
		draw.FloydSteinberg.Draw(pi, pi.Bounds(), img, image.Pt(0, 0))
		anim.Image[i] = pi
		anim.Delay[i] = 10
		anim.Disposal[i] = gif.DisposalNone
	}

	Optimize(anim.Image)

	f, err := os.OpenFile("out.gif", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err == nil {
		err = gif.EncodeAll(f, anim)
	}
	if err == nil {
		err = f.Sync()
	}
	if err == nil {
		err = f.Close()
	}

	if err == nil {
		f, err = os.OpenFile("out.png", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	}
	if err == nil {
		err = png.Encode(f, images[0])
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

func PaletteFrom(images []image.Image) color.Palette {
	ret := color.Palette{
		color.Black,
		color.Transparent,
		color.White,
	}

	colors := map[color.RGBA64]int{}
	for _, img := range images {
		r := img.Bounds()
		for x := r.Min.X; x <= r.Max.X; x++ {
			for y := r.Min.Y; y <= r.Max.Y; y++ {
				r, g, b, a := img.At(x, y).RGBA()
				colors[color.RGBA64{
					R: uint16(r),
					G: uint16(g),
					B: uint16(b),
					A: uint16(a),
				}]++
			}
		}
	}

	colorList := maps.Keys(colors)

	slices.SortFunc(colorList, func(a, b color.RGBA64) bool {
		return colors[b] < colors[a]
	})
	for len(ret) < 255 && len(colorList) > 0 {
		ret = append(ret, colorList[0])
		colorList = colorList[1:]
	}

	return ret
}

func Optimize(imgs []*image.Paletted) {
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
		perc := n * 100 / len(imgs)
		if perc%10 == 0 && perc > lp {
			lp = perc
			log.Printf("resizing %d%%...", perc)
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

		perc := n * 100 / len(imgs)
		if perc%10 == 0 && perc > lp {
			lp = perc
			log.Printf("optimizing %d%%...", perc)
		}
	}
}
