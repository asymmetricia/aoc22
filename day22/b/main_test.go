package main

import (
	"fmt"
	"testing"

	"github.com/asymmetricia/aoc22/coord"
	"github.com/stretchr/testify/require"
)

func Test_normalize(t *testing.T) {
	tests := []struct {
		pos  position
		want position
	}{
		{
			position{
				side:   Top,
				pos:    coord.C(-1, 0),
				facing: coord.West,
			},
			position{
				side:   West,
				pos:    coord.C(0, 49),
				facing: coord.East,
			},
		}, {
			position{Top, coord.C(10, -1), coord.North},
			position{North, coord.C(0, 10), coord.East},
		}, {
			position{Top, coord.C(50, 10), coord.East},
			position{East, coord.C(0, 10), coord.East},
		}, {position{Top, coord.C(-1, 0), coord.West},
			position{West, coord.C(0, 49), coord.East},
		}, {
			position{Top, coord.C(10, 50), coord.South},
			position{South, coord.C(10, 0), coord.South},
		}, {
			position{South, coord.C(50, 10), coord.East},
			position{East, coord.C(10, 49), coord.North},
		}, {
			position{South, coord.C(10, 50), coord.South},
			position{Bottom, coord.C(10, 0), coord.South},
		}, {
			position{South, coord.C(-1, 10), coord.West},
			position{West, coord.C(10, 0), coord.South},
		}, {
			position{North, coord.C(0, 50), coord.South},
			position{East, coord.C(0, 0), coord.South},
		}, {
			position{North, coord.C(50, 10), coord.East},
			position{Bottom, coord.C(10, 49), coord.North},
		}, {
			position{Bottom, coord.C(50, 10), coord.East},
			position{East, coord.C(49, 39), coord.West},
		}, {
			position{West, coord.C(-1, 10), coord.West},
			position{Top, coord.C(0, 39), coord.East},
		},
	}

	for _, tt := range tests {
		got := normalize(tt.pos)
		require.Equal(t, tt.want, got)
	}
}
func Test_normalize_identity(t *testing.T) {
	for _, side := range []side{Top, Bottom, North, East, South, West} {
		tests := []struct {
			c   coord.Coord
			f   coord.Direction
			exp coord.Coord
		}{
			{coord.C(-1, 0), coord.West, coord.C(0, 0)},
			{coord.C(50, 0), coord.East, coord.C(49, 0)},
			{coord.C(10, -1), coord.North, coord.C(10, 0)},
			{coord.C(10, 50), coord.South, coord.C(10, 49)},
		}
		for _, tt := range tests {
			t.Run(fmt.Sprintf("side %s @ %s", side, tt.c), func(t *testing.T) {
				from := position{side, tt.c, tt.f}
				step := normalize(from)
				back := step
				back.facing = back.facing.CW(false).CW(false)
				back.pos = back.pos.Move(back.facing)
				to := normalize(back)
				require.Equalf(t, from.side, to.side, "{%v} to {%v} then {%v}", from, step, back)
				require.Equal(t, tt.exp, to.pos, "pos")
				require.Equal(t, tt.f.CW(false).CW(false), to.facing, "facing")
			})
		}
	}
}
