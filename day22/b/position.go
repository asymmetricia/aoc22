package main

import (
	"fmt"

	"github.com/asymmetricia/aoc22/coord"
)

type position struct {
	side   side
	pos    coord.Coord
	facing coord.Direction
}

func (p position) String() string {
	return fmt.Sprintf("at %s on side %s facing %s", p.pos.String(), p.side.String(), p.facing.String())
}
