package main

import (
	"bytes"
	"fmt"
	"image/color"
	"image/gif"
	"os"
	"strconv"
	"strings"

	"github.com/asymmetricia/aoc22/aoc"
	"github.com/asymmetricia/aoc22/framebuffer"
	"github.com/asymmetricia/aoc22/term"
	"github.com/sirupsen/logrus"
)

type Throw int

const (
	Rock     Throw = 1
	Paper          = 2
	Scissors       = 3
)

func (t Throw) String() string {
	switch t {
	case Rock:
		return "rock"
	case Paper:
		return "paper"
	case Scissors:
		return "scissors"
	}
	return "(invalid)"
}

type Outcome int

const (
	Lose Outcome = 0
	Draw         = 3
	Win          = 6
)

func (o Outcome) String() string {
	switch o {
	case Lose:
		return "lose"
	case Draw:
		return "draw"
	case Win:
		return "win"
	}
	return "(invalid)"
}

type Game map[Outcome]uint
type World map[Throw]Game

func (w World) Record(o Outcome, t Throw) {
	if _, ok := w[t]; ok {
		w[t][o]++
	} else {
		w[t] = Game{o: 1}
	}
}

func (w World) Score() uint {
	var ret uint
	for t, game := range w {
		for o, count := range game {
			ret += (uint(t) + uint(o)) * count
		}
	}
	return ret
}

func frame(w World, o Outcome, t Throw, n int, total int, final bool) *framebuffer.Framebuffer {
	var ret = &framebuffer.Framebuffer{}
	for y, outcome := range []Outcome{Lose, Draw, Win} {
		for x, throw := range []Throw{Rock, Paper, Scissors} {
			frameColor := aoc.TolVibrantRed
			if outcome == Draw {
				frameColor = aoc.TolVibrantOrange
			} else if outcome == Win {
				frameColor = aoc.TolVibrantTeal
			}
			var bodyColor color.Color = aoc.TolVibrantMagenta
			if throw == Paper {
				bodyColor = color.White
			}
			if throw == Scissors {
				bodyColor = aoc.TolVibrantCyan
			}
			amt := w[throw][outcome]
			framebuffer.TextBox{
				Top:        y * (aoc.LineHeight),
				Left:       x * (aoc.GlyphWidth*3 + 2),
				Title:      []rune(throw.String()),
				Body:       []rune(fmt.Sprintf("%3d", amt)),
				BodyBlock:  true,
				Footer:     []rune(outcome.String()),
				FrameColor: frameColor,
				BodyColor:  bodyColor,
			}.On(ret)
		}
	}

	if final {
		framebuffer.TextBox{
			Top:       3 * aoc.LineHeight,
			Center:    true,
			Title:     []rune("Final Score"),
			Body:      []rune(strconv.Itoa(int(w.Score()))),
			BodyBlock: true,
		}.On(ret)
	} else {
		framebuffer.TextBox{
			Top:     3 * aoc.LineHeight,
			Center:  true,
			Footer:  []rune(fmt.Sprintf("game %d of %d", n, total)),
			Body:    []rune(fmt.Sprintf("Throw %8s in order to %4s", t.String(), o.String())),
			BodyPad: true,
		}.On(ret)
		framebuffer.TextBox{
			Top:    3*aoc.LineHeight + 5,
			Center: true,
			Title:  []rune("Score"),
			Body:   []rune(fmt.Sprintf(" %5d ", w.Score())),
		}.On(ret)
	}
	return ret
}

var log = logrus.StandardLogger()

func outcome(game string) (Outcome, Throw) {
	switch game {
	case "A X":
		return Lose, Scissors
	case "A Y":
		return Draw, Rock
	case "A Z":
		return Win, Paper
	case "B X":
		return Lose, Rock
	case "B Y":
		return Draw, Paper
	case "B Z":
		return Win, Scissors
	case "C X":
		return Lose, Paper
	case "C Y":
		return Draw, Scissors
	case "C Z":
		return Win, Rock
	}
	panic("bad game instructions " + game)
}

func solution(inputName string, input []byte) int {
	input = bytes.TrimSpace(input)
	input = bytes.Replace(input, []byte("\r"), []byte(""), -1)
	lines := strings.Split(strings.TrimSpace(string(input)), "\n")
	log.Printf("read %d input lines", len(lines))
	var score int
	world := World{}

	var frames []*framebuffer.Framebuffer
	for i, line := range lines {
		outcome, throw := outcome(line)
		world.Record(outcome, throw)
		score += int(outcome) + int(throw)
		if i < 100 ||
			i >= 100 && i < 200 && i%2 == 0 ||
			i >= 200 && i < 300 && i%4 == 0 ||
			i >= 300 && i < 400 && i%8 == 0 ||
			i >= 400 && i < 1000 && i%11 == 0 ||
			i >= 1000 && i%47 == 0 {
			frames = append(frames, frame(world, outcome, throw, i, len(lines), false))
		}
	}

	// doubled last frames improve the experience when converting to video
	frames = append(frames, frame(world, 0, 0, 0, 0, true))
	frames = append(frames, frame(world, 0, 0, 0, 0, true))

	anim := &gif.GIF{}
	delay := 100
	framesRect := framebuffer.Rect(frames)
	for i, frame := range frames {
		anim.Image = append(anim.Image, frame.TypeSetRect(framesRect, aoc.TypesetOpts{Scale: 2}))
		anim.Disposal = append(anim.Disposal, gif.DisposalNone)
		anim.Delay = append(anim.Delay, delay)
		delay = delay * 3 / 4
		if delay < 3 {
			delay = 3
		}
		term.ClearLine()
		fmt.Printf("\r%d/%d", i+1, len(frames))
	}
	println()

	anim.Delay[len(anim.Delay)-2] = 300

	aoc.SaveGIF(anim, fmt.Sprintf("day2b-%s.gif", inputName), log)

	return score
}

func main() {
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
