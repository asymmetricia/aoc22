package main

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"

	"github.com/asymmetricia/aoc22/aoc"
	"github.com/asymmetricia/aoc22/coord"
	"github.com/asymmetricia/aoc22/term"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
)

var log = logrus.StandardLogger()

type side int

func (s side) String() string {
	return map[side]string{
		Top:    "top",
		North:  "north",
		East:   "east",
		West:   "west",
		South:  "south",
		Bottom: "bottom",
	}[s]
}

const (
	Top side = iota
	North
	East
	South
	West
	Bottom
)

type step struct {
	Steps int
	Turn  rune
}

type position struct {
	side   side
	pos    coord.Coord
	facing coord.Direction
}

func (p position) String() string {
	return fmt.Sprintf("at %s on side %s facing %s", p.pos.String(), p.side.String(), p.facing.String())
}

func (s step) String() string {
	if s.Steps == 0 {
		return string(s.Turn)
	}
	return strconv.Itoa(s.Steps)
}

type rotation int

const (
	same rotation = iota
	right
	left
	flip
)

type corner struct {
	side side
	edge coord.Direction
	rot  rotation
}

func (c corner) rotate(d coord.Direction) coord.Direction {
	switch c.rot {
	case same:
		return d
	case flip:
		return d.CW(false).CW(false)
	case right:
		return d.CW(false)
	case left:
		return d.CCW(false)
	default:
		panic(c.rot)
	}
}

type transform func(position) position

var mtx = map[side]map[coord.Direction]transform{
	Top: {
		coord.North: func(p position) position { return position{North, coord.C(0, p.pos.X), coord.East} },
		coord.South: func(p position) position { return position{South, coord.C(p.pos.X, 0), coord.South} },
		coord.West:  func(p position) position { return position{West, coord.C(0, 49-p.pos.Y), coord.East} },
		coord.East:  func(p position) position { return position{East, coord.C(0, p.pos.Y), coord.East} },
	},
	North: {
		coord.North: func(p position) position { return position{West, coord.C(p.pos.X, 49), coord.North} },
		coord.East:  func(p position) position { return position{Bottom, coord.C(p.pos.Y, 49), coord.North} },
		coord.West:  func(p position) position { return position{Top, coord.C(p.pos.Y, 0), coord.South} },
		coord.South: func(p position) position { return position{East, coord.C(p.pos.X, 0), coord.South} },
	},
	East: {
		coord.North: func(p position) position { return position{North, coord.C(p.pos.X, 49), coord.North} },
		coord.East:  func(p position) position { return position{Bottom, coord.C(49, 49-p.pos.Y), coord.West} },
		coord.West:  func(p position) position { return position{Top, coord.C(49, p.pos.Y), coord.West} },
		coord.South: func(p position) position { return position{South, coord.C(49, p.pos.X), coord.West} },
	},
	West: {
		coord.North: func(p position) position { return position{South, coord.C(0, p.pos.X), coord.East} },
		coord.East:  func(p position) position { return position{Bottom, coord.C(0, p.pos.Y), coord.East} },
		coord.West:  func(p position) position { return position{Top, coord.C(0, 49-p.pos.Y), coord.East} },
		coord.South: func(p position) position { return position{North, coord.C(p.pos.X, 0), coord.South} },
	},
	South: {
		coord.North: func(p position) position { return position{Top, coord.C(p.pos.X, 49), coord.North} },
		coord.East:  func(p position) position { return position{East, coord.C(p.pos.Y, 49), coord.North} },
		coord.West:  func(p position) position { return position{West, coord.C(p.pos.Y, 0), coord.South} },
		coord.South: func(p position) position { return position{Bottom, coord.C(p.pos.X, 0), coord.South} },
	},
	Bottom: {
		coord.North: func(p position) position { return position{South, coord.C(p.pos.X, 49), coord.North} },
		coord.East:  func(p position) position { return position{East, coord.C(49, 49-p.pos.Y), coord.West} },
		coord.West:  func(p position) position { return position{West, coord.C(49, p.pos.Y), coord.West} },
		coord.South: func(p position) position { return position{North, coord.C(49, p.pos.X), coord.West} },
	},
}

func normalize(p position) position {
	if p.pos.X >= 0 && p.pos.X <= 49 &&
		p.pos.Y >= 0 && p.pos.Y <= 49 {
		return p
	}

	return mtx[p.side][p.facing](p)
}

const debug = true

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

	dim := 50
	layout := map[side]coord.Coord{
		Top:    coord.C(50, 0),
		East:   coord.C(100, 0),
		South:  coord.C(50, 50),
		West:   coord.C(0, 100),
		Bottom: coord.C(50, 100),
		North:  coord.C(0, 150),
	}

	maps := map[side]*coord.DenseWorld{}

	for side, tl := range layout {
		var sidelines []string
		for _, line := range lines[tl.Y : tl.Y+dim] {
			sidelines = append(sidelines, line[tl.X:tl.X+dim])
		}
		maps[side] = coord.Load(sidelines, true).(*coord.DenseWorld)
	}

	globalFromMap := func(s side, c coord.Coord) coord.Coord {
		return c.Plus(layout[s])
	}

	var steps = lines[201]
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

	startX := slices.Index((*maps[Top])[0], '.')

	pos := position{
		side:   Top,
		pos:    coord.C(startX, 0),
		facing: coord.East,
	}

	for _, step := range stepList {
		log.Print(pos)
		switch step.Turn {
		case 'R':
			pos.facing = pos.facing.CW(false)
			continue
		case 'L':
			pos.facing = pos.facing.CCW(false)
			continue
		case 0:
		default:
			log.Fatalf("%c", step.Turn)
		}

	stepCount:
		for step.Steps > 0 {
			if debug && (pos.pos.X < 1 || pos.pos.X > 48 || pos.pos.Y < 1 || pos.pos.Y > 48) {
				term.Clear()
				term.MoveCursor(1, 1)
				println(pos.side.String())
				f := rune(strings.ToUpper(pos.facing.String())[0])
				maps[pos.side].Set(pos.pos, f)
				for side, tl := range map[side][2]int{
					Top:    {51, 2},
					North:  {1, 152},
					East:   {101, 2},
					West:   {1, 102},
					South:  {51, 52},
					Bottom: {51, 102},
				} {
					for i, line := range strings.Split(maps[side].String(), "\n") {
						term.MoveCursor(tl[0]*2, tl[1]+i)
						for _, r := range line {
							if r == f {
								print(term.Scolor(0, 255, 255) + string(f) + string(f) + term.ScolorReset())
							} else {
								print(string(r), string(r))
							}
						}
					}
				}
				maps[pos.side].Set(pos.pos, '.')
				b := make([]byte, 1)
				os.Stdin.Read(b)
			}

			step.Steps--
			nextpos := normalize(position{
				side:   pos.side,
				pos:    pos.pos.Move(pos.facing),
				facing: pos.facing,
			})
			switch maps[nextpos.side].At(nextpos.pos) {
			case '.':
				pos = nextpos
				continue
			case '#':
				break stepCount
			default:
				log.Fatalf("bad position %v", nextpos)
			}
		}
	}

	value := map[coord.Direction]int{
		coord.East:      0,
		coord.SouthWest: 1,
		coord.West:      2,
		coord.North:     3,
	}

	pp := globalFromMap(pos.side, pos.pos)

	score := 1000*(pp.Y+1) + 4*(pp.X+1) + value[pos.facing]
	if score >= 197036 {
		panic("nope")
	}
	if score <= 134060 {
		panic("nope")
	}

	return score
	//return -1
	//world := coord.Load(carta, false).(*coord.SparseWorld)
	//start := position{
	//	pos:    coord.C(slices.Index((*world)[0], '.'), 0),
	//	facing: coord.East,
	//}
	//
	//var stepList []step
	//var accum step
	//for _, i := range steps {
	//	if i == 'R' {
	//		stepList = append(stepList, accum, step{Turn: 'R'})
	//		accum.Steps = 0
	//	} else if i == 'L' {
	//		stepList = append(stepList, accum, step{Turn: 'L'})
	//		accum.Steps = 0
	//	} else {
	//		accum.Steps = accum.Steps*10 + int(i-'0')
	//	}
	//}
	//stepList = append(stepList, accum)
	//
	//for i, step := range stepList {
	//	start = step.Move(start, world)
	//	log.Printf("after %d, %v", i, start)
	//}
	//
	//log.Print(start)
	//log.Print(stepList)
	//
	//value := map[coord.Direction]int{
	//	coord.East:      0,
	//	coord.SouthWest: 1,
	//	coord.West:      2,
	//	coord.North:     3,
	//}
	//
	//return 1000*(start.pos.Y+1) + 4*(start.pos.X+1) + value[start.facing]
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
