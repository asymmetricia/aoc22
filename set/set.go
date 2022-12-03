package set

import (
	"sort"
	"strings"

	"golang.org/x/exp/constraints"
)

type Set[K comparable] map[K]bool

func (a Set[K]) Intersect(b Set[K]) Set[K] {
	ret := Set[K]{}
	for k := range a {
		if b[k] {
			ret[k] = true
		}
	}
	return ret
}

func (a Set[K]) Copy() Set[K] {
	ret := Set[K]{}
	for k := range a {
		ret[k] = true
	}
	return ret
}

func (a Set[K]) Union(b Set[K]) Set[K] {
	ret := Set[K]{}
	for k := range a {
		ret[k] = true
	}
	for k := range b {
		ret[k] = true
	}
	return ret
}

func Items[K constraints.Ordered](a Set[K]) []K {
	var ret []K
	for k := range a {
		ret = append(ret, k)
	}

	// have to, otherwise map traversal is unstable
	sort.Slice(ret, func(i, j int) bool {
		return ret[i] < ret[j]
	})
	return ret
}

func (a Set[K]) Difference(b Set[K]) Set[K] {
	ret := Set[K]{}
	for k := range a {
		if !b[k] {
			ret[k] = true
		}
	}
	return ret
}

// Item returns a "random" item from the set, or the zero value if there are
// no items in the set.
func (a Set[K]) Item() K {
	for s := range a {
		return s
	}
	var zero K
	return zero
}

func FromItems[K comparable](items []K) Set[K] {
	s := Set[K]{}
	for _, i := range items {
		s[i] = true
	}
	return s
}

func FromString(s string) Set[rune] {
	return FromItems([]rune(s))
}

func FromWords(s string) Set[string] {
	return FromItems(strings.Fields(s))
}
