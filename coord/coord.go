package coord

import (
	"fmt"
	"strconv"
	"strings"
)

type Coord struct {
	X, Y int
}

func (c Coord) String() string {
	return fmt.Sprintf("(%d,%d)", c.X, c.Y)
}

func C(x, y int) Coord {
	return Coord{x, y}
}

func (c Coord) Neighbors(diag bool) []Coord {
	if diag {
		return []Coord{
			c.North(), c.NorthEast(),
			c.East(), c.SouthEast(),
			c.South(), c.SouthWest(),
			c.West(), c.NorthWest(),
		}
	}
	return []Coord{
		c.North(), c.East(), c.South(), c.West(),
	}
}

func MustFromComma(xy string) Coord {
	c, e := FromComma(xy)
	if e != nil {
		panic(e)
	}
	return c
}

func FromComma(xy string) (Coord, error) {
	parts := strings.Split(strings.TrimSpace(xy), ",")
	if len(parts) != 2 {
		return Coord{}, fmt.Errorf("expected two ,-separated parts, got %d", len(parts))
	}
	var ret Coord
	var err error
	ret.X, err = strconv.Atoi(parts[0])
	if err != nil {
		return Coord{}, fmt.Errorf("bad X coordinate %q: %w", parts[0], err)
	}
	ret.Y, err = strconv.Atoi(parts[1])
	if err != nil {
		return Coord{}, fmt.Errorf("bad Y coordinate %q: %w", parts[1], err)
	}
	return ret, nil
}

func (c Coord) Move(d Direction) Coord {
	switch d {
	case North:
		return c.North()
	case NorthEast:
		return c.NorthEast()
	case East:
		return c.East()
	case SouthEast:
		return c.SouthEast()
	case South:
		return c.South()
	case SouthWest:
		return c.SouthWest()
	case West:
		return c.West()
	case NorthWest:
		return c.NorthWest()
	}
	panic("bad direction " + strconv.Itoa(int(d)))
}

func (c Coord) North() Coord {
	return Coord{c.X, c.Y - 1}
}
func (c Coord) South() Coord {
	return Coord{c.X, c.Y + 1}
}
func (c Coord) East() Coord {
	return Coord{c.X + 1, c.Y}
}
func (c Coord) West() Coord {
	return Coord{c.X - 1, c.Y}
}
func (c Coord) NorthEast() Coord {
	return Coord{c.X + 1, c.Y + 1}
}
func (c Coord) SouthEast() Coord {
	return Coord{c.X + 1, c.Y - 1}
}
func (c Coord) NorthWest() Coord {
	return Coord{c.X - 1, c.Y + 1}
}
func (c Coord) SouthWest() Coord {
	return Coord{c.X - 1, c.Y - 1}
}

func (c Coord) Execute(steps []string) Coord {
	for _, step := range steps {
		if dir, ok := DirectionStrings[step]; ok {
			c = c.Move(dir)
		} else {
			panic(step)
		}
	}
	return c
}

func (c Coord) Plus(a Coord) Coord {
	return Coord{c.X + a.X, c.Y + a.Y}
}
func (c Coord) Equal(a Coord) bool {
	return c.X == a.X && c.Y == a.Y
}
