package main

import (
	"bytes"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"sort"
	"strconv"
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

type State struct {
	Network map[string]*Valve `json:"-"`
	Pos     string
	//Goal         string
	ElephantPos string
	//ElephantGoal string
	Open   []string
	Closed []string
	//Open   map[string]bool
	OpenCk string
	//Minutes int
	Minute int
}

func (s *State) IsOpen(valve string) bool {
	for _, o := range s.Open {
		if o == valve {
			return true
		}
	}
	return false
}

func (s *State) MarkOpen(valve string) {
	s.Open = append(s.Open, valve)
	ci := slices.Index(s.Closed, valve)
	if ci == -1 {
		panic("opening open valve")
	}
	s.Closed = append(s.Closed[:ci], s.Closed[ci+1:]...)
	//s.Open[valve] = true
	//list := maps.Keys(s.Open)
	sort.Strings(s.Open)
	s.OpenCk = strings.Join(s.Open, ",")
}

func (s *State) CK() string {
	var ret strings.Builder

	if strings.Compare(s.Pos, s.ElephantPos) == -1 {
		ret.WriteString(s.Pos)
		ret.WriteRune('|')
		ret.WriteString(s.ElephantPos)
	} else {
		ret.WriteString(s.ElephantPos)
		ret.WriteRune('|')
		ret.WriteString(s.Pos)
	}
	ret.WriteRune('|')
	ret.WriteString(strconv.Itoa(s.Minute))
	ret.WriteRune('|')
	ret.WriteString(s.OpenCk)
	return ret.String()
}

func (s *State) Copy() *State {
	return &State{
		Network:     s.Network,
		Pos:         s.Pos,
		ElephantPos: s.ElephantPos,
		//Minutes:     s.Minutes,
		Minute: s.Minute,
		Closed: slices.Clone(s.Closed),
		Open:   slices.Clone(s.Open),
		//Open:   maps.Clone(s.Open),
		OpenCk: s.OpenCk,
	}
}

func (s *State) HumanAction() []*State {
	pos := s.Pos
	set := func(n *State, v string) { n.Pos = v }

	return act(s, pos, set)
}
func (s *State) ElephantAction() []*State {
	pos := s.ElephantPos
	set := func(n *State, v string) { n.ElephantPos = v }

	return act(s, pos, set)
}

func act(s *State, pos string, setPosition func(s *State, v string)) []*State {
	var moves []*State

	// consider staying here & opening the valve
	// minute 26 is our last minute. If it's minute 24, we'll open for 25 to produce.
	if s.Minute <= 24 && !s.IsOpen(pos) {
		//if s.Minutes > 1 && !s.IsOpen(pos) {
		n := s.Copy()
		n.MarkOpen(pos)
		moves = append(moves, n)
	}

	// consider moving toward another closed valve
	nextSteps := map[string]bool{}
	for _, id := range s.Closed {
		// not a distant closed valve
		if id == pos {
			continue
		}

		path := paths[[2]string{pos, id}]

		// we need len(path) minutes to get there, 1 minute to open valve, 1 minute to
		// produce. otherwise no point. We _have_ 26 - s.Minute minutes of action left.
		if len(path)+2 >= (26 - s.Minute) {
			continue
		}
		if !nextSteps[path[0]] {
			nextSteps[path[0]] = true
			n := s.Copy()
			setPosition(n, path[0])
			moves = append(moves, n)
		}
	}

	return moves
}

var cache = map[string][]*State{}
var checks = 1
var hits = 0
var last = time.Now()

func Act(s *State, depth int) []*State {
	if s.Minute == 26 {
		return nil
	}

	//checks++
	//var ck string
	//ck = s.CK()
	//if cached, ok := cache[ck]; ok {
	//	hits++
	//	return cached
	//}

	var candidates []*State
	hNext := s.HumanAction()
	if len(hNext) == 0 {
		eNext := s.ElephantAction()
		if len(eNext) == 0 {
			return nil
		}
	}

	for _, state := range hNext {
		eNext := state.ElephantAction()
		if len(eNext) == 0 {
			candidates = append(candidates, state)
		} else {
			candidates = append(candidates, eNext...)
		}
	}

	for _, cs := range candidates {
		cs.Minute++
	}

	var best int
	var bestCandidate []*State
	for i, cs := range candidates {
		result := Act(cs, depth+1)
		if Value(result) > best {
			best = Value(result)
			bestCandidate = result
		}

		//if Matches(append([]*State{s}, result...), 26, [][2]string{{"AA", "AA"},{"DD", "II"}, {"DD", "JJ"}, {"EE", "JJ"}}) {
		//	log.Print(s)
		//	for _, s := range result {
		//		log.Print(s)
		//	}
		//}

		if depth == 1 {
			log.Printf("chr %d%% (%d/%d) depth %d, %d/%d, best -> %d", hits*100/checks, hits, checks, depth, i+1, len(candidates), best)
			log.Print(Path(result))
			depths := map[int]int{}
			for _, s := range cache {
				depths[len(s)]++
			}
			log.Print(depths)
			last = time.Now()
			//log.Printf("depth %d, %d/%d, best -> %d", depth, i, len(candidates), best)
		}
	}

	ret := append([]*State{s}, bestCandidate...)
	//cache[ck] = ret
	return ret
}

func (s *State) String() string {
	var opens []string
	rate := 0
	for _, valve := range s.Open {
		rate += s.Network[valve].Rate
		if s.Network[valve].Rate > 0 {
			opens = append(opens, valve)
		}
	}
	sort.Strings(opens)
	return fmt.Sprintf("{@%d me=%s, el=%s, open %v, rate %d}", s.Minute, s.Pos, s.ElephantPos, opens, rate)
}

func Value(states []*State) int {
	ret := 0
	for i, state := range states {
		for _, open := range state.Open {
			if i == len(states)-1 {
				ret += state.Network[open].Rate * (26 - state.Minute)
			} else {
				ret += state.Network[open].Rate
			}
		}
	}
	return ret
}

func Path(states []*State) string {
	var ret strings.Builder
	for i, state := range states {
		if i > 0 {
			ret.WriteString("->")
		}
		ret.WriteRune('[')
		ret.WriteString(state.Pos)
		ret.WriteRune(',')
		ret.WriteString(state.ElephantPos)
		ret.WriteRune(']')
	}
	return ret.String()
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
		// 30min total
	}

	start := &State{
		Network:     network,
		Pos:         "AA",
		ElephantPos: "AA",
		//Open:        map[string]bool{},
		//Minutes:     26,
		Minute: 0,
	}

	for id := range network {
		start.Closed = append(start.Closed, id)
	}

	for id, valve := range network {
		if valve.Rate == 0 {
			start.MarkOpen(id)
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

	cache = map[string][]*State{}
	checks = 1
	hits = 0
	if name == "test" {
		for i := 0; i < 10; i++ {
			path := Act(start, 1)
			log.Printf("%v -> %d", path, Value(path))
			if Value(path) != 1707 {
				for _, state := range path {
					log.Print(state)
				}
				panic("nope")
			}
		}
	}
	log.Print(network)
	path := Act(start, 1)
	for _, s := range path {
		log.Print(s)
	}
	log.Print(Value(path))

	return Value(path)
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
