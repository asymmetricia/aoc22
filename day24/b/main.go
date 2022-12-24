package main

import (
	"bytes"
	"os"
	"strings"
	"unicode"

	"github.com/sirupsen/logrus"
	"golang.org/x/exp/maps"

	"github.com/asymmetricia/aoc22/aoc"
	"github.com/asymmetricia/aoc22/coord"
)

var log = logrus.StandardLogger()

type state struct {
	world     *coord.SparseWorld
	blizzards map[coord.Direction][]coord.Coord
}

func (s state) step() state {
	ret := state{
		world:     s.world.Copy().(*coord.SparseWorld),
		blizzards: map[coord.Direction][]coord.Coord{},
	}

	for _, locations := range s.blizzards {
		for _, loc := range locations {
			ret.world.Set(loc, 0)
		}
	}

	minx, miny, maxx, maxy := s.world.Rect()

	for dir, locations := range s.blizzards {
		var nextlocs []coord.Coord
		for _, loc := range locations {
			next := loc.Move(dir)
			if s.world.At(next) == '#' {
				switch dir {
				case coord.East:
					next.X = minx + 1
				case coord.West:
					next.X = maxx - 1
				case coord.North:
					next.Y = maxy - 1
				case coord.South:
					next.Y = miny + 1
				}
			}
			nextlocs = append(nextlocs, next)
		}
		ret.blizzards[dir] = nextlocs
	}

	for dir, locs := range ret.blizzards {
		for _, loc := range locs {
			switch ret.world.At(loc) {
			case -1:
				ret.world.Set(loc, map[coord.Direction]rune{
					coord.East:  '>',
					coord.West:  '<',
					coord.North: '^',
					coord.South: 'v',
				}[dir])
			case '>':
				fallthrough
			case '<':
				fallthrough
			case 'v':
				fallthrough
			case '^':
				ret.world.Set(loc, '2')
			case '2':
				ret.world.Set(loc, '3')
			case '4':
				ret.world.Set(loc, '4')
			}
		}
	}

	return ret
}

type pos struct {
	c coord.Coord
	t int
}

func (a pos) Equal(b pos) bool {
	return a.c == b.c
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

	initial := state{
		world: coord.Load(lines, false).(*coord.SparseWorld),
	}

	for _, c := range initial.world.Find('.') {
		initial.world.Set(c, 0)
	}

	initial.blizzards = map[coord.Direction][]coord.Coord{
		coord.East:  initial.world.Find('>'),
		coord.West:  initial.world.Find('<'),
		coord.South: initial.world.Find('v'),
		coord.North: initial.world.Find('^'),
	}

	if name == "test" {
		initial.world.Print()
		initial.step().world.Print()
	}

	graph := map[int]coord.SparseWorld{
		0: maps.Clone(*initial.world),
	}
	state := initial

	minx, miny, maxx, maxy := initial.world.Rect()

	path := aoc.Dijkstra(
		pos{coord.C(1, 0), 0},
		pos{coord.C(maxx-1, maxy), 99999},
		func(a pos) []pos {
			if a.t+1 >= len(graph) {
				state = state.step()
				graph[len(graph)] = *state.world
			}

			var ret []pos
			ns := graph[a.t+1]
			for _, n := range append(a.c.Neighbors(false), a.c) {
				if n.Y == maxy && n.X == maxx-1 ||
					n.Y == miny && n.X == minx+1 ||
					n.X > minx && n.X < maxx &&
						n.Y > miny && n.Y < maxy &&
						ns.At(n) == -1 {
					ret = append(ret, pos{n, a.t + 1})
				}
			}
			return ret
		}, aoc.ConstantCost[pos], func(q *aoc.PQueue[pos], dist map[pos]int, prev map[pos]pos, current pos) {
		})

	path2 := aoc.Dijkstra(
		pos{coord.C(maxx-1, maxy), len(path) - 1},
		pos{coord.C(minx+1, miny), 99999},
		func(a pos) []pos {
			if a.t+1 >= len(graph) {
				state = state.step()
				graph[len(graph)] = *state.world
			}

			var ret []pos
			ns := graph[a.t+1]
			for _, n := range append(a.c.Neighbors(false), a.c) {
				if n.Y == maxy && n.X == maxx-1 ||
					n.Y == miny && n.X == minx+1 ||
					n.X > minx && n.X < maxx &&
						n.Y > miny && n.Y < maxy &&
						ns.At(n) == -1 {
					ret = append(ret, pos{n, a.t + 1})
				}
			}
			return ret
		}, aoc.ConstantCost[pos], func(q *aoc.PQueue[pos], dist map[pos]int, prev map[pos]pos, current pos) {
		})

	path3 := aoc.Dijkstra(
		pos{coord.C(minx+1, miny), len(path) - 1 + len(path2) - 1},
		pos{coord.C(maxx-1, maxy), 99999},
		func(a pos) []pos {
			if a.t+1 >= len(graph) {
				state = state.step()
				graph[len(graph)] = *state.world
			}

			var ret []pos
			ns := graph[a.t+1]
			for _, n := range append(a.c.Neighbors(false), a.c) {
				if n.Y == maxy && n.X == maxx-1 ||
					n.Y == miny && n.X == minx+1 ||
					n.X > minx && n.X < maxx &&
						n.Y > miny && n.Y < maxy &&
						ns.At(n) == -1 {
					ret = append(ret, pos{n, a.t + 1})
				}
			}
			return ret
		}, aoc.ConstantCost[pos], func(q *aoc.PQueue[pos], dist map[pos]int, prev map[pos]pos, current pos) {
		})

	log.Print(len(path)-1, len(path2)-1, len(path3)-1)
	log.Print(path)

	if name == "test" && len(path)-1 != 18 {
		panic("nope")
	}
	return len(path) - 1 + len(path2) - 1 + len(path3) - 1
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

	input := aoc.Input(2022, 24)
	log.Printf("input solution: %d", solution("input", input))
}
