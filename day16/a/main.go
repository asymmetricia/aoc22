package main

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"sort"
	"strings"
	"unicode"

	"github.com/sirupsen/logrus"
	"golang.org/x/exp/maps"

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

// Explore considers moving to each possible unopened valve. It returns the best
// discovered order of valve openings, and an integer describing the total rate
// released in this process.
func Explore(pos string, network map[string]*Valve, unopened map[string]bool, minutes int, openRate int) ([]string, int) {
	var best int
	var bestPath []string
	var ret = 0

	// we could not actually get here lol
	if minutes < 0 {
		return nil, math.MinInt
	}

	// open this valve
	if minutes > 0 && network[pos].Rate > 0 {
		minutes--
		ret += openRate
		openRate += network[pos].Rate
	}

	// nowhere left to go, let them run
	if len(unopened) == 0 {
		return nil, ret + openRate*minutes
	}

	uo := maps.Keys(unopened)
	sort.Strings(uo)
	for i, next := range uo {
		// Compute a path to the next unopened valve
		path := aoc.Dijkstra(pos, next, func(a string) []string {
			return maps.Keys(network[a].Neighbors)
		}, func(a, b string) int {
			return 1
		})
		path = path[1:]

		if len(path) > minutes {
			continue
		}

		// we will open this valve
		candidateUnopened := map[string]bool{}
		maps.Copy(candidateUnopened, unopened)
		delete(candidateUnopened, next)

		candidate, value := Explore(next, network, candidateUnopened, minutes-len(path), openRate)
		value += openRate * len(path)
		if value > best {
			best = value
			bestPath = append([]string{next}, candidate...)
		}
		if pos == "AA" {
			log.Printf("%d/%d == %d", i, len(uo), ret+value)
		}
	}
	return bestPath, ret + best
}

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
		// 30min total
	}

	valves := map[string]bool{}
	for id, valve := range network {
		if valve.Rate > 0 {
			valves[id] = true
		}
		for _, peer := range valve.Peers {
			valve.Neighbors[peer] = network[peer]
		}
	}

	log.Print(network)
	path, value := Explore("AA", network, valves, 30, 0)
	log.Print(path, value)

	return value
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

	input := aoc.Input(2022, 16)
	log.Printf("input solution: %d", solution("input", input))
}
