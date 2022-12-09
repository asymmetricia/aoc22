package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
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

func frame(
	state string,
	trees [][]int,
	visible map[int]map[int]bool,
	cursorX int,
	cursorY int,
	bestX int,
	bestY int,
	bestScore int,
) *canvas.Canvas {
	ret := &canvas.Canvas{}

	canvas.TextBox{
		Title:  []rune("Dendrological Survey"),
		Width:  99,
		Height: 99,
		Footer: []rune("DENDRO v1.2"),
	}.On(ret)

	canvas.TextBox{
		Title: []rune("Status"),
		Left:  101,
		Width: 4 * aoc.GlyphWidth,
		Body:  []rune("[ ] Loading\n[ ] Perimeter Scan\n[ ] Seeking Best Tree\n[ ] Done"),
	}.On(ret)

	switch state {
	case "loading":
		ret.PrintAt(103, 1, "X", aoc.TolVibrantOrange)
		ret.PrintAt(106, 1, "Loading", aoc.TolVibrantOrange)
	case "perimeter":
		ret.PrintAt(103, 1, "✔", aoc.TolVibrantTeal)
		ret.PrintAt(106, 1, "Loading", aoc.TolVibrantTeal)
		ret.PrintAt(103, 2, "X", aoc.TolVibrantOrange)
		ret.PrintAt(106, 2, "Perimeter Scan", aoc.TolVibrantOrange)
	case "seeking":
		ret.PrintAt(103, 1, "✔", aoc.TolVibrantTeal)
		ret.PrintAt(106, 1, "Loading", aoc.TolVibrantTeal)
		ret.PrintAt(103, 2, "✔", aoc.TolVibrantTeal)
		ret.PrintAt(106, 2, "Perimeter Scan", aoc.TolVibrantTeal)
		ret.PrintAt(103, 3, "X", aoc.TolVibrantOrange)
		ret.PrintAt(106, 3, "Seeking Best Tree", aoc.TolVibrantOrange)
	case "done":
		ret.PrintAt(103, 1, "✔", aoc.TolVibrantTeal)
		ret.PrintAt(106, 1, "Loading", aoc.TolVibrantTeal)
		ret.PrintAt(103, 2, "✔", aoc.TolVibrantTeal)
		ret.PrintAt(106, 2, "Perimeter Scan", aoc.TolVibrantTeal)
		ret.PrintAt(103, 3, "✔", aoc.TolVibrantTeal)
		ret.PrintAt(106, 3, "Seeking Best Tree", aoc.TolVibrantTeal)
		ret.PrintAt(103, 4, "✔", aoc.TolVibrantTeal)
		ret.PrintAt(106, 4, "Done", aoc.TolVibrantTeal)
	}

	visCount := 0
	for yy, row := range trees {
		for xx, height := range row {
			c := aoc.TolVibrantGrey
			if visible != nil {
				if row, ok := visible[yy]; ok && row[xx] {
					c = aoc.TolVibrantTeal
					visCount++
				}
			}
			ret.PrintAt(xx+1, yy+1, strconv.Itoa(height), c)
		}
	}

	if cursorX > 0 && cursorY == -1 {
		for y := 0; y < 99; y++ {
			ret.PrintAt(cursorX+1, y+1, string(aoc.LineV), aoc.TolVibrantMagenta)
		}
	} else if cursorX == -1 && cursorY > 0 {
		for x := 0; x < 99; x++ {
			ret.PrintAt(x+1, cursorY+1, string(aoc.LineH), aoc.TolVibrantMagenta)
		}
	} else if cursorX > 0 && cursorY > 0 {
		canvas.TextBox{
			Top:       cursorY,
			Left:      cursorX,
			Body:      []rune(fmt.Sprintf("%d", trees[cursorY][cursorX])),
			BodyColor: aoc.TolVibrantRed,
		}.On(ret)
	}

	vtColor := aoc.TolVibrantGrey
	if state == "perimeter" {
		vtColor = aoc.TolVibrantOrange
	}
	canvas.TextBox{
		Top:        6,
		Left:       101,
		Title:      []rune("Visible Trees"),
		Body:       []rune(fmt.Sprintf("%4d", visCount)),
		BodyBlock:  true,
		FrameColor: vtColor,
	}.On(ret)

	bestFrameColor := aoc.TolVibrantGrey
	bestXDisp := " -- "
	if bestX >= 0 {
		bestFrameColor = aoc.TolVibrantOrange
		bestXDisp = fmt.Sprintf(" %2d ", bestX)
	}
	canvas.TextBox{
		Top:        6 + 2 + 8,
		Left:       101,
		Title:      []rune("Best Tree X"),
		Body:       []rune(bestXDisp),
		BodyBlock:  true,
		FrameColor: bestFrameColor,
	}.On(ret)

	bestYDisp := " -- "
	if bestY >= 0 {
		bestYDisp = fmt.Sprintf(" %2d ", bestY)
	}
	canvas.TextBox{
		Top:        6 + 2 + 8 + 2 + 8,
		Left:       101,
		Title:      []rune("Best Tree Y"),
		Body:       []rune(bestYDisp),
		BodyBlock:  true,
		FrameColor: bestFrameColor,
	}.On(ret)

	bestScoreDisp := "------"
	if bestScore >= 0 {
		bestScoreDisp = fmt.Sprintf("%6d", bestScore)

		canvas.TextBox{
			Top:       bestY,
			Left:      bestX,
			Body:      []rune(fmt.Sprintf("%d", trees[bestY][bestX])),
			BodyColor: aoc.TolVibrantCyan,
		}.On(ret)
	}

	padding := (4*aoc.GlyphWidth - len(bestScoreDisp)) / 2
	for i := 0; i < padding; i++ {
		bestScoreDisp = " " + bestScoreDisp
	}

	canvas.TextBox{
		Top:        6 + 2 + 8 + 2 + 8 + 2 + 8,
		Left:       101,
		Title:      []rune("Best Tree Score"),
		Body:       []rune(bestScoreDisp),
		Width:      4 * aoc.GlyphWidth,
		FrameColor: bestFrameColor,
	}.On(ret)

	if state == "done" {
		canvas.TextBox{
			Top:        40,
			Center:     true,
			Title:      []rune("Best Tree Score"),
			Body:       []rune(strings.TrimSpace(bestScoreDisp)),
			BodyBlock:  true,
			FrameColor: aoc.TolVibrantTeal,
			TitleColor: aoc.TolVibrantMagenta,
			BodyColor:  aoc.TolVibrantCyan,
		}.On(ret)
		canvas.TextBox{
			Top:        50,
			Center:     true,
			Title:      []rune("Perimeter-Visible Tree Count"),
			Body:       []rune(strconv.Itoa(visCount)),
			BodyBlock:  true,
			FrameColor: aoc.TolVibrantTeal,
			TitleColor: aoc.TolVibrantMagenta,
			BodyColor:  aoc.TolVibrantCyan,
		}.On(ret)
	}

	return ret
}

func perimeterScan(trees [][]int) ([]*canvas.Canvas, map[int]map[int]bool) {
	var ret []*canvas.Canvas
	visible := map[int]map[int]bool{}
	for y, row := range trees {
		visible[y] = map[int]bool{}
		visible[y][0] = true
		last := row[0]
		for x, height := range row {
			if height > last {
				last = height
				visible[y][x] = true
			}
		}

		visible[y][len(row)-1] = true
		last = row[len(row)-1]
		for x := len(row) - 1; x >= 0; x-- {
			if row[x] > last {
				last = row[x]
				visible[y][x] = true
			}
		}

		ret = append(ret, frame("perimeter", trees, visible, -1, y, -1, -1, -1))
	}

	for x := 0; x < len(trees[0]); x++ {
		visible[0][x] = true
		last := trees[0][x]
		for y, row := range trees {
			if row[x] > last {
				last = row[x]
				visible[y][x] = true
			}
		}

		visible[len(trees)-1][x] = true
		last = trees[len(trees)-1][x]
		for y := len(trees) - 1; y >= 0; y-- {
			if trees[y][x] > last {
				last = trees[y][x]
				visible[y][x] = true
			}
		}
		ret = append(ret, frame("perimeter", trees, visible, x, -1, -1, -1, -1))
	}

	return ret, visible
}

func solution(name string, input []byte) int {
	// trim trailing space only
	input = bytes.Replace(input, []byte("\r"), []byte(""), -1)
	input = bytes.TrimRightFunc(input, unicode.IsSpace)
	lines := strings.Split(strings.TrimRightFunc(string(input), unicode.IsSpace), "\n")
	log.Printf("read %d %s lines", len(lines), name)

	var frames []*canvas.Canvas

	trees := [][]int{}
	for y, row := range lines {
		trees = append(trees, make([]int, len(row)))
		for x, height := range row {
			trees[y][x] = int(height - '0')
		}
		frames = append(frames, frame("loading", trees, nil, -1, -1, -1, -1, -1))
	}

	perimFrames, perimTrees := perimeterScan(trees)
	frames = append(frames, perimFrames...)

	var best coord.Coord
	var bestScore int = math.MinInt
	for y, row := range trees {
		for x := range row {
			visibleTrees := visible(trees, x, y)
			score := len(visibleTrees[coord.North]) *
				len(visibleTrees[coord.East]) *
				len(visibleTrees[coord.South]) *
				len(visibleTrees[coord.West])
			if score > bestScore {
				best = coord.C(x, y)
				bestScore = score
			}
			if rand.Intn(20) == 0 {
				frames = append(frames, frame("seeking", trees, perimTrees, x, y, best.X, best.Y, bestScore))
			}
		}
		if y%10 == 0 {
			log.Printf("y=%d", y)
		}
	}

	frames = append(frames, frame("done", trees, perimTrees, -1, -1, best.X, best.Y, bestScore))

	log.Print("computation finished")

	canvas.RenderGif(frames, nil, "day08b-"+name+".gif", log)

	return bestScore
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
