package main

import (
	"bytes"
	"fmt"
	"os"
	"sort"
	"strconv"
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
}

func (v *Valve) String() string {
	return fmt.Sprintf("{rate=%d, neighbors=%v}", v.Rate, v.Peers)
}

type State struct {
	Network      map[string]*Valve `json:"-"`
	Pos          string
	Goal         string
	ElephantPos  string
	ElephantGoal string
	Open         map[string]bool
	Minutes      int
}

func (s *State) CK() string {
	var ret strings.Builder
	ret.WriteString(s.Pos)
	ret.WriteRune('|')
	ret.WriteString(s.Goal)
	ret.WriteRune('|')
	ret.WriteString(s.ElephantPos)
	ret.WriteRune('|')
	ret.WriteString(s.ElephantGoal)
	ret.WriteRune('|')
	ret.WriteString(strconv.Itoa(s.Minutes))
	ret.WriteRune('|')
	opens := maps.Keys(s.Open)
	sort.Strings(opens)
	ret.WriteString(strings.Join(opens, ","))
	return ret.String()
}

func (s *State) Copy() *State {
	return &State{
		Network:      s.Network,
		Pos:          s.Pos,
		Goal:         s.Goal,
		ElephantPos:  s.ElephantPos,
		ElephantGoal: s.ElephantGoal,
		Minutes:      s.Minutes,
		Open:         maps.Clone(s.Open),
	}
}

func (s *State) HumanAction(valves []string) []*State {
	//ck := "ha_" + s.CK()

	//checks++
	//if cache, ok := cache[ck]; ok {
	//	hits++
	//	return cache
	//}

	if s.Pos == s.Goal && !s.Open[s.Pos] {
		ret := s.Copy()
		ret.Open[s.Pos] = true
		return []*State{ret}
	}

	if s.Pos != s.Goal {
		path := paths[[2]string{s.Pos, s.Goal}]
		ret := s.Copy()
		ret.Pos = path[0]
		return []*State{ret}
	}

	var ret []*State
	for _, goal := range valves {
		if s.Open[goal] || s.ElephantGoal == goal {
			continue
		}

		path := paths[[2]string{s.Pos, goal}]

		//skip goals that are too far away
		//if len(path)+2 > s.Minutes {
		//	continue
		//}

		n := s.Copy()
		n.Goal = goal
		n.Pos = path[0]
		ret = append(ret, n)
	}

	//cache[ck] = ret
	return ret
}
func (s *State) ElephantAction(valves []string) []*State {
	//ck := "ea_" + s.CK()

	//checks++
	//if cache, ok := cache[ck]; ok {
	//	hits++
	//	return cache
	//}

	if s.ElephantPos == s.ElephantGoal && !s.Open[s.ElephantPos] {
		ret := s.Copy()
		ret.Open[s.ElephantPos] = true
		return []*State{ret}
	}

	if s.ElephantPos != s.ElephantGoal {
		path := paths[[2]string{s.ElephantPos, s.ElephantGoal}]
		ret := s.Copy()
		ret.ElephantPos = path[0]
		return []*State{ret}
	}

	var ret []*State
	for _, goal := range valves {
		if s.Open[goal] || s.Goal == goal {
			continue
		}

		path := paths[[2]string{s.ElephantPos, goal}]

		//skip goals that are too far away
		// len(path) == time to get there
		// 1 == time to open valve
		// 1 == time for opening valve to be worthwhile
		//if len(path)+2 > s.Minutes {
		//	continue
		//}

		n := s.Copy()
		n.ElephantGoal = goal
		n.ElephantPos = path[0]
		ret = append(ret, n)
	}

	//cache[ck] = ret
	return ret
}

var cache = map[string][]*State{}
var checks = 0
var hits = 0

func Act(s *State, depth int) []*State {
	if s.Minutes == 0 {
		return nil
	}

	ck := s.CK()
	checks++
	if cached, ok := cache[ck]; ok {
		hits++
		return cached
	}

	var candidates []*State
	for _, state := range s.HumanAction(maps.Keys(s.Network)) {
		ecs := state.ElephantAction(maps.Keys(s.Network))
		if len(ecs) > 0 {
			candidates = append(candidates, ecs...)
		} else {
			candidates = append(candidates, state)
		}
	}

	var best int
	var bestCandidate []*State
	//type seenKey struct {
	//	P, G, P2, G2 string
	//}
	//seen := map[seenKey]bool{}
	for i, cs := range candidates {
		cs.Minutes--

		//sk := seenKey{cs.Pos, cs.Goal, cs.ElephantPos, cs.ElephantGoal}
		//rk := seenKey{cs.ElephantPos, cs.ElephantGoal, cs.Pos, cs.Goal}
		//if seen[sk] || seen[rk] {
		//	continue
		//}
		//seen[sk] = true

		if depth <= 1 {
			log.Printf("chr %d%% (%d/%d) depth %d, %d/%d, best -> %d", hits*100/checks, hits, checks, depth, i, len(candidates), best)
			//log.Printf("depth %d, %d/%d, best -> %d", depth, i, len(candidates), best)
		}
		result := Act(cs, depth+1)
		if Value(result) > best {
			best = Value(result)
			bestCandidate = result
		}
	}

	ret := append([]*State{s}, bestCandidate...)
	cache[ck] = ret
	return ret
}

func (s State) String() string {
	var opens []string
	rate := 0
	for valve := range s.Open {
		rate += s.Network[valve].Rate
		if s.Network[valve].Rate > 0 {
			opens = append(opens, valve)
		}
	}
	sort.Strings(opens)
	return fmt.Sprintf("@%d me=%s->%s, el=%s->%s, open %v, rate %d", s.Minutes, s.Pos, s.Goal, s.ElephantPos, s.ElephantGoal, opens, rate)
}

func Value(states []*State) int {
	ret := 0
	for i, state := range states {
		for open := range state.Open {
			if i == len(states)-1 {
				ret += state.Network[open].Rate * state.Minutes
			} else {
				ret += state.Network[open].Rate
			}
		}
	}
	return ret
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

	start := State{
		Network:      network,
		Pos:          "AA",
		Goal:         "AA",
		ElephantPos:  "AA",
		ElephantGoal: "AA",
		Open:         map[string]bool{},
		Minutes:      26,
	}

	for id, valve := range network {
		if valve.Rate == 0 {
			start.Open[id] = true
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

	//cache = map[string][]*State{}
	//checks = 0
	//hits = 0
	if name == "test" {
		for i := 0; i < 10; i++ {
			path := Act(&start, 1)
			if Value(path) != 1707 {
				panic("nope")
			}
		}
	}
	log.Print(network)
	path := Act(&start, 1)
	for _, s := range path {
		log.Print(s)
	}
	log.Print(Value(path))

	return Value(path)
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
