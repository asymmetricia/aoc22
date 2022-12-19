package main

import (
	"bytes"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/sirupsen/logrus"
	"golang.org/x/exp/maps"

	"github.com/asymmetricia/aoc22/aoc"
)

var log = logrus.StandardLogger()

type Cost struct {
	Ore, Clay, Obsidian uint8
	Geode               uint16
}

func (c Cost) Have(b Cost) bool {
	return c.Ore >= b.Ore && c.Clay >= b.Clay && c.Obsidian >= b.Obsidian && c.Geode >= b.Geode
}

func (c Cost) Minus(b Cost) Cost {
	return Cost{c.Ore - b.Ore, c.Clay - b.Clay, c.Obsidian - b.Obsidian, c.Geode - b.Geode}
}

func (c Cost) Plus(b Cost) Cost {
	return Cost{c.Ore + b.Ore, c.Clay + b.Clay, c.Obsidian + b.Obsidian, c.Geode + b.Geode}
}

type Blueprint struct {
	Ore      Cost
	Clay     Cost
	Obsidian Cost
	Geode    Cost
}

type State struct {
	Robots    Cost
	Resources Cost
}

func (s State) Wait() State {
	return State{
		Robots:    s.Robots,
		Resources: s.Resources.Plus(s.Robots),
	}
}

func (s State) Ore(c Cost) State {
	return State{
		Robots:    s.Robots.Plus(Cost{Ore: 1}),
		Resources: s.Resources.Plus(s.Robots).Minus(c),
	}
}

func (s State) Clay(c Cost) State {
	return State{
		Robots:    s.Robots.Plus(Cost{Clay: 1}),
		Resources: s.Resources.Plus(s.Robots).Minus(c),
	}
}

func (s State) Obsidian(c Cost) State {
	return State{
		Robots:    s.Robots.Plus(Cost{Obsidian: 1}),
		Resources: s.Resources.Plus(s.Robots).Minus(c),
	}
}

func (s State) Geode(c Cost) State {
	return State{
		Robots:    s.Robots.Plus(Cost{Geode: 1}),
		Resources: s.Resources.Plus(s.Robots).Minus(c),
	}
}

func (s State) TheoreticalBest(minutes uint16, geode Cost) uint16 {
	maxBots := minutes - 2
	if !s.Resources.Have(geode) {
		maxBots--
	}
	return maxBots*(maxBots+1)/2 + minutes*s.Robots.Geode + s.Resources.Geode
}

var last = time.Now()

func (bp *Blueprint) Geodes(minute uint16, initial State) uint16 {
	states := map[State]bool{
		initial: true,
	}

	var max Cost
	for _, c := range []Cost{bp.Ore, bp.Clay, bp.Obsidian, bp.Geode} {
		max.Ore = aoc.Max(max.Ore, c.Ore)
		max.Clay = aoc.Max(max.Clay, c.Clay)
		max.Obsidian = aoc.Max(max.Obsidian, c.Obsidian)
	}

	for minute > 1 {
		next := map[State]bool{}

		var geodes uint16
		for s := range states {
			if s.Resources.Geode+s.Robots.Geode*(minute-1) > geodes {
				geodes = s.Resources.Geode + s.Robots.Geode*(minute-1)
			}
		}

		for state := range states {
			// if the best this branch could possibly do is less than something another
			// branch can _definitely_ do, abandon it.
			if state.TheoreticalBest(minute, bp.Geode) < geodes {
				continue
			}

			// If we can build a geode bot, just do that.
			geode := state.Resources.Have(bp.Geode)
			if geode {
				next[state.Geode(bp.Geode)] = true
			}

			ore := state.Resources.Have(bp.Ore)
			if ore && state.Robots.Ore < max.Ore {
				next[state.Ore(bp.Ore)] = true
			}

			clay := state.Resources.Have(bp.Clay)
			if clay && state.Robots.Clay < max.Clay {
				next[state.Clay(bp.Clay)] = true
			}

			obsidian := state.Resources.Have(bp.Obsidian)
			if obsidian && state.Robots.Obsidian < max.Obsidian {
				next[state.Obsidian(bp.Obsidian)] = true
			}

			// if we could build _all_ of them, don't consider not building any of them
			if !(ore && clay && obsidian && geode) {
				next[state.Wait()] = true
			}
		}

		minute--
		states = next
	}

	var best uint16
	for state := range states {
		state.Resources = state.Resources.Plus(state.Robots)
		if state.Resources.Geode > best {
			best = state.Resources.Geode
		}
	}
	return best
}

func solution(name string, input []byte) int {
	// trim trailing space only
	input = bytes.Replace(input, []byte("\r"), []byte(""), -1)
	input = bytes.TrimRightFunc(input, unicode.IsSpace)
	lines := strings.Split(strings.TrimRightFunc(string(input), unicode.IsSpace), "\n")
	log.Printf("read %d %s lines", len(lines), name)

	// Blueprint 1: Each ore robot costs 4 ore. Each clay robot costs 3 ore. Each obsidian robot costs 3 ore and 18 clay.
	// Each geode robot costs 4 ore and 8 obsidian.

	blueprints := map[int]*Blueprint{}

	// 24 minutes
	for _, line := range lines {
		var id int
		bp := &Blueprint{}
		fmt.Sscanf(line, "Blueprint %d: Each ore robot costs %d ore. Each clay robot costs %d ore. Each "+
			"obsidian robot costs %d ore and %d clay. Each geode robot costs %d ore and %d obsidian",
			&id, &bp.Ore.Ore, &bp.Clay.Ore, &bp.Obsidian.Ore, &bp.Obsidian.Clay, &bp.Geode.Ore, &bp.Geode.Obsidian)
		blueprints[id] = bp
	}

	mu := &sync.Mutex{}
	var result int = 1
	ids := maps.Keys(blueprints)
	sort.Ints(ids)
	if len(ids) > 3 {
		ids = ids[0:3]
	}
	wg := &sync.WaitGroup{}
	for _, id := range ids {
		wg.Add(1)
		go func(id int) {
			bp := blueprints[id]
			bpg := bp.Geodes(32, State{
				Robots: Cost{Ore: 1},
			})
			mu.Lock()
			result *= int(bpg)
			log.Printf("%d yields %d", id, bpg)
			mu.Unlock()
			wg.Done()
		}(id)
	}
	wg.Wait()

	return int(result)
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

	input := aoc.Input(2022, 19)
	log.Printf("input solution: %d", solution("input", input))
}
