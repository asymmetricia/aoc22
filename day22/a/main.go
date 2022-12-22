package main

import (
	"bytes"
	"os"
	"strconv"
	"strings"
	"unicode"

	"github.com/asymmetricia/aoc22/coord"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"

	"github.com/asymmetricia/aoc22/aoc"
)

var log = logrus.StandardLogger()

type step struct {
	Steps int
	Turn  rune
}

type position struct {
	pos    coord.Coord
	facing coord.Direction
}

func (s step) String() string {
	if s.Steps == 0 {
		return string(s.Turn)
	}
	return strconv.Itoa(s.Steps)
}

func (s step) Move(from position, world *coord.DenseWorld) position {
	if s.Steps == 0 {
		if s.Turn == 'R' {
			return position{from.pos, from.facing.CW(false)}
		} else if s.Turn == 'L' {
			return position{from.pos, from.facing.CCW(false)}
		}
		return from
	}

steps:
	for s.Steps > 0 {
		s.Steps--
		pos := from.pos.Move(from.facing)
		switch world.At(pos) {
		case '#':
			break steps
		case '.':
			from.pos = pos
		case -1:
			fallthrough
		case 0:
			fallthrough
		case ' ':
			switch from.facing {
			case coord.North:
				for y := len(*world) - 1; y > 0; y-- {
					switch world.At(coord.C(pos.X, y)) {
					case '.':
						from.pos = coord.C(pos.X, y)
						continue steps
					case '#':
						break steps
					}
				}
			case coord.South:
				for y := 0; y < len(*world); y++ {
					switch world.At(coord.C(pos.X, y)) {
					case '.':
						from.pos = coord.C(pos.X, y)
						continue steps
					case '#':
						break steps
					}
				}
			case coord.East:
				for x := 0; x < pos.X; x++ {
					switch world.At(coord.C(x, pos.Y)) {
					case '.':
						from.pos = coord.C(x, pos.Y)
						continue steps
					case '#':
						break steps
					}
				}
			case coord.West:
				row := (*world)[pos.Y]
				for x := len(row) - 1; x > pos.X; x-- {
					switch world.At(coord.C(x, pos.Y)) {
					case '.':
						from.pos = coord.C(x, pos.Y)
						continue steps
					case '#':
						break steps
					}
				}
			}
		}
	}

	return from
}

func solution(name string, input []byte) int {
	// trim trailing space only
	input = bytes.Replace(input, []byte("\r"), []byte(""), -1)
	input = bytes.TrimRightFunc(input, unicode.IsSpace)
	lines := strings.Split(strings.TrimRightFunc(string(input), unicode.IsSpace), "\n")
	uniq := map[string]bool{}
	for _, line := range lines {
		uniq[line] = true
	}
	log.Printf("read %d %s lines (%d unique)", len(lines), name, len(uniq))

	var carta []string
	var steps string
	for i, line := range lines {
		if line == "" {
			steps = lines[i+1]
			break
		}
		carta = append(carta, line)
	}

	world := coord.Load(carta, true).(*coord.DenseWorld)
	start := position{
		pos:    coord.C(slices.Index((*world)[0], '.'), 0),
		facing: coord.East,
	}

	var stepList []step
	var accum step
	for _, i := range steps {
		if i == 'R' {
			stepList = append(stepList, accum, step{Turn: 'R'})
			accum.Steps = 0
		} else if i == 'L' {
			stepList = append(stepList, accum, step{Turn: 'L'})
			accum.Steps = 0
		} else {
			accum.Steps = accum.Steps*10 + int(i-'0')
		}
	}
	stepList = append(stepList, accum)

	for i, step := range stepList {
		start = step.Move(start, world)
		log.Printf("after %d, %v", i, start)
	}

	log.Print(start)
	log.Print(stepList)

	value := map[coord.Direction]int{
		coord.East:      0,
		coord.SouthWest: 1,
		coord.West:      2,
		coord.North:     3,
	}

	return 1000*(start.pos.Y+1) + 4*(start.pos.X+1) + value[start.facing]
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

	input := aoc.Input(2022, 22)
	log.Printf("input solution: %d", solution("input", input))
}
