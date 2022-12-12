package main

import (
	"bytes"
	"os"
	"strings"
	"unicode"

	"github.com/sirupsen/logrus"

	"github.com/asymmetricia/aoc22/aoc"
	"github.com/asymmetricia/aoc22/coord"
	"github.com/asymmetricia/aoc22/set"
)

var log = logrus.StandardLogger()

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

	neighbors := func(from coord.Coord) []coord.Coord {
		var ret []coord.Coord
		curHeight := world.At(from)
		for _, neighbor := range []coord.Coord{from.North(), from.East(), from.South(), from.West()} {
			neighborHeight := world.At(neighbor)
			if neighborHeight >= 0 && (neighborHeight == curHeight+1 || neighborHeight <= curHeight) {
				ret = append(ret, neighbor)
			}
		}
		log.Printf("%v -> %v", from, ret)
		return ret
	}

	path := aoc.AStarGraph[coord.Coord](
		start,
		set.Set[coord.Coord]{end: true},
		neighbors,
		func(a, b coord.Coord) int {
			return 1
		},
		func(a coord.Coord) int {
			return 1
		},
		func(
			openSet map[coord.Coord]bool,
			cameFrom map[coord.Coord]coord.Coord,
			gScore map[coord.Coord]int,
			fScore map[coord.Coord]int,
			current coord.Coord,
		) {
		},
	)

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
