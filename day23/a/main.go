package main

import (
	"bytes"
	"os"
	"strings"
	"unicode"

	"github.com/asymmetricia/aoc22/aoc"
	"github.com/asymmetricia/aoc22/coord"
	"github.com/sirupsen/logrus"
)

var log = logrus.StandardLogger()

type world struct {
	world    *coord.SparseWorld
	proposal map[coord.Coord][]coord.Coord
	cons     []consideration
}

type consideration struct {
	Neighbors []coord.Direction
	Direction coord.Direction
}

func (c consideration) consider(w *world, elves []coord.Coord) []coord.Coord {
	var ret []coord.Coord
elves:
	for _, elf := range elves {
		hasNeigh := false
		for _, neigh := range elf.Neighbors(true) {
			if w.world.At(neigh) == '#' {
				hasNeigh = true
				break
			}
		}

		if !hasNeigh {
			ret = append(ret, elf)
			// elf does not need to move
			continue
		}

		for _, dir := range c.Neighbors {
			if w.world.At(elf.Move(dir)) == '#' {
				ret = append(ret, elf)
				continue elves
			}
		}
		pos := elf.Move(c.Direction)
		w.proposal[pos] = append(w.proposal[pos], elf)
	}
	return ret
}

// If no other Elves are in one of those eight positions, the Elf does not do
// anything during this round. Otherwise, the Elf looks in each of four
// directions in the following order and proposes moving one step in the first
// valid direction:
//
//    If there is no Elf in the N, NE, or NW adjacent positions, the Elf proposes moving north one step.
//    If there is no Elf in the S, SE, or SW adjacent positions, the Elf proposes moving south one step.
//    If there is no Elf in the W, NW, or SW adjacent positions, the Elf proposes moving west one step.
//    If there is no Elf in the E, NE, or SE adjacent positions, the Elf proposes moving east one step.

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

	w := &world{
		world:    coord.Load(lines, false).(*coord.SparseWorld),
		proposal: map[coord.Coord][]coord.Coord{},
		cons: []consideration{
			{[]coord.Direction{coord.North, coord.NorthEast, coord.NorthWest}, coord.North},
			{[]coord.Direction{coord.South, coord.SouthEast, coord.SouthWest}, coord.South},
			{[]coord.Direction{coord.West, coord.NorthWest, coord.SouthWest}, coord.West},
			{[]coord.Direction{coord.East, coord.NorthEast, coord.SouthEast}, coord.East},
		}}

	for _, dot := range w.world.Find('.') {
		w.world.Set(dot, 0)
	}

	ec := len(w.world.Find('#'))

	count := 0
	round := func() bool {
		log := log.WithField("round", count+1)
		init := w.world.Find('#')
		if len(init) != ec {
			panic("elf count changed")
		}
		elves := init
		for _, cons := range w.cons {
			elves = cons.consider(w, elves)
		}
		w.cons[0], w.cons[1], w.cons[2], w.cons[3] = w.cons[1], w.cons[2], w.cons[3], w.cons[0]
		if len(init) == len(elves) {
			return false
		}

		for to, from := range w.proposal {
			if len(from) != 1 {
				continue
			}
			if w.world.At(from[0]) != '#' {
				log.Fatalf("no elf at from pos %v", from[0])
			}
			w.world.Set(from[0], 0)
			w.world.Set(to, '#')
			//log.Printf("from %v to %v, %d elves", from[0], to, len(*w.world))
		}
		w.proposal = map[coord.Coord][]coord.Coord{}
		//if count%100 == 0 {
		//term.Clear()
		//term.MoveCursor(1, 1)
		//w.world.Print()
		//println(len(init))
		//}
		//os.Stdin.Read([]byte{0})
		count++
		if count == 10 {
			return false
		}
		return true
	}

	for round() {
	}

	w.world.Print()

	minx, miny, maxx, maxy := w.world.Rect()
	return (maxx-minx+1)*(maxy-miny+1) - len(w.world.Find('#'))
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

	input := aoc.Input(2022, 23)
	log.Printf("input solution: %d", solution("input", input))
}
