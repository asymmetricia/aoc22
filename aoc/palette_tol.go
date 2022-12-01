package aoc

import (
	"golang.org/x/exp/constraints"
	"image/color"
)

// TolVibrant is Paul Tol's Vibrant qualitative color palette that should be
// color-blind accessible; see https://personal.sron.nl/~pault/
var TolVibrant = color.Palette{
	color.Black,
	color.Transparent,
	color.White,
	TolVibrantBlue,
	TolVibrantCyan,
	TolVibrantTeal,
	TolVibrantOrange,
	TolVibrantRed,
	TolVibrantMagenta,
	TolVibrantGrey,
}

var (
	TolVibrantBlue    = color.RGBA{0, 119, 187, 255}
	TolVibrantCyan    = color.RGBA{51, 187, 238, 255}
	TolVibrantTeal    = color.RGBA{0, 153, 136, 255}
	TolVibrantOrange  = color.RGBA{238, 119, 51, 255}
	TolVibrantRed     = color.RGBA{204, 51, 17, 255}
	TolVibrantMagenta = color.RGBA{238, 51, 119, 255}
	TolVibrantGrey    = color.RGBA{187, 187, 187, 255}
)

// TolSequentialSmoothRainbow is Paul Tol's "smooth rainbow" sequential color palette.
var TolSequentialSmoothRainbow = color.Palette{
	color.Black,
	color.Transparent,
	color.White,
	color.RGBA{232, 236, 251, 255},
	color.RGBA{221, 216, 239, 255},
	color.RGBA{209, 193, 225, 255},
	color.RGBA{195, 168, 209, 255},
	color.RGBA{181, 143, 194, 255},
	color.RGBA{167, 120, 180, 255},
	color.RGBA{155, 98, 167, 255},
	color.RGBA{140, 78, 153, 255},
	color.RGBA{111, 76, 155, 255},
	color.RGBA{96, 89, 169, 255},
	color.RGBA{85, 104, 184, 255},
	color.RGBA{78, 121, 197, 255},
	color.RGBA{77, 138, 198, 255},
	color.RGBA{78, 150, 188, 255},
	color.RGBA{84, 158, 179, 255},
	color.RGBA{89, 165, 169, 255},
	color.RGBA{96, 171, 158, 255},
	color.RGBA{105, 177, 144, 255},
	color.RGBA{119, 183, 125, 255},
	color.RGBA{140, 188, 104, 255},
	color.RGBA{166, 190, 84, 255},
	color.RGBA{190, 188, 72, 255},
	color.RGBA{209, 181, 65, 255},
	color.RGBA{221, 170, 60, 255},
	color.RGBA{228, 156, 57, 255},
	color.RGBA{231, 140, 53, 255},
	color.RGBA{230, 121, 50, 255},
	color.RGBA{228, 99, 45, 255},
	color.RGBA{223, 72, 40, 255},
	color.RGBA{218, 34, 34, 255},
	color.RGBA{184, 34, 30, 255},
	color.RGBA{149, 33, 27, 255},
	color.RGBA{114, 30, 23, 255},
	color.RGBA{82, 26, 19, 255},
}

// TolScale returns a sequential rainbow colorfrom TolSequentialSmoothRainbow for
// the given value, scaled to min and max. Out of bound values (< min or > max)
// are clamped.
func TolScale[K constraints.Integer | constraints.Float](min, max, val K) color.RGBA {
	scale := max - min
	adj := val - min
	// 34 possible colors
	// i should go from 0 ... 33.99999, so when we truncate down, we end up with 0..33
	// we can't guarantee we don't get 34 exactly, but we clamp after truncating so it's OK
	i := int(float32(adj) / float32(scale) * 34)
	if i < 0 {
		i = 0
	}
	if i > 33 {
		i = 33
	}
	return TolSequentialSmoothRainbow[3+i].(color.RGBA)
}
