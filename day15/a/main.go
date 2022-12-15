package main

import (
	"bytes"
	"math"
	"os"
	"strings"
	"unicode"

	"github.com/sirupsen/logrus"

	"github.com/asymmetricia/aoc22/aoc"
	"github.com/asymmetricia/aoc22/coord"
)

var log = logrus.StandardLogger()

func solution(name string, input []byte) int {
	// trim trailing space only
	input = bytes.Replace(input, []byte("\r"), []byte(""), -1)
	input = bytes.TrimRightFunc(input, unicode.IsSpace)
	lines := strings.Split(strings.TrimRightFunc(string(input), unicode.IsSpace), "\n")
	log.Printf("read %d %s lines", len(lines), name)

	var targetY = 2000000
	if name == "test" {
		targetY = 10
	}
	sensors := map[coord.Coord]int{}
	beacons := map[coord.Coord]bool{}
	for _, line := range lines {
		fields := strings.Fields(line)
		xstr := aoc.Before(aoc.After(fields[2], "="), ",")
		ystr := aoc.Before(aoc.After(fields[3], "="), ":")
		clXStr := aoc.Before(aoc.After(fields[8], "="), ",")
		clYStr := aoc.After(fields[9], "=")
		//log.Print(line, " ", xstr, " ", ystr, " ", clXStr, " ", clYStr)
		sensor := coord.MustFromComma(xstr + "," + ystr)
		beacon := coord.MustFromComma(clXStr + "," + clYStr)
		sensors[sensor] = sensor.TaxiDistance(beacon)
		beacons[beacon] = true
	}

	minX := math.MaxInt
	maxX := math.MinInt
	for sensor, dist := range sensors {
		if sensor.X-dist < minX {
			minX = sensor.X - dist
		}
		if sensor.X+dist > maxX {
			maxX = sensor.X + dist
		}
	}

	log.Printf("%d to %d", minX, maxX)

	count := 0
xpos:
	for x := minX; x <= maxX; x++ {
		c := coord.C(x, targetY)
		if beacons[c] {
			continue
		}
		for sensor, dist := range sensors {
			if c.TaxiDistance(sensor) <= dist {
				count++
				continue xpos
			}
		}
	}

	return count
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

	input := aoc.Input(2022, 15)
	log.Printf("input solution: %d", solution("input", input))
}
