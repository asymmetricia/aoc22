package aoc

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/png"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"

	"github.com/sirupsen/logrus"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
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
	WriteGIF(g, filename, log...)
}

func WriteGIF(g *gif.GIF, filename string, log ...logrus.FieldLogger) {
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
	log.Print(len(colorList))
	slices.SortFunc(colorList, func(a, b color.RGBA64) bool {
		return colors[b] < colors[a]
	})
	for len(ret) < 255 && len(colorList) > 0 {
		ret = append(ret, colorList[0])
		colorList = colorList[1:]
	}

	return ret
}

func RenderPng(image image.Image, filename string) error {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	if err := png.Encode(f, image); err != nil {
		f.Close()
		return err
	}

	if err := f.Sync(); err != nil {
		f.Close()
		return err
	}

	return f.Close()
}

type MP4Encoder struct {
	cmd   *exec.Cmd
	stdin io.WriteCloser
}

func NewMP4Encoder(filename string, framerate int, log ...logrus.FieldLogger) (*MP4Encoder, error) {
	ffmpeg, err := exec.LookPath("ffmpeg")
	if err != nil {
		return nil, fmt.Errorf("could not find ffmpeg in path: %w", err)
	}

	cmd := exec.Command(ffmpeg,
		"-loglevel", "info",
		"-stats",
		"-y",
		"-hide_banner",
		"-f", "image2pipe",
		"-c:v", "png",
		"-framerate", strconv.Itoa(framerate),
		"-i", "-",
		"-c:v", "libx264",
		"-movflags", "+faststart",
		"-preset", "veryslow",
		"-tune", "animation",
		"-threads", "0",
		"-crf", "22",
		"-f", "mp4",
		"-vf", "format=yuv420p",
		filename,
	)

	in, _ := cmd.StdinPipe()

	if len(log) > 0 {
		outPipe, _ := cmd.StdoutPipe()
		errPipe, _ := cmd.StderrPipe()

		go func(r io.Reader) {
			sc := bufio.NewScanner(r)
			for sc.Scan() {
				log[0].Printf("ffmpeg: %s", sc.Text())
			}
		}(outPipe)

		go func(r io.Reader) {
			sc := bufio.NewScanner(r)
			for sc.Scan() {
				log[0].Errorf("ffmpeg: %s", sc.Text())
			}
		}(errPipe)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("could not start ffmpeg: %w", err)
	}

	return &MP4Encoder{
		cmd:   cmd,
		stdin: in,
	}, nil
}

func (enc *MP4Encoder) Encode(image image.Image) error {
	err := png.Encode(enc.stdin, image)
	if err != nil {
		return err
	}
	if f, ok := enc.stdin.(*os.File); ok {
		return f.Sync()
	}
	return nil
}

func (enc *MP4Encoder) Close() error {
	if err := enc.stdin.Close(); err != nil {
		return err
	}
	return enc.cmd.Wait()
}
