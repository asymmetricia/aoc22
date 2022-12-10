package main

import (
	"bytes"
	"os"
	"strconv"
	"strings"
	"unicode"

	"github.com/sirupsen/logrus"

	"github.com/asymmetricia/aoc22/aoc"
	"github.com/asymmetricia/aoc22/canvas"
)

var log = logrus.StandardLogger()

type CPU struct {
	X           int
	Cycle       int
	SS          int
	Signal      int
	FrameBuffer [240]bool
	Frames      []*canvas.Canvas
}

func (c *CPU) Noop() {
	c.Cycle++
	c.Draw()
}

func (c *CPU) Addx(i int) {
	c.Cycle++
	c.Draw()
	c.Cycle++
	c.Draw()
	c.X += i
}

func (c *CPU) Draw() {
	if (c.Cycle+20)%40 == 0 {
		c.Signal += c.Cycle * c.X
	}

	cycle := c.Cycle - 1
	if cycle%40 >= c.X-1 && cycle%40 <= c.X+1 {
		c.FrameBuffer[cycle] = true
	}
	cnv := &canvas.Canvas{}
	canvas.TextBox{
		Title:      []rune("Elfosonic"),
		FrameColor: aoc.TolVibrantBlue,
		Width:      40,
		Height:     6,
	}.On(cnv)
	for y := 0; y < 6; y++ {
		for x := 0; x < 40; x++ {
			if c.FrameBuffer[40*y+x] {
				cnv.PrintAt(x+1, y+1, "#", aoc.TolVibrantTeal)
			} else {
				cnv.PrintAt(x+1, y+1, ".", aoc.TolVibrantGrey)
			}
		}
	}
	canvas.TextBox{
		Left:  42,
		Title: []rune("Cycle"),
		Body:  []rune(strconv.Itoa(c.Cycle)),
		Width: 6,
	}.On(cnv)
	canvas.TextBox{
		Top:   2,
		Left:  42,
		Title: []rune("X"),
		Body:  []rune(strconv.Itoa(c.X)),
		Width: 6,
	}.On(cnv)
	canvas.TextBox{
		Top:   4,
		Left:  42,
		Title: []rune("Signal"),
		Body:  []rune(strconv.Itoa(c.Signal)),
		Width: 6,
	}.On(cnv)
	cnv.PrintAt(42, 7, " ⓒElfTek", aoc.TolVibrantMagenta)
	c.Frames = append(c.Frames, cnv)
}

func solution(name string, input []byte) int {
	// trim trailing space only
	input = bytes.Replace(input, []byte("\r"), []byte(""), -1)
	input = bytes.TrimRightFunc(input, unicode.IsSpace)
	lines := strings.Split(strings.TrimRightFunc(string(input), unicode.IsSpace), "\n")
	log.Printf("read %d %s lines", len(lines), name)

	cpu := &CPU{X: 1}
	for _, line := range lines {
		parts := strings.Fields(line)
		switch parts[0] {
		case "noop":
			cpu.Noop()
		case "addx":
			immed, err := strconv.Atoi(parts[1])
			if err != nil {
				panic(err)
			}
			cpu.Addx(immed)
		}
	}

	for y := 0; y < 6; y++ {
		for x := 0; x < 40; x++ {
			if cpu.FrameBuffer[y*40+x] {
				print("#")
			} else {
				print(".")
			}
		}
		println()
	}

	canvas.RenderGif(cpu.Frames, nil, "day10-"+name+".gif", log)

	return cpu.SS
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
