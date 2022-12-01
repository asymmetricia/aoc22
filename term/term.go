package term

import (
	"image/color"
	"strconv"
)

func ClearLine() {
	print("\x1b[2K")
}

// Clear clears but does not move the cursor
func Clear() {
	print("\x1b[2J")
}

func MoveCursor(x, y int) {
	print("\x1b[", y, ";", x, "H")
}

func HideCursor() {
	print("\x1b[?25l")
}

func ShowCursor() {
	print("\x1b[?25h")
}

func ColorC(color color.Color) {
	print(ScolorC(color))
}

func ScolorC(color color.Color) string {
	r, g, b, _ := color.RGBA()
	return Scolor(int(r>>8), int(g>>8), int(b>>8))
}
func Color(r, g, b int) {
	print(Scolor(r, g, b))
}
func Scolor(r, g, b int) string {
	ret := []byte{
		'\x1b', '[', '3', '8', ';', '2', ';',
	}
	ret = append(ret, []byte(strconv.Itoa(r))...)
	ret = append(ret, ';')
	ret = append(ret, []byte(strconv.Itoa(g))...)
	ret = append(ret, ';')
	ret = append(ret, []byte(strconv.Itoa(b))...)
	ret = append(ret, 'm')
	return string(ret)
}

func ColorReset() {
	print(ScolorReset())
}

func ScolorReset() string {
	return "\x1b[0m"
}
