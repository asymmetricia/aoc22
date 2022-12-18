package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"strings"
	"sync"
	"unicode"

	"github.com/sirupsen/logrus"
	"golang.org/x/exp/maps"

	"github.com/asymmetricia/aoc22/aoc"
	"github.com/asymmetricia/aoc22/isovox"
)

var log = logrus.StandardLogger()

type coord struct {
	x, y, z int
}

func (c coord) ivc() isovox.Coord {
	return isovox.Coord{c.x, c.y, c.z}
}

func bounds(world map[coord]int, f ...func(coord, int) bool) (min, max coord) {
	min = coord{math.MaxInt, math.MaxInt, math.MaxInt}
	max = coord{math.MinInt, math.MinInt, math.MinInt}
	for c, v := range world {
		if len(f) > 0 && !f[0](c, v) {
			continue
		}

		min.x = aoc.Min(min.x, c.x)
		min.y = aoc.Min(min.y, c.y)
		min.z = aoc.Min(min.z, c.z)
		max.x = aoc.Max(max.x, c.x)
		max.y = aoc.Max(max.y, c.y)
		max.z = aoc.Max(max.z, c.z)
	}
	return
}

func render(world, extra map[coord]int) image.Image {
	world = maps.Clone(world)
	maps.Copy(world, extra)
	min, max := bounds(world, func(c coord, v int) bool {
		return v == 1
	})
	ivw := &isovox.World{Voxels: map[isovox.Coord]*isovox.Voxel{}}

	for _, x := range []int{min.x - 1, max.x + 1} {
		for _, y := range []int{min.y - 1, max.y + 1} {
			for _, z := range []int{min.z - 1, max.z + 1} {
				ivw.Voxels[isovox.Coord{x, y, z}] = &isovox.Voxel{Color: color.Transparent}
			}
		}
	}

	cyan := color.NRGBA{51, 187, 238, 0x10}
	for c, v := range world {
		switch v {
		case 1:
			ivw.Voxels[c.ivc()] = &isovox.Voxel{Color: aoc.TolVibrantRed}
		case 2:
			ivw.Voxels[c.ivc()] = &isovox.Voxel{Color: cyan}
		case 3:
			ivw.Voxels[c.ivc()] = &isovox.Voxel{Color: aoc.TolVibrantOrange}
		}
	}

	return ivw.Render(24)
}

func solution(name string, input []byte) int {
	ich := make(chan image.Image)
	ffmpegErr := aoc.RenderMP4(ich, "day18-b-"+name+".mp4", 60, log)
	ffmpegWg := &sync.WaitGroup{}
	ffmpegWg.Add(1)
	defer ffmpegWg.Wait()
	go func(ffmpegErr <-chan error) {
		defer ffmpegWg.Done()
		for err := range ffmpegErr {
			log.Error(err)
		}
	}(ffmpegErr)

	// trim trailing space only
	input = bytes.Replace(input, []byte("\r"), []byte(""), -1)
	input = bytes.TrimRightFunc(input, unicode.IsSpace)
	lines := strings.Split(strings.TrimRightFunc(string(input), unicode.IsSpace), "\n")
	log.Printf("read %d %s lines", len(lines), name)

	world := map[coord]int{}
	for _, line := range lines {
		var x, y, z int
		fmt.Sscanf(line, "%d,%d,%d", &x, &y, &z)
		world[coord{x, y, z}] = 1
	}

	var min, max coord
	for c := range world {
		if (min == max && min == coord{}) {
			min = c
			max = c
		}
		if c.x < min.x {
			min.x = c.x
		}
		if c.y < min.y {
			min.y = c.y
		}
		if c.z < min.z {
			min.z = c.z
		}
		if c.x > max.x {
			max.x = c.x
		}
		if c.y > max.y {
			max.y = c.y
		}
		if c.z > max.z {
			max.z = c.z
		}
	}
	min.x--
	min.y--
	min.z--
	max.x++
	max.y++
	max.z++

	log.Print(min, max)
	world[min] = 2
	changed := true
	water := 1
	for changed {
		changed = false
		toAdd := map[coord]int{}
		for c, v := range world {
			if v == 2 {
				x, y, z := c.x, c.y, c.z
				for _, n := range []coord{
					{x + 1, y, z}, {x - 1, y, z},
					{x, y + 1, z}, {x, y - 1, z},
					{x, y, z + 1}, {x, y, z - 1},
				} {
					if n.x < min.x || n.x > max.x ||
						n.y < min.y || n.y > max.y ||
						n.z < min.z || n.z > max.z {
						continue
					}
					if world[n] == 0 && toAdd[n] == 0 {
						water++
						changed = true
						toAdd[n] = 2
						if water < 10 ||
							water < 1000 && water%5 == 0 ||
							water%50 == 0 {
							ich <- render(world, toAdd)
						}
					}
				}
			}
		}
		log.Print(water)
		maps.Copy(world, toAdd)
	}

	surfaces := 0
	for z := min.z; z <= max.z; z++ {
		exp := false
		for x := min.x; x <= max.x; x++ {
			for y := min.y; y <= max.y; y++ {
				if world[coord{x, y, z}] != 1 {
					continue
				}
				neighbors := []coord{
					{x + 1, y, z}, {x - 1, y, z},
					{x, y + 1, z}, {x, y - 1, z},
					{x, y, z + 1}, {x, y, z - 1},
				}
				for _, neighbor := range neighbors {
					if world[neighbor] == 2 {
						world[coord{x, y, z}] = 3
						exp = true
						surfaces++
					}
				}
			}
		}
		log.Printf("counting surfaces, z=%d", z)
		if exp {
			ich <- render(world, nil)
		}
	}

	close(ich)

	f, err := os.OpenFile("day18-b-"+name+".png", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err == nil {
		err = png.Encode(f, render(world, nil))
	}
	if err == nil {
		err = f.Sync()
	}
	if err == nil {
		err = f.Close()
	}

	return surfaces
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

	input := aoc.Input(2022, 18)
	log.Printf("input solution: %d", solution("input", input))
}
