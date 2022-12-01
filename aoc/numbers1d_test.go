package aoc

import (
	"reflect"
	"testing"
)

func TestUnique(t *testing.T) {
	tests := []struct {
		name string
		in   []int
		want []int
	}{
		{"empty", nil, nil},
		{"one", []int{1}, []int{1}},
		{"two unique", []int{1, 2}, []int{1, 2}},
		{"two dupe", []int{1, 1}, []int{1}},
		{"three", []int{1, 2, 3}, []int{1, 2, 3}},
		{"three dupe", []int{1, 2, 2}, []int{1, 2}},
		{"four dupe", []int{1, 2, 2, 3}, []int{1, 2, 3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Unique(tt.in); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Unique() = %v, want %v", got, tt.want)
			}
		})
	}
}
