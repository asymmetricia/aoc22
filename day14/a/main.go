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

func solution(name string, input []byte) int {
	// trim trailing space only
	input = bytes.Replace(input, []byte("\r"), []byte(""), -1)
	input = bytes.TrimRightFunc(input, unicode.IsSpace)
	lines := strings.Split(strings.TrimRightFunc(string(input), unicode.IsSpace), "\n")
	log.Printf("read %d %s lines", len(lines), name)

	world := &coord.DenseWorld{}
	for _, line := range lines {
		fields := strings.Split(line, " -> ")
		var prev coord.Coord
		for i, field := range fields {
			c := coord.MustFromComma(field)
			if i == 0 {
				prev = c
				world.Set(prev, '#')
				continue
			}
			for prev.X < c.X {
				prev.X++
				world.Set(prev, '#')
			}
			for prev.X > c.X {
				prev.X--
				world.Set(prev, '#')
			}
			for prev.Y < c.Y {
				prev.Y++
				world.Set(prev, '#')
			}
			for prev.Y > c.Y {
				prev.Y--
				world.Set(prev, '#')
			}
		}
	}

	var sandStart = coord.MustFromComma("500,0")
	world.Set(sandStart, 'v')

	particles := 0
	pos := sandStart
	for {
		if pos.South().Y >= len(*world) {
			break
		}
		if world.At(pos.South()) == 0 {
			pos = pos.South()
			continue
		}
		if world.At(pos.SouthWest()) == 0 {
			pos = pos.SouthWest()
			continue
		}
		if world.At(pos.SouthEast()) == 0 {
			pos = pos.SouthEast()
			continue
		}
		world.Set(pos, '+')
		particles++
		pos = sandStart
	}

	return particles
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

	input := aoc.Input(2022, 14)
	log.Printf("input solution: %d", solution("input", input))
}
