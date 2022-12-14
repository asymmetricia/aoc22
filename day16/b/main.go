package main

import (
	"bytes"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
	"unicode"

	"github.com/sirupsen/logrus"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"

	"github.com/asymmetricia/aoc22/aoc"
)

var log = logrus.StandardLogger()

type Valve struct {
	Rate      int
	Neighbors map[string]*Valve
	Peers     []string
	Open      bool
}

func (v *Valve) String() string {
	return fmt.Sprintf("{rate=%d, neighbors=%v}", v.Rate, v.Peers)
}

// Explore returns a list of every legal path
func Explore(pos string, network map[string]*Valve, unopened map[string]bool, minutes int, n int) [][]string {
	// we could not actually get here lol
	if minutes < 0 {
		return nil
	}

	// open this valve
	if minutes > 0 && network[pos].Rate > 0 {
		minutes--
	}

	// nowhere left to go, let them run
	if len(unopened) == 0 || n == 0 {
		return [][]string{nil}
	}

	var ret [][]string
	uo := maps.Keys(unopened)
	sort.Strings(uo)
	for _, next := range uo {
		path := paths[[2]string{pos, next}]

		if len(path) > minutes {
			continue
		}

		// we will open this valve
		candidateUnopened := map[string]bool{}
		maps.Copy(candidateUnopened, unopened)
		delete(candidateUnopened, next)

		subpaths := Explore(next, network, candidateUnopened, minutes-len(path), n-1)
		for _, sp := range subpaths {
			ret = append(ret, slices.Insert(sp, 0, next))
		}
	}
	return ret
}

type World struct {
	Pos     []string
	Goal    []string
	Open    map[string]int
	Minutes int
}

func simulate(network map[string]*Valve, order []string) (map[string]int, int) {
	const actorCount = 1
	w := World{
		Open:    map[string]int{},
		Minutes: 26,
	}
	for i := 0; i < actorCount; i++ {
		w.Pos = append(w.Pos, "AA")
		w.Goal = append(w.Goal, order[0])
		order = order[1:]
	}
	var score int

	for w.Minutes > 0 {
		w.Minutes--
		for _, v := range w.Open {
			score += v
		}

		for i, pos := range w.Pos {
			// if we're at our goal...
			if pos == w.Goal[i] {
				_, open := w.Open[pos]
				// and it is not open...
				if !open {
					// open it
					w.Open[pos] = network[pos].Rate
				} else if len(order) > 0 {
					// otherwise pick a new goal
					w.Goal[i] = order[0]
					order = order[1:]
				}
			}
			// if we are not at our goal (maybe because we picked a new one?)
			if pos != w.Goal[i] {
				// move toward it
				w.Pos[i] = paths[[2]string{pos, w.Goal[i]}][0]
			}
		}
	}
	return w.Open, score
}

var paths = map[[2]string][]string{}

func solution(name string, input []byte) int {
	// trim trailing space only
	input = bytes.Replace(input, []byte("\r"), []byte(""), -1)
	input = bytes.TrimRightFunc(input, unicode.IsSpace)
	lines := strings.Split(strings.TrimRightFunc(string(input), unicode.IsSpace), "\n")
	log.Printf("read %d %s lines", len(lines), name)

	network := map[string]*Valve{}

	for _, line := range lines {
		fields := strings.Fields(line)
		name := fields[1]
		rate := aoc.MustAtoi(aoc.Before(aoc.After(fields[4], "="), ";"))
		peers := fields[9:]
		for i := range peers {
			peers[i] = strings.Trim(peers[i], ",")
		}
		network[name] = &Valve{Rate: rate, Peers: peers, Neighbors: map[string]*Valve{}}
	}

	valveId := map[string]int{}
	valves := map[string]bool{}
	for id, valve := range network {
		valveId[id] = len(valveId) + 1
		if valve.Rate > 0 {
			valves[id] = true
		}
		for _, peer := range valve.Peers {
			valve.Neighbors[peer] = network[peer]
		}
	}

	paths = map[[2]string][]string{}
	for a := range network {
		for b := range network {
			if a == b {
				continue
			}

			path := aoc.Dijkstra(a, b, func(a string) []string {
				return maps.Keys(network[a].Neighbors)
			}, aoc.ConstantCost[string])
			path = path[1:]
			paths[[2]string{a, b}] = path
		}
	}

	best := 0
	var bestPath [2][]string

	type path struct {
		path   []string
		value  int
		valves int64
	}
	var ourPaths []path
	for i := 1; i <= 15; i++ {
		p := Explore("AA", network, valves, 26, i)
		for _, p := range p {
			valves, v := simulate(network, p)
			var vb int64
			for valve := range valves {
				vb |= 1 << valveId[valve]
			}
			ourPaths = append(ourPaths, path{p, v, vb})
		}
	}

	last := time.Now()
	for i, path := range ourPaths {
		if time.Since(last) > time.Second {
			log.Printf("%d/%d (%d%%)", i+1, len(ourPaths), (i+1)*100/len(ourPaths))
			last = time.Now()
		}
		for _, pair := range ourPaths[i+1:] {
			if path.valves&pair.valves > 0 {
				continue
			}
			if path.value+pair.value > best {
				best = path.value + pair.value
				bestPath = [2][]string{path.path, pair.path}
			}
		}
	}

	log.Printf("%v -> %d", bestPath, best)

	return best
}

func main() {
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02T15:04:05",
	})
	test, err := os.ReadFile("test")
	if err == nil {
		ts := solution("test", test)
		log.Printf("test solution: %d", ts)
		if ts != 1707 {
			panic("nope")
		}
	} else {
		log.Warningf("no test data present")
	}

	input := aoc.Input(2022, 16)
	log.Printf("input solution: %d", solution("input", input))
}
