package main

import (
	"bytes"
	"fmt"
	"os"
	"sort"
	"strings"
	"unicode"

	"github.com/sirupsen/logrus"
	"golang.org/x/exp/maps"

	"github.com/asymmetricia/aoc22/aoc"
)

var log = logrus.StandardLogger()

type Cost struct {
	Ore, Clay, Obsidian int
}

func (c Cost) Have(b Cost) bool {
	return c.Ore >= b.Ore && c.Clay >= b.Clay && c.Obsidian >= b.Obsidian
}

func (c Cost) Minus(b Cost) Cost {
	return Cost{c.Ore - b.Ore, c.Clay - b.Clay, c.Obsidian - b.Obsidian}
}

func (c Cost) Plus(b Cost) Cost {
	return Cost{c.Ore + b.Ore, c.Clay + b.Clay, c.Obsidian + b.Obsidian}
}

type cacheKey struct {
	Minutes, Ore, Clay, Obsidian, Geode int
	Resources                           Cost
}

type Blueprint struct {
	Ore      Cost
	Clay     Cost
	Obsidian Cost
	Geode    Cost

	cache map[cacheKey]int
}

func (bp *Blueprint) Geodes(minutes int, ore, clay, obsidian, geode int, resources Cost) int {
	if bp.cache == nil {
		bp.cache = map[cacheKey]int{}
	}

	ck := cacheKey{minutes, ore, clay, obsidian, geode, resources}
	if cv, ok := bp.cache[ck]; ok {
		return cv
	}

	if minutes == 1 {
		return geode
	}

	best := 0

	if resources.Have(bp.Ore) {
		c := bp.Geodes(minutes-1, ore+1, clay, obsidian, geode, resources.Minus(bp.Ore).Plus(Cost{ore, clay, obsidian}))
		if c > best {
			best = c
		}
	}
	if resources.Have(bp.Clay) {
		c := bp.Geodes(minutes-1, ore, clay+1, obsidian, geode, resources.Minus(bp.Clay).Plus(Cost{ore, clay, obsidian}))
		if c > best {
			best = c
		}
	}
	if resources.Have(bp.Obsidian) {
		c := bp.Geodes(minutes-1, ore, clay, obsidian+1, geode, resources.Minus(bp.Obsidian).Plus(Cost{ore, clay, obsidian}))
		if c > best {
			best = c
		}
	}
	if resources.Have(bp.Geode) {
		c := bp.Geodes(minutes-1, ore, clay, obsidian, geode+1, resources.Minus(bp.Geode).Plus(Cost{ore, clay, obsidian}))
		if c > best {
			best = c
		}
	}

	c := bp.Geodes(minutes-1, ore, clay, obsidian, geode, resources.Plus(Cost{ore, clay, obsidian}))
	if c > best {
		best = c
	}

	bp.cache[ck] = best + geode
	return best + geode
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
		log.Printf("%+v", bp)
		//fields := strings.Fields(line)
	}

	var qlsum int
	var best, bestId int
	ids := maps.Keys(blueprints)
	sort.Ints(ids)
	for _, id := range ids {
		bp := blueprints[id]
		bpg := bp.Geodes(24, 1, 0, 0, 0, Cost{})
		bp.cache = nil
		qlsum += id * bpg
		log.Printf("%d yields %d", id, bpg)
		if bpg > best {
			best = bpg
			bestId = id
		}
	}

	log.Printf("%d -> %d", bestId, best)

	return qlsum
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
