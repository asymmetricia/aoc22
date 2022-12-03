package canvas

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTextBox_On(t *testing.T) {
	tests := []struct {
		name   string
		tb     TextBox
		expect string
	}{
		{
			"titled",
			TextBox{
				Title: []rune("test"),
				Body:  []rune("this is a text box\nwith a title"),
			},
			"\x1b[38;2;187;187;187m┏test━━━━━━━━━━━━━━┓\n┃this is a text box┃\n┃with a title      ┃\n┗━━━━━━━━━━━━━━━━━━┛\n",
		},
		{
			"footered",
			TextBox{
				Title:  []rune("test"),
				Body:   []rune("this is a text box\nwith a title"),
				Footer: []rune("and a footer"),
			},
			"\x1b[38;2;187;187;187m┏test━━━━━━━━━━━━━━┓\n┃this is a text box┃\n┃with a title      ┃\n┗━━━━━━and a footer┛\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fb := &Canvas{}
			tt.tb.On(fb)
			got := fb.String()

			assert.Equal(t, tt.expect, got)
			if got != tt.expect {
				t.Error("expected:")
				for _, line := range strings.Split(tt.expect, "\n") {
					t.Error(line)
				}

				t.Error("got:")
				for _, line := range strings.Split(got, "\n") {
					t.Error(line)
				}
				t.Fail()
			}
		})
	}
}
