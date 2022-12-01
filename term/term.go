package term

import (
	"fmt"
	"image/color"
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
	r, g, b, _ := color.RGBA()
	Color(int(r>>8), int(g>>8), int(b>>8))
}

func Color(r, g, b int) {
	print(Scolor(r, g, b))
}
func Scolor(r, g, b int) string {
	return fmt.Sprint("\x1b[38;2;", r, ";", g, ";", b, "m")
}

func ColorReset() {
	print(ScolorReset())
}

func ScolorReset() string {
	return "\x1b[0m"
}
