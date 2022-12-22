package coord

type Direction int

func (d Direction) CW(fortyFive bool) Direction {
	if fortyFive {
		return map[Direction]Direction{
			North:     NorthEast,
			NorthEast: East,
			East:      SouthEast,
			SouthEast: South,
			South:     SouthWest,
			SouthWest: West,
			West:      NorthWest,
			NorthWest: North,
		}[d]
	}
	return map[Direction]Direction{
		North: East,
		East:  South,
		South: West,
		West:  North,
	}[d]
}

func (d Direction) CCW(fortyFive bool) Direction {
	if fortyFive {
		return map[Direction]Direction{
			North:     NorthWest,
			NorthWest: West,
			West:      SouthWest,
			SouthWest: South,
			South:     SouthEast,
			SouthEast: East,
			East:      NorthEast,
			NorthEast: North,
		}[d]
	}
	return map[Direction]Direction{
		North: West,
		West:  South,
		South: East,
		East:  North,
	}[d]
}

func (d Direction) String() string {
	for k, v := range DirectionStrings {
		if v == d {
			return k
		}
	}
	return "(bad direction)"
}

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
