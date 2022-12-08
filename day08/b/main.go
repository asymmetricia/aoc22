package main

import (
	"bytes"
	"image/color"
	"math"
	"os"
	"strconv"
	"strings"
	"unicode"

	"github.com/sirupsen/logrus"

	"github.com/asymmetricia/aoc22/aoc"
	"github.com/asymmetricia/aoc22/canvas"
	"github.com/asymmetricia/aoc22/coord"
)

var log = logrus.StandardLogger()

func visible(trees [][]int, x, y int) map[coord.Direction][]coord.Coord {
	ret := map[coord.Direction][]coord.Coord{}

	height := trees[y][x]

	for north := y - 1; north >= 0; north-- {
		ret[coord.North] = append(ret[coord.North], coord.C(x, north))
		if trees[north][x] >= height {
			break
		}
	}

	for south := y + 1; south < len(trees); south++ {
		ret[coord.South] = append(ret[coord.South], coord.C(x, south))
		if trees[south][x] >= height {
			break
		}
	}

	for west := x - 1; west >= 0; west-- {
		ret[coord.West] = append(ret[coord.West], coord.C(west, y))
		if trees[y][west] >= height {
			break
		}
	}

	for east := x + 1; east < len(trees[0]); east++ {
		ret[coord.East] = append(ret[coord.East], coord.C(east, y))
		if trees[y][east] >= height {
			break
		}
	}

	return ret
}

func frame(vis map[int]map[int]int, trees [][]int, x, y int, v map[coord.Direction][]coord.Coord, best coord.Coord) *canvas.Canvas {
	ret := &canvas.Canvas{}

	for yy, row := range trees {
		for xx, height := range row {
			var c color.Color = color.White
			if yy < y || yy == y && xx < x {
				c = aoc.TolVibrantGrey
			}
			if xx == x && yy == y {
				c = aoc.TolVibrantOrange
			}
			ret.PrintAt(xx*2, yy, strconv.Itoa(height), c)
		}
	}

	for _, visibleTrees := range v {
		for _, tree := range visibleTrees {
			xx := tree.X
			yy := tree.Y
			ret.PrintAt(xx*2, yy, strconv.Itoa(trees[yy][xx]), aoc.TolVibrantCyan)
		}
	}

	if best.Y > 0 && best.X > 0 {
		canvas.TextBox{
			Top:       best.Y - 1,
			Left:      best.X*2 - 1,
			Body:      []rune("+"),
			BodyColor: aoc.TolVibrantRed,
		}.On(ret)
	}

	return ret
}

func solution(name string, input []byte) int {
	// trim trailing space only
	input = bytes.Replace(input, []byte("\r"), []byte(""), -1)
	input = bytes.TrimRightFunc(input, unicode.IsSpace)
	lines := strings.Split(strings.TrimRightFunc(string(input), unicode.IsSpace), "\n")
	log.Printf("read %d %s lines", len(lines), name)

	trees := [][]int{}
	for y, row := range lines {
		trees = append(trees, make([]int, len(row)))
		for x, height := range row {
			trees[y][x] = int(height - '0')
		}
	}

	framesMap := map[int]map[int]*canvas.Canvas{}
	visibility := map[int]map[int]int{}

	var best coord.Coord
	for y, row := range trees {
		visibility[y] = map[int]int{}
		framesRow := map[int]*canvas.Canvas{}
		for x := range row {
			visibleTrees := visible(trees, x, y)
			score := len(visibleTrees[coord.North]) *
				len(visibleTrees[coord.East]) *
				len(visibleTrees[coord.South]) *
				len(visibleTrees[coord.West])
			visibility[y][x] = score
			if score > visibility[best.Y][best.X] {
				log.Print(score)
				best = coord.C(x, y)
			}
			if x > 0 && x < len(row)-1 && y > 0 && y < len(trees)-1 {
				framesRow[x] = frame(visibility, trees, x, y, visibleTrees, best)
			}
		}
		framesMap[y] = framesRow
		if y%10 == 0 {
			log.Printf("y=%d", y)
		}
	}

	log.Print("computation finished")

	var frames []*canvas.Canvas
	for y, row := range trees {
		for x := range row {
			if framesMap[y][x] != nil {
				frames = append(frames, framesMap[y][x])
			}
		}
	}
	canvas.RenderGif(frames, map[int]float32{10: 50, 20: 25, 40: 12, 150: 3, 300: 1.0 / 2, math.MaxInt: 1.0 / 4}, "day08b-"+name+".gif", log)

	return visibility[best.Y][best.X]
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

	input, err := os.ReadFile("input")
	if err != nil {
		log.WithError(err).Fatal("could not read input")
	}
	log.Printf("input solution: %d", solution("input", input))
}
