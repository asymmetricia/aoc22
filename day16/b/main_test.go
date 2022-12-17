package main

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func Test_permutations(t *testing.T) {
	tests := []struct {
		in   []string
		want [][]string
	}{
		{[]string{"a", "b"}, [][]string{{"a", "b"}, {"b", "a"}}},
		{[]string{"a", "b", "c"}, [][]string{
			{"a", "b", "c"}, {"a", "c", "b"}, {"c", "a", "b"},
			{"c", "b", "a"}, {"b", "c", "a"}, {"b", "a", "c"}}},
		{[]string{"a", "b", "c", "d"}, [][]string{
			{"a", "b", "c", "d"}, {"a", "b", "d", "c"}, {"a", "d", "b", "c"}, {"d", "a", "b", "c"},
			{"d", "a", "c", "b"},
			{"a", "d", "c", "b"},
			{"a", "c", "d", "b"},
			{"a", "c", "b", "d"},
			{"c", "a", "b", "d"},
			{"c", "a", "d", "b"},
			{"c", "d", "a", "b"},
			{"d", "c", "a", "b"},
			{"d", "c", "b", "a"},
			{"c", "d", "b", "a"},
			{"c", "b", "d", "a"},
			{"c", "b", "a", "d"},
			{"b", "c", "a", "d"},
			{"b", "c", "d", "a"},
			{"b", "d", "c", "a"},
			{"d", "b", "c", "a"},
			{"d", "b", "a", "c"},
			{"b", "d", "a", "c"},
			{"b", "a", "d", "c"},
			{"b", "a", "c", "d"},
		}},
	}
	for _, tt := range tests {
		var got [][]string
		for s := range permutations(tt.in, len(tt.in), len(tt.in), logrus.StandardLogger()) {
			got = append(got, s)
		}
		require.ElementsMatch(t, tt.want, got)
	}
}
