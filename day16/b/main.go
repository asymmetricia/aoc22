package main

import (
	"bytes"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strings"
	"unicode"

	"github.com/asymmetricia/aoc22/aoc"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/maps"
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

	for _, valve := range network {
		for _, peer := range valve.Peers {
			valve.Neighbors[peer] = network[peer]
		}
	}

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

	return -1
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
