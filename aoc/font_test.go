package aoc

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTypesetBytes(t *testing.T) {
	tests := []struct {
		name string
		line string
		opts []TypesetOpts
		want [][]byte
	}{
		{"render T", "T", nil, [][]byte{
			[]byte(" ######"),
			[]byte("    #"),
			[]byte("    #"),
			[]byte("    #"),
			[]byte("    #"),
			[]byte("    #"),
			[]byte("    #"),
			[]byte("    #"),
		}},
		{"render Tricia", "Tricia", nil, [][]byte{
			[]byte(" ######"),
			[]byte("    #"),
			[]byte("    #              #               #      ###"),
			[]byte("    #    # ##             ###            #   #"),
			[]byte("    #    ##  #           #   #             ###"),
			[]byte("    #    #   #     #     #         #      #  #"),
			[]byte("    #    #         #     #   #     #     #   #"),
			[]byte("    #    #         #      ###      #      ### #"),
		}},
		{"render T 2X", "T", []TypesetOpts{{2, false}}, [][]byte{
			[]byte("  ############"),
			[]byte("  ############"),
			[]byte("        ##"),
			[]byte("        ##"),
			[]byte("        ##"),
			[]byte("        ##"),
			[]byte("        ##"),
			[]byte("        ##"),
			[]byte("        ##"),
			[]byte("        ##"),
			[]byte("        ##"),
			[]byte("        ##"),
			[]byte("        ##"),
			[]byte("        ##"),
			[]byte("        ##"),
			[]byte("        ##"),
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TypesetBytes(tt.line, tt.opts...)
			require.Equal(t, len(tt.want), len(got), "height")
			for y := 0; y < len(got); y++ {
				require.Equalf(t, string(tt.want[y]), string(got[y]), "row %d", y)
			}
		})
	}
}
