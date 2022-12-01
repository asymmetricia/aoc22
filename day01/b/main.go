package main

import (
	"io/ioutil"
	"log"
	"sort"
	"strconv"
	"strings"
)

func main() {
	in, err := ioutil.ReadFile("input")
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(in), "\n")

	var elves [][]int
	var accum []int
	for _, line := range lines {
		if line == "" {
			elves = append(elves, accum)
			accum = nil
		} else {
			i, err := strconv.Atoi(line)
			if err != nil {
				panic(err)
			}
			accum = append(accum, i)
		}
	}
	if len(accum) > 0 {
		elves = append(elves, accum)
	}

	var counts []int
	for _, calories := range elves {
		total := 0
		for _, c := range calories {
			total += c
		}
		counts = append(counts, total)
	}
	sort.Ints(counts)
	log.Print(counts[len(counts)-3] + counts[len(counts)-2] + counts[len(counts)-1])
}
