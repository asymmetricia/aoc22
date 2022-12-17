package main

import (
	"bytes"
	"os"
	"strings"
	"unicode"

	"github.com/sirupsen/logrus"

	"github.com/asymmetricia/aoc22/aoc"
	"github.com/asymmetricia/aoc22/coord"
)

var log = logrus.StandardLogger()

/*
####

.#.
###
.#.

..#
..#
###

#
#
#
#

##
##
*/

var symbols = [][]coord.Coord{
	{{0, 0}, {1, 0}, {2, 0}, {3, 0}},
	{{1, 0}, {0, -1}, {1, -1}, {2, -1}, {1, -2}},
	{{0, 0}, {1, 0}, {2, 0}, {2, -1}, {2, -2}},
	{{0, 0}, {0, -1}, {0, -2}, {0, -3}},
	{{0, 0}, {1, 0}, {0, -1}, {1, -1}},
}

func solution(name string, input []byte) int {
	// trim trailing space only
	input = bytes.Replace(input, []byte("\r"), []byte(""), -1)
	input = bytes.TrimRightFunc(input, unicode.IsSpace)
	lines := strings.Split(strings.TrimRightFunc(string(input), unicode.IsSpace), "\n")
	log.Printf("read %d %s lines", len(lines), name)
	log.Print(len(lines[0]))

	world := coord.SparseWorld{}

	for y := 0; y > -5; y-- {
		world.Set(coord.C(-1, y), '|')
		world.Set(coord.C(7, y), '|')
	}
	for x := -1; x <= 7; x++ {
		world.Set(coord.C(x, 1), '-')
	}

	clear := func(symbol []coord.Coord, pos coord.Coord) bool {
		for _, c := range symbol {
			c = c.Plus(pos)
			if c.X < 0 || c.X > 6 {
				return false
			}
			if c.Y > 0 {
				return false
			}
			if world.At(c) == '#' {
				return false
			}
		}
		return true
	}

	place := func(symbol []coord.Coord, pos coord.Coord, r rune) int {
		for _, c := range symbol {
			world.Set(c.Plus(pos), r)
		}
		var min int
		world.Each(func(c coord.Coord) (stop bool) {
			if world.At(c) == '#' && c.Y < min {
				min = c.Y
			}
			return false
		})
		return min
	}

	jidx := 0
	idx := 0
	height := 0
	// 15815
	// 31644
	for idx < 2022 {
		pos := coord.C(2, height-3)
		sym := symbols[idx%5]
		for {
			dir := lines[0][jidx]
			jidx++
			if jidx >= len(lines[0]) {
				jidx = 0
			}

			if dir == '>' {
				if clear(sym, pos.East()) {
					pos = pos.East()
				}
			} else {
				if clear(sym, pos.West()) {
					pos = pos.West()
				}
			}
			if clear(sym, pos.South()) {
				pos = pos.South()
			} else {
				break
			}

		}

		height = place(sym, pos, '#') - 1
		idx++
	}

	return -height
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

	input := aoc.Input(2022, 17)
	log.Printf("input solution: %d", solution("input", input))
}
