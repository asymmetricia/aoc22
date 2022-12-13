package main

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPacket_Compare(t *testing.T) {
	tests := []struct {
		a    string
		b    string
		want int
	}{
		{"[1,1,3,1,1]", "[1,1,5,1,1]", -1},
		{"[2,3,4]", "[4]", -1},
		{"[[1],[2,3,4]]", "[[1],4]", -1},
		{"[9]", "[[8,7,6]]", 1},
		{"[[4,4],4,4]", "[[4,4],4,4,4]", -1},
		{"[[[]]]", "[[]]", 1},
	}
	for _, tt := range tests {
		var a, b Packet
		json.Unmarshal([]byte(tt.a), &a)
		json.Unmarshal([]byte(tt.b), &b)
		require.Equal(t, tt.want, a.Compare(b))
	}
}
