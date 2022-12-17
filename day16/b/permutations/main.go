package main

import (
	"bytes"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"reflect"
	"sort"
	"strings"
	"time"
	"unicode"

	"github.com/asymmetricia/aoc22/aoc"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

var log = logrus.StandardLogger()

type Valve struct {
	Rate      int
	Neighbors map[string]*Valve
	Peers     []string
}

func (v *Valve) String() string {
	return fmt.Sprintf("{rate=%d, neighbors=%v}", v.Rate, v.Peers)
}

var paths = map[[2]string][]string{}

func permutations(in []string, length int, maxlength int, log logrus.FieldLogger) <-chan []string {
	ch := make(chan []string)
	go func(in []string, length int, maxlength int, ch chan<- []string) {
		defer close(ch)
		if len(in) == 1 {
			ch <- []string{in[0]}
			return
		}

		if length == 1 {
			for _, v := range in {
				ch <- []string{v}
			}
			return
		}

		rest := make([]string, 0, len(in)-1)
		for i, v := range in {
			rest = rest[:0]
			rest = append(rest, in[:i]...)
			rest = append(rest, in[i+1:]...)
			for end := range permutations(rest, length-1, maxlength, log) {
				ch <- slices.Insert(end, 0, v)
			}
		}
	}(in, length, maxlength, ch)
	return ch
}

type World struct {
	Pos     []string
	Goal    []string
	Open    map[string]int
	Minutes int
}

func simulate(network map[string]*Valve, order []string) (map[string]int, int) {
	const actorCount = 2
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

	var valves []string
	for id, valve := range network {
		for _, peer := range valve.Peers {
			valve.Neighbors[peer] = network[peer]
		}
		if valve.Rate > 0 {
			valves = append(valves, id)
		}
	}
	sort.Strings(valves)

	// pre-compute shortest paths
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

	const workers = 16
	var orders = make(chan []string)
	var resultsCh []reflect.SelectCase
	type result struct {
		Order  []string
		Opened map[string]int
		Value  int
	}
	for i := 0; i < workers; i++ {
		rch := make(chan result)
		go func(orders <-chan []string, ch chan<- result) {
			defer close(ch)
			for order := range orders {
				opened, value := simulate(network, order)
				ch <- result{order, opened, value}
			}
		}(orders, rch)
		resultsCh = append(resultsCh, reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(rch),
		})
	}

	n := len(valves)/2 + len(valves)%2

	go func() {
		defer close(orders)
		for order := range permutations(valves, n, n, log) {
			orders <- order
		}
	}()

	var total int64 = 1
	for i := 0; i < n; i++ {
		total *= int64(len(valves) - i)
	}

	count := 0
	var lastPerc int64 = -1
	last := time.Now()

	var results []result
	for {
		i, resultI, ok := reflect.Select(resultsCh)
		if !ok {
			resultsCh = append(resultsCh[:i], resultsCh[i+1:]...)
			if len(resultsCh) == 0 {
				break
			}
		}
		count++
		if time.Since(last) > time.Second {
			perc := int64(count) * 100 / total
			if lastPerc != perc {
				log.Printf("%d/%d (%d%%)", count, total, int64(count)*100/total)
			}
			lastPerc = perc
			last = time.Now()
		}

		result := resultI.Interface().(result)
		results = append(results, result)
	}

	var best int
	var a, b result

	for i, result := range results {
		pairs:
		for _, pair := range results[i+1:] {
			for v := range result.Opened {
				if _, ok := pair.Opened[v]; ok {
					continue pairs
				}
			}
			for v := range pair.Opened {
				if _, ok := result.Opened[v]; ok {
					continue pairs
				}
			}
			if result.Value + pair.Value > best {
				best = result.Value + pair.Value
				a = result
				b = pair
			}
		}
	}

	log.Printf("%v + %v = %d", a, b, best)

	return best
}

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

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
