package main

import (
	"io/ioutil"
	"log"
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

	best := 0
	bestTotal := 0
	for i, calories := range elves {
		total := 0
		for _, c := range calories {
			total += c
		}
		if total > bestTotal {
			best = i
			bestTotal = total
		}
	}
	log.Print(best, bestTotal)
}
