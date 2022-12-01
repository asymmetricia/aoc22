package aoc

import (
	"math"

	"github.com/asymmetricia/aoc21/coord"
)

// AStarGraph finds the path from start to end along the grpah defined by edges
// returns from calling neighbors against each cell such that the path minimizes
// the total cost.
//
// If any callbacks are defined, they're called just before each time a cell is
// picked from the open set.
func AStarGraph[Cell comparable](
	start Cell,
	end Cell,
	neighbors func(a Cell) []Cell,
	cost func(a, b Cell) int,
	heuristic func(a Cell) int,
	callback ...func(
		openSet map[Cell]bool,
		cameFrom map[Cell]Cell,
		gScore map[Cell]int,
		fScore map[Cell]int,
		current Cell),
) []Cell {
	openSet := map[Cell]bool{start: true}
	cameFrom := map[Cell]Cell{}
	gScore := map[Cell]int{}
	fScore := map[Cell]int{
		start: 0,
	}

	for len(openSet) > 0 {
		var current Cell
		var curFscore = math.MaxInt
		first := true

		for c := range openSet {
			fs, ok := fScore[c]
			if !ok {
				fs = math.MaxInt
			}

			if first || fs < curFscore {
				first = false
				current = c
				curFscore = fs
			}
		}

		for _, cb := range callback {
			cb(openSet, cameFrom, gScore, fScore, current)
		}

		if current == end {
			break
		}

		delete(openSet, current)

		for _, neighbor := range neighbors(current) {
			curGS, ok := gScore[current]
			if !ok {
				curGS = math.MaxInt
			}

			neighGS, ok := gScore[neighbor]
			if !ok {
				neighGS = math.MaxInt32
			}

			tentativeGScore := curGS + cost(current, neighbor)
			if tentativeGScore < neighGS {
				cameFrom[neighbor] = current
				gScore[neighbor] = tentativeGScore
				fScore[neighbor] = tentativeGScore + heuristic(neighbor)
				openSet[neighbor] = true
			}
		}
	}

	ret := []Cell{end}
	cursor := end
	for {
		if cursor == start {
			break
		}
		cursor = cameFrom[cursor]
		ret = append(ret, cursor)
	}
	for i := 0; i < len(ret)/2; i++ {
		ret[i], ret[len(ret)-1-i] = ret[len(ret)-1-i], ret[i]
	}

	return ret
}

func AStarGrid[Cell any](
	grid map[coord.Coord]Cell,
	start coord.Coord,
	end coord.Coord,
	cost func(from, to coord.Coord) int,
	heuristic func(from coord.Coord) int,
	diag bool,
	callback ...func(
		openSet map[coord.Coord]bool,
		cameFrom map[coord.Coord]coord.Coord,
		gScore map[coord.Coord]int,
		fScore map[coord.Coord]int,
		current coord.Coord,
	)) []coord.Coord {
	return AStarGraph(start, end,
		func(a coord.Coord) []coord.Coord {
			var ret []coord.Coord
			for _, n := range a.Neighbors(diag) {
				if _, ok := grid[n]; ok {
					ret = append(ret, n)
				}
			}
			return ret
		},
		cost,
		heuristic,
		callback...)
}
