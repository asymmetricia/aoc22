package aoc

import "golang.org/x/exp/constraints"

func MaxFn[T any, K constraints.Ordered](a []T, fn func(T) K) (ret K) {
	if len(a) == 0 {
		return ret
	}

	best := fn(a[0])
	for _, item := range a[1:] {
		i := fn(item)
		if best < i {
			best = i
		}
	}
	return best
}

func Max[K constraints.Ordered](a ...K) K {
	first := true
	var max K
	for _, a := range a {
		if first || a > max {
			max = a
			first = false
		}
	}
	return max
}

func Min[K constraints.Ordered](a ...K) K {
	first := true
	var min K
	for _, a := range a {
		if first || a < min {
			min = a
			first = false
		}
	}
	return min
}

func Abs[K constraints.Signed | constraints.Float](a K) K {
	if a < 0 {
		return -a
	}
	return a
}
