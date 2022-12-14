package main

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
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

	world := &coord.DenseWorld{}
	for _, line := range lines {
		fields := strings.Split(line, " -> ")
		var prev coord.Coord
		for i, field := range fields {
			c := coord.MustFromComma(field)
			if i == 0 {
				prev = c
				world.Set(prev, '#')
				continue
			}
			for prev.X < c.X {
				prev.X++
				world.Set(prev, '#')
			}
			for prev.X > c.X {
				prev.X--
				world.Set(prev, '#')
			}
			for prev.Y < c.Y {
				prev.Y++
				world.Set(prev, '#')
			}
			for prev.Y > c.Y {
				prev.Y--
				world.Set(prev, '#')
			}
		}
	}

	var sandStart = coord.MustFromComma("500,0")
	world.Set(sandStart, 'v')

	y := len(*world) + 1

	var worlds []*coord.DenseWorld
	particles := 0
	for {
		moved := false

		if world.At(sandStart) == 'V' {
			moved = true
			world.Set(sandStart, 'v')
		} else if world.At(sandStart) == 'v' {
			particles++
			moved = true
			world.Set(sandStart, '-')
		}

		for _, minus := range world.Find('-') {
			world.Set(minus, '+')
		}

		sands := world.Find('+')
		sort.Slice(sands, func(i, j int) bool {
			return sands[i].Y > sands[j].Y
		})
		for _, sand := range sands {
			if sand.South().Y < y {
				if world.At(sand.South()) == 0 {
					world.Set(sand, 0)
					world.Set(sand.South(), '-')
				} else if world.At(sand.SouthWest()) == 0 {
					world.Set(sand, 0)
					world.Set(sand.SouthWest(), '-')
				} else if world.At(sand.SouthEast()) == 0 {
					world.Set(sand, 0)
					world.Set(sand.SouthEast(), '-')
				}
			}
		}
		if world.At(sandStart) == 0 {
			world.Set(sandStart, 'V')
		}
		if particles%50 == 0 {
			worlds = append(worlds, world.Crop())
		}
		if !moved {
			world.Set(sandStart, 'v')
			worlds = append(worlds, world.Crop())
			break
		}
	}

	log.Print(particles)

	last := worlds[len(worlds)-1]
	vpos := last.Find('v')[0]
	_, _, maxx, maxy := last.Rect()

	worlds = append(worlds, worlds[len(worlds)-1])

	anim := &gif.GIF{
		Image:    make([]*image.Paletted, len(worlds)),
		Delay:    make([]int, len(worlds)),
		Disposal: make([]byte, len(worlds)),
	}
	const scale = 4
	wg := &sync.WaitGroup{}
	for i, world := range worlds {
		wg.Add(1)
		go func(i int, world *coord.DenseWorld) {
			defer wg.Done()
			img := image.NewPaletted(image.Rect(0, 0, maxx*scale, maxy*scale), aoc.TolVibrant)
			var wvpos coord.Coord
			if pos := world.Find('v'); len(pos) > 0 {
				wvpos = pos[0]
			} else if pos = world.Find('V'); len(pos) > 0 {
				wvpos = pos[0]
			}
			adj := vpos.Minus(wvpos)
			world.Each(func(c coord.Coord) (stop bool) {
				view := c.Plus(adj)
				var col color.Color = aoc.TolVibrantGrey
				switch world.At(c) {
				case 0:
					col = color.Black
				case '-':
					col = aoc.TolVibrantRed
				case '+':
					col = aoc.TolVibrantOrange
				default:
					col = aoc.TolVibrantGrey
				}
				draw.Draw(
					img,
					image.Rect(view.X*scale, view.Y*scale, (view.X+1)*scale, (view.Y+1)*scale),
					image.NewUniform(col),
					image.Pt(0, 0),
					draw.Over)
				return false
			})

			aoc.Typeset(img, image.Point{}, strconv.Itoa(len(world.Find('+'))), aoc.TolVibrantTeal, aoc.TypesetOpts{Scale: scale})

			anim.Image[i] = img
			anim.Delay[i] = 2
			anim.Disposal[i] = gif.DisposalNone
		}(i, world)
	}

	wg.Wait()

	anim.Delay[len(anim.Delay)-2] = 500

	aoc.SaveGIF(anim, "day14-"+name+"-b.gif", log)

	return particles
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

	input := aoc.Input(2022, 14)
	log.Printf("input solution: %d", solution("input", input))
}
