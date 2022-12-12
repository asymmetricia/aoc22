package main

import (
	"bytes"
	"math/rand"
	"os"
	"strings"
	"unicode"

	"github.com/sirupsen/logrus"

	"github.com/asymmetricia/aoc22/aoc"
	"github.com/asymmetricia/aoc22/canvas"
	"github.com/asymmetricia/aoc22/coord"
	"github.com/asymmetricia/aoc22/set"
)

var log = logrus.StandardLogger()

func frame(
	world coord.World,
	start coord.Coord,
	end coord.Coord,
	openSet map[coord.Coord]bool,
	from map[coord.Coord]coord.Coord,
) *canvas.Canvas {
	ret := &canvas.Canvas{}
	world.Each(func(c coord.Coord) (stop bool) {
		col := aoc.TolVibrantGrey
		ret.PrintAt(c.X, c.Y, string(world.At(c)), col)
		return false
	})

	for current := range openSet {
		cursor := current
		ok := true
		for ok {
			ret.PrintAt(cursor.X, cursor.Y, string(world.At(cursor)), aoc.TolVibrantCyan)
			cursor, ok = from[cursor]
		}
		ret.PrintAt(current.X, current.Y, string(world.At(current)), aoc.TolVibrantMagenta)
	}

	if openSet == nil {
		cursor, ok := start, true
		prev := cursor
		for ok {
			if cursor == start {
				ret.PrintAt(cursor.X, cursor.Y, "S", aoc.TolVibrantOrange)
				cursor, ok = from[cursor]
				continue
			}

			if cursor == end {
				ret.PrintAt(cursor.X, cursor.Y, "E", aoc.TolVibrantTeal)
				cursor, ok = from[cursor]
				break
			}

			next := from[cursor]
			symbol := '?'

			if prev.East() == cursor && cursor.North() == next {
				symbol = aoc.LineBR
			} else if prev.East() == cursor && cursor.South() == next {
				symbol = aoc.LineTR
			} else if prev.East() == cursor && cursor.East() == next ||
				prev.West() == cursor && cursor.West() == next {
				symbol = aoc.LineH
			} else if prev.West() == cursor && cursor.North() == next {
				symbol = aoc.LineBL
			} else if prev.West() == cursor && cursor.South() == next {
				symbol = aoc.LineTL
			} else if prev.North() == cursor && cursor.East() == next {
				symbol = aoc.LineTL
			} else if prev.North() == cursor && cursor.West() == next {
				symbol = aoc.LineTR
			} else if prev.North() == cursor && cursor.North() == next ||
				prev.South() == cursor && cursor.South() == next {
				symbol = aoc.LineV
			} else if prev.South() == cursor && cursor.East() == next {
				symbol = aoc.LineBL
			} else if prev.South() == cursor && cursor.West() == next {
				symbol = aoc.LineBR
			}

			ret.PrintAt(cursor.X, cursor.Y, string(symbol), aoc.TolVibrantOrange)

			prev = cursor
			cursor, ok = from[cursor]
		}
	}

	if start.X >= 0 {
		ret.PrintAt(start.X, start.Y, "E", aoc.TolVibrantTeal)
	}
	ret.PrintAt(end.X, end.Y, "E", aoc.TolVibrantTeal)

	return ret
}

func solution(name string, input []byte) int {
	// trim trailing space only
	input = bytes.Replace(input, []byte("\r"), []byte(""), -1)
	input = bytes.TrimRightFunc(input, unicode.IsSpace)
	lines := strings.Split(strings.TrimRightFunc(string(input), unicode.IsSpace), "\n")
	log.Printf("read %d %s lines", len(lines), name)

	world := coord.Load(lines, true)
	start := world.Find('S')[0]
	world.Set(start, 'a')
	end := world.Find('E')[0]
	world.Set(end, 'z')

	var frames []*canvas.Canvas

	neighbors := func(from coord.Coord) []coord.Coord {
		var ret []coord.Coord
		curHeight := world.At(from)
		for _, neighbor := range []coord.Coord{from.North(), from.East(), from.South(), from.West()} {
			neighborHeight := world.At(neighbor)
			if neighborHeight >= 0 && neighborHeight >= curHeight-1 {
				ret = append(ret, neighbor)
			}
		}
		return ret
	}

	goals := set.FromItems(world.Find('a'))
	path := aoc.AStarGraph[coord.Coord](
		end,
		goals,
		neighbors,
		func(a, b coord.Coord) int {
			return 1
		},
		func(a coord.Coord) int {
			return a.X
		},
		func(
			openSet map[coord.Coord]bool,
			cameFrom map[coord.Coord]coord.Coord,
			gScore map[coord.Coord]int,
			fScore map[coord.Coord]int,
			current coord.Coord,
		) {
			if goals[current] {
				f := frame(world, current, end, nil, cameFrom)
				f.Timing = 20
				frames = append(frames, f)
			} else if rand.Intn(10) == 0 {
				frames = append(frames, frame(world, coord.C(-1, -1), end, openSet, cameFrom))
			}
		},
	)

	canvas.RenderGif(frames, "day12-"+name+"-b.gif", log)

	return len(path) - 1
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
