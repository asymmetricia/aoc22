package main

import (
	"fmt"
	"math/rand"
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
		}, {
			position{Top, coord.C(-1, 0), coord.West},
			position{West, coord.C(0, 49), coord.East},
		}, {
			position{Top, coord.C(10, 50), coord.South},
			position{South, coord.C(10, 0), coord.South},
		},

		{
			position{North, coord.C(10, -1), coord.North},
			position{West, coord.C(10, 49), coord.North},
		}, {
			position{North, coord.C(50, 10), coord.East},
			position{Bottom, coord.C(10, 49), coord.North},
		}, {
			position{North, coord.C(-1, 10), coord.West},
			position{Top, coord.C(10, 0), coord.South},
		}, {
			position{North, coord.C(0, 50), coord.South},
			position{East, coord.C(0, 0), coord.South},
		},

		{
			position{East, coord.C(10, -1), coord.North},
			position{North, coord.C(10, 49), coord.North}},
		{
			position{East, coord.C(50, 10), coord.East},
			position{Bottom, coord.C(49, 39), coord.West}},
		{
			position{East, coord.C(-1, 10), coord.West},
			position{Top, coord.C(49, 10), coord.West}},
		{
			position{East, coord.C(10, 50), coord.South},
			position{South, coord.C(49, 10), coord.West}},

		{
			position{West, coord.C(10, -1), coord.North},
			position{South, coord.C(0, 10), coord.East}},
		{
			position{West, coord.C(50, 10), coord.East},
			position{Bottom, coord.C(0, 10), coord.East}},
		{
			position{West, coord.C(-1, 10), coord.West},
			position{Top, coord.C(0, 39), coord.East}},
		{
			position{West, coord.C(10, 50), coord.South},
			position{North, coord.C(10, 0), coord.South}},

		{
			position{South, coord.C(10, -1), coord.North},
			position{Top, coord.C(10, 49), coord.North},
		}, {
			position{South, coord.C(50, 10), coord.East},
			position{East, coord.C(10, 49), coord.North},
		}, {
			position{South, coord.C(-1, 10), coord.West},
			position{West, coord.C(10, 0), coord.South},
		}, {
			position{South, coord.C(10, 50), coord.South},
			position{Bottom, coord.C(10, 0), coord.South},
		},

		{
			position{Bottom, coord.C(10, -1), coord.North},
			position{South, coord.C(10, 49), coord.North}},
		{
			position{Bottom, coord.C(50, 10), coord.East},
			position{East, coord.C(49, 39), coord.West}},
		{
			position{Bottom, coord.C(-1, 10), coord.West},
			position{West, coord.C(49, 10), coord.West}},
		{
			position{Bottom, coord.C(10, 50), coord.South},
			position{North, coord.C(49, 10), coord.West}},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s heading %s", tt.pos.side, tt.pos.facing), func(t *testing.T) {
			got := normalize(tt.pos)
			require.Equal(t, tt.want, got)
		})
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

func Test_normalize_wrap(t *testing.T) {
	for _, side := range []side{Top, Bottom, North, East, South, West} {
		for _, dir := range []coord.Direction{coord.North, coord.East, coord.South, coord.West} {
			t.Run(fmt.Sprintf("%s facing %s", side, dir), func(t *testing.T) {
				start := position{
					side:   side,
					pos:    coord.C(rand.Intn(50), rand.Intn(50)),
					facing: dir}
				pos := start
				for i := 0; i < 200; i++ {
					if i > 0 {
						require.NotEqual(t, start, pos)
					}
					pos.pos = pos.pos.Move(pos.facing)
					pos = normalize(pos)
				}
				require.Equal(t, start, pos)
			})
		}
	}
}

func Test_globalFromMap(t *testing.T) {
	tests := []struct {
		s    side
		c    coord.Coord
		want coord.Coord
	}{
		{Top, coord.C(5, 6), coord.C(55, 6)},
		{East, coord.C(5, 6), coord.C(105, 6)},
		{South, coord.C(5, 6), coord.C(55, 56)},
		{West, coord.C(5, 6), coord.C(5, 106)},
		{Bottom, coord.C(5, 6), coord.C(55, 106)},
		{North, coord.C(5, 6), coord.C(5, 156)},
	}
	for _, tt := range tests {
		require.Equal(t, tt.want, globalFromMap(tt.s, tt.c))
	}
}
