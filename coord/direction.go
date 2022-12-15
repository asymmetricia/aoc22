package coord

type Direction int

var Directions = []Direction{
	North, NorthEast, East, SouthEast, South, SouthWest, West, NorthWest,
}

const (
	North Direction = iota
	NorthEast
	East
	SouthEast
	South
	SouthWest
	West
	NorthWest
)

var DirectionStrings = map[string]Direction{
	"n":  North,
	"ne": NorthEast,
	"e":  East,
	"se": SouthEast,
	"s":  South,
	"sw": SouthWest,
	"w":  West,
	"nw": NorthWest,
}
