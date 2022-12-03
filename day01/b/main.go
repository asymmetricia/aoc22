package main

import (
	"image/gif"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/asymmetricia/aoc22/aoc"
	"github.com/asymmetricia/aoc22/canvas"
	"github.com/asymmetricia/aoc22/set"
	"github.com/asymmetricia/aoc22/term"
)

func frame(elves [][]int, selected []int, totals []int, cursor int) canvas.Canvas {
	width := len(strconv.Itoa(aoc.MaxFn(elves, func(elf []int) int { return aoc.Max(elf...) })))
	height := aoc.MaxFn(elves, func(elf []int) int { return len(elf) })
	if len(totals) > 0 {
		height++
	}

	selSet := set.FromItems(selected)

	var ret canvas.Canvas
	const bw = 4
	x := bw - 1
	y := 0
	size := 1
	for i, elf := range elves {
		for snackIdx, snack := range elf {
			color := aoc.TolVibrantGrey
			if selSet[i] {
				color = aoc.TolVibrantCyan
			} else if cursor == i {
				color = aoc.TolVibrantMagenta
			}
			ret.PrintAt(x*(width+2), y*(height+1)+snackIdx, strconv.Itoa(snack), color)
		}
		if len(totals) > i {
			ret.PrintAt(x*(width+2), y*(height+1)+len(elf), strconv.Itoa(totals[i]), aoc.TolVibrantRed)
		}

		if x < ((size-1)*bw + 1) {
			if y == size-1 {
				if x == 0 {
					size++
					y = 0
					x = size*bw - 1
				} else {
					x--
				}
			} else {
				y++
				x = size*bw - 1
			}
		} else {
			x--
		}
	}

	return ret
}

func main() {
	in, err := ioutil.ReadFile("input")
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(in), "\n")

	var frames []canvas.Canvas

	var elves [][]int
	var accum []int
	for _, line := range lines {
		if line == "" && len(accum) > 0 {
			elves = append(elves, accum)
			accum = nil
			frames = append(frames, frame(elves, nil, nil, len(elves)-1))
		} else {
			i, err := strconv.Atoi(line)
			if err != nil {
				panic(err)
			}
			accum = append(accum, i)
		}
	}
	if len(accum) > 0 {
		elves = append(elves, accum)
	}

	var counts []int
	for i, calories := range elves {
		total := 0
		for _, c := range calories {
			total += c
		}
		counts = append(counts, total)
		frames = append(frames, frame(elves, nil, counts, i))
	}

	var first, second, third int
	first = 0
	if counts[1] > counts[0] {
		first = 1
		second = 0
	} else {
		second = 1
	}
	if counts[2] > counts[first] {
		first, second, third = 2, first, second
	} else if counts[2] > counts[second] {
		second, third = 2, second
	}

	for i, count := range counts {
		if count > counts[first] {
			first, second, third = i, first, second
		} else if count > counts[second] {
			second, third = i, second
		} else if count > counts[third] {
			third = i
		}
		frames = append(frames, frame(elves, []int{first, second, third}, counts, i))
	}

	soln := counts[first] + counts[second] + counts[third]
	log.Print(soln)

	txt := strconv.Itoa(soln)
	message := "┏━"
	for i := 0; i < len(txt); i++ {
		message += "━"
	}
	message += "━┓\n┃ " + txt + " ┃\n┗━"
	for i := 0; i < len(txt); i++ {
		message += "━"
	}
	message += "━┛"
	lastFrame := frames[len(frames)-1].Copy()
	rect := lastFrame.Rect()
	for y, line := range strings.Split(message, "\n") {
		lastFrame.BlockPrintAt(
			rect.Dx()/2-utf8.RuneCountInString(line)/2*aoc.GlyphWidth,
			rect.Dy()/2-aoc.LineHeight+y*aoc.LineHeight,
			line,
			aoc.TolVibrantTeal)
	}
	frames = append(frames, lastFrame)

	anim := &gif.GIF{}
	for i, frame := range frames {
		anim.Image = append(anim.Image, frame.Render())
		anim.Delay = append(anim.Delay, 3)
		anim.Disposal = append(anim.Disposal, gif.DisposalNone)
		print("\r")
		term.ClearLine()
		print(i+1, "/", len(frames))
	}
	println()

	anim.Delay[len(anim.Delay)-1] = 600
	aoc.SaveGIF(anim, "day1b.gif")
}
