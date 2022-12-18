package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"os"
	"sort"
	"strings"
	"unicode"

	"github.com/asymmetricia/aoc22/aoc"
	"github.com/fogleman/fauxgl"
	"github.com/nfnt/resize"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/maps"
)

var log = logrus.StandardLogger()

type coord struct {
	x, y, z int
}

const (
	scale  = 2   // optional supersampling
	width  = 800 // output width in pixels
	height = 800 // output height in pixels
	fovy   = 30  // vertical field of view in degrees
	near   = 1   // near clipping plane
	far    = 300 // far clipping plane
)

var (
	eye    = fauxgl.V(50, 60, -100)              // camera position
	center = fauxgl.V(15, 15, 20)                // view center position
	up     = fauxgl.V(0, 1, 0)                   // up vector
	light  = fauxgl.V(0.75, 1, -.25).Normalize() // light direction
)

func render(world map[coord]int) image.Image {
	rctx := fauxgl.NewContext(width*scale, height*scale)
	rctx.ClearColorBufferWith(fauxgl.Black)

	aspect := float64(width) / float64(height)
	matrix := fauxgl.LookAt(eye, center, up).Perspective(fovy, aspect, near, far)
	shader := fauxgl.NewPhongShader(matrix, light, eye)
	rctx.Shader = shader

	red := fauxgl.MakeColor(aoc.TolVibrantRed)
	blue := fauxgl.MakeColor(aoc.TolVibrantCyan).Alpha(0.5)
	for c, v := range world {
		voxel := fauxgl.Voxel{
			X: c.x * 2,
			Y: c.z * 2,
			Z: c.y * 2,
		}
		if v == 2 {
			voxel.Color = blue
		} else {
			voxel.Color = red
		}
		rctx.DrawMesh(fauxgl.NewVoxelMesh([]fauxgl.Voxel{voxel}))
	}

	return resize.Resize(width, height, rctx.Image(), resize.Bilinear)
}

func solution(name string, input []byte) int {
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

	var images []image.Image
	log.Print(min, max)
	world[min] = 2
	changed := true
	water := 0
	for changed {
		changed = false
		for z := min.z; z <= max.z; z++ {
			for x := min.x; x <= max.x; x++ {
			cube:
				for y := min.y; y <= max.y; y++ {
					// 0123456789
					// 2222222210
					v := world[coord{x, y, z}]
					if v > 0 {
						continue
					}
					neighbors := []coord{
						coord{x, y + 1, z},
						coord{x, y - 1, z},
						coord{x + 1, y, z},
						coord{x - 1, y, z},
						coord{x, y, z + 1},
						coord{x, y, z - 1},
					}
					for _, neighbor := range neighbors {
						n := world[neighbor]
						if n == 2 {
							water++
							changed = true
							world[coord{x, y, z}] = 2
							if water%100 == 0 {
								log.Print(water)
								images = append(images, render(world))
							}
							continue cube
						}
					}
				}
			}
		}
	}

	fmt.Printf("z=%d\n", (min.z+max.z)/2)
	for y := min.y; y <= max.y; y++ {
		fmt.Printf("%3d ", y)
		for x := min.x; x <= max.x; x++ {
			print(world[coord{x, y, (min.z + max.z) / 2}])
		}
		println()
	}

	surfaces := 0
	for x := min.x; x <= max.x; x++ {
		for y := min.y; y <= max.y; y++ {
			for z := min.z; z <= max.z; z++ {
				if world[coord{x, y, z}] != 1 {
					continue
				}
				neighbors := []coord{
					coord{x, y + 1, z},
					coord{x, y - 1, z},
					coord{x + 1, y, z},
					coord{x - 1, y, z},
					coord{x, y, z + 1},
					coord{x, y, z - 1},
				}
				for _, neighbor := range neighbors {
					if world[neighbor] == 2 {
						surfaces++
					}
				}
			}
		}
	}

	images = append(images, render(world))
	images = append(images, render(world))

	palette := color.Palette{
		color.Black,
		color.Transparent,
		color.White,
	}

	colors := map[color.RGBA64]int{}
	for _, img := range images {
		rect := img.Bounds()
		min, max := rect.Min, rect.Max
		for x := min.X; x <= max.X; x++ {
			for y := min.Y; y <= max.Y; y++ {
				r, g, b, a := img.At(x, y).RGBA()
				colors[color.RGBA64{
					R: uint16(r),
					G: uint16(g),
					B: uint16(b),
					A: uint16(a),
				}]++
			}
		}
	}
	log.Print(len(colors), " colors")
	colorList := maps.Keys(colors)
	sort.Slice(colorList, func(i, j int) bool {
		return colors[colorList[j]] < colors[colorList[i]]
	})
	for len(palette) < 255 && len(colorList) > 0 {
		palette = append(palette, colorList[0])
		colorList = colorList[1:]
	}

	anim := &gif.GIF{}
	for _, img := range images {
		pimg := image.NewPaletted(img.Bounds(), palette)
		draw.FloydSteinberg.Draw(pimg, img.Bounds(), img, image.Point{})
		anim.Image = append(anim.Image, pimg)
	}
	anim.Delay = make([]int, len(anim.Image))
	anim.Disposal = make([]byte, len(anim.Image))
	for i := range anim.Image {
		anim.Delay[i] = 2
		anim.Disposal[i] = gif.DisposalNone
	}
	anim.Delay[len(anim.Delay)-2] = 300

	aoc.SaveGIF(anim, "day18-b-"+name+".gif", log)

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
