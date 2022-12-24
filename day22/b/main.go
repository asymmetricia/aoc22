package main

import (
	"bytes"
	"image"
	"image/color"
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
	"github.com/asymmetricia/aoc22/isovox"
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

var recency = map[side]map[coord.Coord]int{}

func globalFromMap(s side, c coord.Coord) coord.Coord {
	return c.Plus(layout[s])
}

func render(maps map[side]*coord.DenseWorld) image.Image {
	colors := map[rune]color.Color{
		'N': aoc.TolVibrantMagenta,
		'E': aoc.TolVibrantMagenta,
		'W': aoc.TolVibrantMagenta,
		'S': aoc.TolVibrantMagenta,
		'#': aoc.TolVibrantCyan,
	}
	ivx := &isovox.World{Voxels: map[isovox.Coord]*isovox.Voxel{}}
	for side, origin := range layout {
		for y, row := range *maps[side] {
			for x, cell := range row {
				c, ok := colors[cell]
				if _, sideRecOk := recency[side]; !ok && sideRecOk {
					var r int
					r, ok = recency[side][coord.C(x, y)]
					c = aoc.TolScale(0, 1000, r)
					if r > 1000 {
						ok = false
					}
				}

				y := 200 - y - origin.Y
				ivx.Voxels[isovox.Coord{x + origin.X, y, -1}] = &isovox.Voxel{Color: aoc.TolVibrantGrey}
				if ok {
					ivx.Voxels[isovox.Coord{x + origin.X, y, 0}] = &isovox.Voxel{Color: c}
					if c == aoc.TolVibrantMagenta {
						ivx.Voxels[isovox.Coord{x + origin.X, y, 1}] = &isovox.Voxel{Color: aoc.TolVibrantMagenta}
					}
				}
			}
		}
	}

	dim := len(*maps[Top])

	cube1x := -10
	cube1y := 112

	for x := 0; x < 50; x++ {
		for y := 0; y < 50; y++ {
			for z := 0; z < 50; z++ {
				ivx.Voxels[isovox.Coord{cube1x + x, cube1y + y, z}] = &isovox.Voxel{Color: aoc.TolVibrantGrey}
			}
		}
	}

	for y, row := range *maps[Top] {
		for x, cell := range row {
			c, ok := colors[cell]
			if _, sideRecOk := recency[Top]; !ok && sideRecOk {
				var r int
				r, ok = recency[Top][coord.C(x, y)]
				c = aoc.TolScale(0, 1000, r)
				if r > 1000 {
					ok = false
				}
			}
			if ok {
				ivx.Voxels[isovox.Coord{cube1x + x, cube1y + 49 - y, dim}] = &isovox.Voxel{Color: c}
			}
		}
	}

	// east is on YZ plane
	for y, row := range *maps[East] {
		for x, cell := range row {
			c, ok := colors[cell]
			if _, sideRecOk := recency[East]; !ok && sideRecOk {
				var r int
				r, ok = recency[East][coord.C(x, y)]
				c = aoc.TolScale(0, 1000, r)
				if r > 1000 {
					ok = false
				}
			}
			if ok {
				ivx.Voxels[isovox.Coord{
					cube1x + 50,
					cube1y + 49 - y,
					50 - x}] = &isovox.Voxel{Color: c}
			}
		}
	}

	// south is on xz plane
	for y, row := range *maps[South] {
		y = 50 - y
		for x, cell := range row {
			if c, ok := colors[cell]; ok {
				ivx.Voxels[isovox.Coord{
					cube1x + x,
					cube1y,
					y}] = &isovox.Voxel{Color: c}
			}
		}
	}

	cube2x := -62
	cube2y := 60
	for y, row := range *maps[West] {
		y = 50 - y
		for x, cell := range row {
			for z := 0; z < 50; z++ {
				ivx.Voxels[isovox.Coord{
					cube2x + x,
					cube2y + y,
					z}] = &isovox.Voxel{Color: aoc.TolVibrantGrey}
			}
			if c, ok := colors[cell]; ok {
				ivx.Voxels[isovox.Coord{
					cube2x + x,
					cube2y + y,
					dim}] = &isovox.Voxel{Color: c}
			}
		}
	}

	for y, row := range *maps[Bottom] {
		for x, cell := range row {
			if c, ok := colors[cell]; ok {
				ivx.Voxels[isovox.Coord{
					cube2x + 50,
					cube2y + 50 - y,
					49 - x}] = &isovox.Voxel{Color: c}
			}
		}
	}

	for y, row := range *maps[North] {
		for x, cell := range row {
			if c, ok := colors[cell]; ok {
				ivx.Voxels[isovox.Coord{
					cube2x + x,
					cube2y,
					49 - y,
				}] = &isovox.Voxel{Color: c}
			}
		}
	}

	return ivx.Render(6)
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
				f := rune(strings.ToUpper(pos.facing.String())[0])
				term.Clear()
				term.MoveCursor(1, 1)
				println(pos.side.String())
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
				maps[pos.side].Set(pos.pos, '.')
				maps[nextpos.side].Set(nextpos.pos, rune(strings.ToUpper(nextpos.facing.String())[0]))
				for _, recs := range recency {
					for c, v := range recs {
						recs[c] = v + 1
					}
				}
				if _, ok := recency[nextpos.side]; !ok {
					recency[nextpos.side] = map[coord.Coord]int{}
				}
				recency[nextpos.side][nextpos.pos] = 0
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

	aoc.RenderPng(render(maps), "day22-b-"+name+".png")
	for _, m := range maps {
		m.Print()
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
