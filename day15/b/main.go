package main

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"strings"
	"unicode"

	"github.com/sirupsen/logrus"

	"github.com/asymmetricia/aoc22/aoc"
	"github.com/asymmetricia/aoc22/coord"
)

var log = logrus.StandardLogger()

func frame(sensors map[coord.Coord]int, beacons map[coord.Coord]bool, answer coord.Coord) *image.Paletted {
	const reduction = 2000
	ret := image.NewPaletted(image.Rect(0, 0, 2000, 2000), aoc.TolVibrant)
	dot := func(x, y int, color color.RGBA, scale int) {
		dotSize := 5 * scale
		if x < 0 || x > 2000 || y < 0 || y > 2000 {
			return
		}
		draw.Draw(ret, image.Rect(x-dotSize, y-dotSize, x+dotSize, y+dotSize), image.NewUniform(color), image.Point{}, draw.Over)
	}
	for sensor, dist := range sensors {
		x := sensor.X / reduction
		y := sensor.Y / reduction
		for _, point := range sensor.TaxiPerimeter(dist + 1) {
			ret.Set(point.X/reduction, point.Y/reduction, aoc.TolVibrantBlue)
		}
		dot(x, y, aoc.TolVibrantTeal, 1)
	}
	for beacon := range beacons {
		dot(beacon.X/reduction, beacon.Y/reduction, aoc.TolVibrantMagenta, 1)
	}
	dot(answer.X/reduction, answer.Y/reduction, aoc.TolVibrantCyan, 2)
	return ret
}

func solution(name string, input []byte) int {
	// trim trailing space only
	input = bytes.Replace(input, []byte("\r"), []byte(""), -1)
	input = bytes.TrimRightFunc(input, unicode.IsSpace)
	lines := strings.Split(strings.TrimRightFunc(string(input), unicode.IsSpace), "\n")
	log.Printf("read %d %s lines", len(lines), name)

	sensors := map[coord.Coord]int{}
	sensorBound := map[coord.Coord][2]coord.Coord{}
	beacons := map[coord.Coord]bool{}
	for _, line := range lines {
		fields := strings.Fields(line)
		xstr := aoc.Before(aoc.After(fields[2], "="), ",")
		ystr := aoc.Before(aoc.After(fields[3], "="), ":")
		clXStr := aoc.Before(aoc.After(fields[8], "="), ",")
		clYStr := aoc.After(fields[9], "=")
		sensor := coord.MustFromComma(xstr + "," + ystr)
		beacon := coord.MustFromComma(clXStr + "," + clYStr)
		dist := sensor.TaxiDistance(beacon)
		sensors[sensor] = dist
		beacons[beacon] = true
		sensorBound[sensor] = [2]coord.Coord{{sensor.X - dist, sensor.Y - dist}, {sensor.X + dist, sensor.Y + dist}}
	}

	minX := 0
	maxX := 4000000
	minY := 0
	maxY := 4000000
	if name == "test" {
		maxX = 20
		maxY = 20
	}

	var testCoord = func(c coord.Coord) bool {
		for sensor, dist := range sensors {
			if sensor.TaxiDistance(c) <= dist {
				return false
			}
		}
		return true
	}

	var answer coord.Coord
	dist := 0
outer:
	for {
		dist++
		if dist%1000 == 0 {
			log.Print(dist)
		}
		for sensor, beaconDist := range sensors {
			for _, c := range sensor.TaxiPerimeter(beaconDist + dist) {
				if c.X >= minX && c.X <= maxX && c.Y >= minY && c.Y <= maxY && testCoord(c) {
					log.Printf("found %v at dist=+%d", c, dist)
					answer = c
					break outer
				}
			}
		}
	}

	img := frame(sensors, beacons, answer)
	f, err := os.OpenFile("day15-b-"+name+".png", os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	png.Encode(f, img)
	f.Sync()
	f.Close()

	return answer.X*4000000 + answer.Y
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
