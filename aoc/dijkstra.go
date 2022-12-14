package aoc

func Path[Cell comparable](from Cell, prev map[Cell]Cell) []Cell {
	ret := []Cell{from}
	cursor := from
	for {
		var ok bool
		cursor, ok = prev[cursor]
		if !ok {
			break
		}
		ret = append(ret, cursor)
	}
	for i := 0; i < len(ret)/2; i++ {
		ret[i], ret[len(ret)-1-i] = ret[len(ret)-1-i], ret[i]
	}

	return ret
}

func ConstantCost[K any](K, K) int {
	return 1
}

type Equaler[Cell comparable] interface {
	Equal(Cell) bool
}

// Dijkstra implements a generic Dijkstra's Algorithm, which is guaranteed to
// find the shortest path from start to end, with edges given by repeated calls
// to neighbors().
//
// length should return the length of any given edge. callbacks are optional,
// used for status reporting or visualization.
//
// returns the path including the start Cell; or nil if there is no path.
//
// If Cell type implements an `(Cell) Equal(Cell) bool` method, then this is used
// to compare reachable cells to the end state. Otherwise, simple equality is
// used.
func Dijkstra[Cell comparable](
	start Cell,
	end Cell,
	neighbors func(a Cell) []Cell,
	length func(a, b Cell) int,
	callback ...func(
		q *PQueue[Cell],
		dist map[Cell]int,
		prev map[Cell]Cell,
		current Cell)) []Cell {

	dist := map[Cell]int{}
	dist[start] = 0
	q := &PQueue[Cell]{}
	q.AddWithPriority(start, 0)

	prev := map[Cell]Cell{}

	for q.Head != nil {
		u := q.Pop()

		for _, cb := range callback {
			cb(q, dist, prev, u)
		}

		if ue, ok := any(u).(Equaler[Cell]); ok && ue.Equal(end) {
			return Path(u, prev)
		} else if u == end {
			return Path(u, prev)
		}

		for _, v := range neighbors(u) {
			alt := dist[u] + length(u, v)

			dv, ok := dist[v]
			if !ok || alt < dv {
				dist[v] = alt
				prev[v] = u
				q.AddWithPriority(v, alt)
			}
		}
	}

	// no path
	return nil
}
