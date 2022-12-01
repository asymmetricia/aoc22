package set

import "sort"

type Set map[string]bool

func (a Set) Intersect(b Set) Set {
	ret := Set{}
	for k := range a {
		if b[k] {
			ret[k] = true
		}
	}
	return ret
}

func (a Set) Copy() Set {
	ret := Set{}
	for k := range a {
		ret[k] = true
	}
	return ret
}

func (a Set) Union(b Set) Set {
	ret := Set{}
	for k := range a {
		ret[k] = true
	}
	for k := range b {
		ret[k] = true
	}
	return ret
}

func (a Set) Items() []string {
	var ret []string
	for k := range a {
		ret = append(ret, k)
	}

	// have to, otherwise map traversal is unstable
	sort.Strings(ret)
	return ret
}

func (a Set) Difference(b Set) Set {
	ret := Set{}
	for k := range a {
		if !b[k] {
			ret[k] = true
		}
	}
	return ret
}

// Item returns a "random" item from the set, or the blank string if there are
// no items in the set.
func (a Set) Item() string {
	for s := range a {
		return s
	}
	return ""
}
