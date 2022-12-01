package bits

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestString(t *testing.T) {
	tests := []struct {
		name   string
		packet string
		want   string
	}{
		{"basic", "D2FE28", "2021"},
		{"ex1", "38006F45291200", "(if (< 10 20) 1 0)"},
		{"ex2", "EE00D40C823060", "(max 1 2 3)"},
		{"ex3", "9C0141080250320F1802104A08", "(if (= (+ 1 3) (* 2 2)) 1 0)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, e := Decode(tt.packet)
			require.NoError(t, e)

			packet, _, e := Parse(d)
			require.NoError(t, e)

			require.Equal(t, tt.want, packet.String())
		})
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		packet  string
		want    Packet
		wantLen int64
	}{
		{"literal", "D2FE28", Packet{
			Version: 6,
			Type:    4,
			Packets: nil,
			Value:   2021,
		}, 21},
		{"length type 0 opeerator", "38006F45291200", Packet{
			Version: 1,
			Type:    6,
			Packets: []Packet{
				{6, 4, nil, 10},
				{2, 4, nil, 20}},
			Value: 0,
		}, 49},
		{"length type 1 operator", "EE00D40C823060", Packet{
			Version: 7,
			Type:    3,
			Packets: []Packet{
				{2, 4, nil, 1},
				{4, 4, nil, 2},
				{1, 4, nil, 3}},
			Value: 0,
		}, 51},
		{"nested operators", "8A004A801A8002F478", Packet{
			Version: 4,
			Type:    2,
			Packets: []Packet{
				{1, 2, []Packet{
					{5, 2, []Packet{
						{6, 4, nil, 15},
					}, 0},
				}, 0}},
			Value: 0,
		}, 69},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			decoded, err := Decode(tt.packet)
			require.NoError(t, err)

			got, gotLen, err := Parse(decoded)
			require.NoError(t, err)
			require.NotNil(t, got)
			require.Equal(t, tt.want, got)
			require.Equal(t, tt.wantLen, gotLen)
		})
	}
}
