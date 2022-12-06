package main

import (
	"bytes"
	"fmt"
	"image/gif"
	"os"
	"strings"
	"unicode"

	"github.com/sirupsen/logrus"

	"github.com/asymmetricia/aoc22/aoc"
	"github.com/asymmetricia/aoc22/canvas"
)

var log = logrus.StandardLogger()

func frame(s []rune, i int, n int, match bool) *canvas.Canvas {
	c := aoc.TolVibrantRed
	if match {
		c = aoc.TolVibrantCyan
	}

	ret := &canvas.Canvas{}
	canvas.TextBox{
		Top:         0,
		Left:        0,
		Title:       []rune("Scanning for Start of Message"),
		Footer:      []rune(fmt.Sprintf("characters %d..%d", i, i+n)),
		TitleColor:  aoc.TolVibrantMagenta,
		FrameColor:  aoc.TolVibrantTeal,
		FooterColor: aoc.TolVibrantMagenta,
		Width:       128,
		Height:      64,
	}.On(ret)
outer:
	for y := 0; y < 64; y++ {
		for x := 0; x < 64; x++ {
			idx := y*64 + x
			posX := x*2 + 2
			posY := y + 1
			if idx >= len(s) {
				break outer
			}
			if idx == i {
				ret.Set(posX-1, posY, canvas.Cell{Value: '[', Color: c})
			}
			if idx == i+n {
				ret.Set(posX-1, posY, canvas.Cell{Value: ']', Color: c})
			}
			if idx >= i && idx <= i+n {
				ret.Set(posX, posY, canvas.Cell{Value: s[idx], Color: c})
			} else {
				ret.Set(posX, posY, canvas.Cell{Value: s[idx], Color: aoc.TolVibrantGrey})
			}
		}
	}

	if match {
		canvas.TextBox{
			Top:        3,
			Center:     true,
			Title:      []rune("Start Of Message Found!"),
			Body:       s[i : i+n],
			BodyBlock:  true,
			BodyColor:  aoc.TolVibrantCyan,
			TitleColor: aoc.TolVibrantTeal,
			FrameColor: aoc.TolVibrantMagenta,
		}.On(ret)
	}

	canvas.TextBox{
		Top:        7 + aoc.LineHeight,
		Center:     true,
		Title:      []rune("Message Starts After"),
		Body:       []rune(fmt.Sprintf("%04d", i+n)),
		BodyBlock:  true,
		BodyColor:  aoc.TolVibrantCyan,
		TitleColor: aoc.TolVibrantTeal,
		FrameColor: aoc.TolVibrantMagenta,
		BodyFont:   aoc.SevenSeg,
		TitleFont:  aoc.Pixl,
	}.On(ret)

	return ret
}

func isMarker(s string) bool {
	if len(s) != 14 {
		panic(s)
	}
	seen := map[byte]bool{}
	for _, c := range []byte(s) {
		if seen[c] {
			return false
		}
		seen[c] = true
	}
	return true
}

func solution(name string, input []byte) int {
	// trim trailing space only
	input = bytes.Replace(input, []byte("\r"), []byte(""), -1)
	input = bytes.TrimRightFunc(input, unicode.IsSpace)
	lines := strings.Split(strings.TrimRightFunc(string(input), unicode.IsSpace), "\n")
	log.Printf("read %d %s lines", len(lines), name)

	var frames []*canvas.Canvas
	for i := 0; i < len(input); i++ {
		match := isMarker(string(input[i : i+14]))
		if match || i < 10 ||
			i < 100 && i%2 == 0 ||
			i < 200 && i%3 == 0 ||
			i < 500 && i%5 == 0 ||
			i%10 == 0 {
			frame := frame([]rune(string(input)), i, 14, match)
			frames = append(frames, frame)
		}
		if !match {
			continue
		}

		log.Printf("%d, %s", i+14, input[i:i+14])
		break
	}

	anim := &gif.GIF{}
	for _, frame := range frames {
		anim.Image = append(anim.Image, frame.Render())
		anim.Delay = append(anim.Delay, 3)
		anim.Disposal = append(anim.Disposal, gif.DisposalNone)
	}
	anim.Delay[len(anim.Delay)-1] = 1000

	// Stick an extra frame on the end b/c of annoying transcoders
	anim.Image = append(anim.Image, frames[len(frames)-1].Render())
	anim.Delay = append(anim.Delay, 3)
	anim.Disposal = append(anim.Disposal, gif.DisposalNone)

	aoc.SaveGIF(anim, "day06b-"+name+".gif", log)

	return -1
}

func main() {
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02T15:04:05",
	})
	test, err := os.ReadFile("test")
	if err == nil {
		log.Printf("test solution: %d", solution("test", test))
	} else {
		log.Warningf("no test data present")
	}

	input, err := os.ReadFile("input")
	if err != nil {
		log.WithError(err).Fatal("could not read input")
	}
	log.Printf("input solution: %d", solution("input", input))
}
