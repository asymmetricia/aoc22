package canvas

import (
	"image/gif"
	"math"

	"github.com/sirupsen/logrus"

	"github.com/asymmetricia/aoc22/aoc"
)

func timing(n int, csPerFrame map[int]float32) float32 {
	// the best timing is the timing T such that there is no timing less than T and greater than N.
	best := math.MaxInt
	for t := range csPerFrame {
		if t > n && t < best {
			best = t
		}
	}

	v, ok := csPerFrame[best]
	if !ok {
		return 3
	}

	return v
}

func RenderGif(canvases []*Canvas, csPerFrame map[int]float32, filename string, log logrus.FieldLogger) {
	width, height := Bounds(canvases)
	anim := &gif.GIF{}
	var i int
	var canvas *Canvas
	var frameTiming float32
	for i, canvas = range canvases {
		if csPerFrame == nil {
			frameTiming = 3
		} else {
			frameTiming = timing(i, csPerFrame)
			if frameTiming < 1 {
				if i%(int(1.0/frameTiming)) != 0 && i != len(canvases)-1 {
					continue
				}
				frameTiming = 3
			}
		}

		anim.Image = append(anim.Image, canvas.RenderRect(width, height))
		anim.Delay = append(anim.Delay, int(frameTiming))
		anim.Disposal = append(anim.Disposal, gif.DisposalNone)

	}

	anim.Image = append(anim.Image, canvas.RenderRect(width, height))
	anim.Delay = append(anim.Delay, int(frameTiming))
	anim.Disposal = append(anim.Disposal, gif.DisposalNone)

	anim.Delay[len(anim.Delay)-2] = 500

	aoc.SaveGIF(anim, filename, log)
}
