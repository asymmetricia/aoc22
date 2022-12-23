package main

import (
	"bytes"
	"image"
	"image/draw"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/asymmetricia/pencil"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"

	"github.com/asymmetricia/aoc22/aoc"
	"github.com/asymmetricia/aoc22/canvas"
	"github.com/asymmetricia/aoc22/coord"
	"github.com/asymmetricia/aoc22/term"
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

func (s step) String() string {
	if s.Steps == 0 {
		return string(s.Turn)
	}
	return strconv.Itoa(s.Steps)
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

const debug = false
const video = false

var layout = map[side]coord.Coord{
	Top:    coord.C(50, 0),   // 50..99   & 0..49
	East:   coord.C(100, 0),  // 100..149 & 0..49
	South:  coord.C(50, 50),  // 50..99   & 50..99
	West:   coord.C(0, 100),  // 0..49    & 100..149
	Bottom: coord.C(50, 100), // 50..99   & 100..149
	North:  coord.C(0, 150),  // 0..49    & 150..199
}

func globalFromMap(s side, c coord.Coord) coord.Coord {
	return c.Plus(layout[s])
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

	dim := 50

	maps := map[side]*coord.DenseWorld{}

	for side, tl := range layout {
		var sidelines []string
		for _, line := range lines[tl.Y : tl.Y+dim] {
			sidelines = append(sidelines, line[tl.X:tl.X+dim])
		}
		maps[side] = coord.Load(sidelines, true).(*coord.DenseWorld)
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
	stepList = append(stepList, accum)

	startX := slices.Index((*maps[Top])[0], '.')

	pos := position{
		side:   Top,
		pos:    coord.C(startX, 0),
		facing: coord.East,
	}

	var last draw.Image
	var enc *aoc.MP4Encoder

	if video {
		frame := canvas.Canvas{}
		frame.PrintAt(0, 0, strings.Join(lines[:200], "\n"), aoc.TolVibrantGrey)
		last = frame.Render()
		var err error
		enc, err = aoc.NewMP4Encoder("out.mp4", 60, log)
		if err != nil {
			log.Fatal(err)
		}
		if err := enc.Encode(last); err != nil {
			log.Fatal(err)
		}
	}

	for i, step := range stepList {
		log.Printf("%d/%d: %s, %+v", i+1, len(stepList), pos, step)
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
		for j := 0; j < step.Steps; j++ {
			if debug {
				term.Clear()
				term.MoveCursor(1, 1)
				println(pos.side.String())
				f := rune(strings.ToUpper(pos.facing.String())[0])
				maps[pos.side].Set(pos.pos, f)
				for side, tl := range layout {
					for i, line := range strings.Split(maps[side].String(), "\n") {
						term.MoveCursor(tl.X*2+1, tl.Y+i+2)
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
				if pos.pos.X < 1 || pos.pos.X > 48 || pos.pos.Y < 1 || pos.pos.Y > 48 {
					os.Stdin.Read([]byte{0})
				} else {
					time.Sleep(100 * time.Millisecond)
				}
			}

			from := pos
			nextpos := normalize(position{
				side:   pos.side,
				pos:    pos.pos.Move(pos.facing),
				facing: pos.facing,
			})

			switch maps[nextpos.side].At(nextpos.pos) {
			case '.':
				if video {
					fg := globalFromMap(from.side, from.pos)
					tg := globalFromMap(nextpos.side, nextpos.pos)
					pencil.Line(last,
						image.Pt(fg.X*aoc.GlyphWidth+aoc.GlyphWidth/2, fg.Y*aoc.LineHeight+aoc.LineHeight/2),
						image.Pt(tg.X*aoc.GlyphWidth+aoc.GlyphWidth/2, tg.Y*aoc.LineHeight+aoc.LineHeight/2),
						aoc.TolVibrantMagenta)
					if err := enc.Encode(last); err != nil {
						log.Fatal(err)
					}
				}
				pos = nextpos
				continue
			case '#':
				break stepCount
			default:
				log.Fatalf("bad position %v", nextpos)
			}
		}
	}

	if video {
		enc.Close()
	}

	value := map[coord.Direction]int{
		coord.East:  0,
		coord.South: 1,
		coord.West:  2,
		coord.North: 3,
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
