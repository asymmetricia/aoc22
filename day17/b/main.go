package main

import (
	"bytes"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strings"
	"unicode"

	"github.com/asymmetricia/aoc22/aoc"
	"github.com/asymmetricia/aoc22/term"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
)

var log = logrus.StandardLogger()

/*
####

.#.
###
.#.

..#
..#
###

#
#
#
#

##
##
*/

var symbols = []symbol{
	{
		0b0011110},
	{
		0b0001000,
		0b0011100,
		0b0001000},
	{
		0b0011100,
		0b0000100,
		0b0000100},
	{
		0b0010000,
		0b0010000,
		0b0010000,
		0b0010000},
	{
		0b0011000,
		0b0011000},
}

type symbol []uint8

func (s symbol) Left() symbol {
	ret := make(symbol, len(s))
	for i := range ret {
		if s[i]&(1<<6) > 0 {
			copy(ret, s)
			return ret
		}
		ret[i] = s[i] << 1
	}
	return ret
}

func (s symbol) Right() symbol {
	ret := make(symbol, len(s))
	for i := range ret {
		if s[i]&(1<<0) > 0 {
			copy(ret, s)
			return ret
		}
		ret[i] = s[i] >> 1
	}
	return ret
}

func printColumn(col []uint8) {
	tw := term.MustWidth()
	for i := uint(0); i < tw; i++ {
		print("-")
	}
	println()
	for s := 6; s >= 0; s-- {
		for _, cell := range col {
			if cell&(1<<s) > 0 {
				print("#")
			} else {
				print(" ")
			}
		}
		println()
	}
	for i := uint(0); i < tw; i++ {
		print("-")
	}
	println()
}

type Snapshot struct {
	// Metadata
	Height      int64
	SymbolCount int64

	// Key Fields
	Column    []uint8
	SymbolNum uint8
	JetIndex  int
}

func (s Snapshot) String() string {
	return fmt.Sprintf("{@%d height=%d sym=%d jet=%d}", s.SymbolCount, s.Height, s.SymbolNum, s.JetIndex)
}

func (s Snapshot) Equals(b Snapshot) bool {
	return s.SymbolNum == b.SymbolNum &&
		s.JetIndex == b.JetIndex &&
		slices.Equal(s.Column, b.Column)
}

func solution(name string, input []byte) int64 {
	var snapshots []Snapshot

	// depth is how much puzzle is below the view point, that we've cut off.
	var depth int64

	// trim trailing space only
	input = bytes.Replace(input, []byte("\r"), []byte(""), -1)
	input = bytes.TrimRightFunc(input, unicode.IsSpace)
	lines := strings.Split(strings.TrimRightFunc(string(input), unicode.IsSpace), "\n")
	log.Printf("read %d %s lines", len(lines), name)
	log.Print(len(lines[0]))

	jetIndex := 0
	symbolCount := int64(0)
	column := []uint8{0xFF, 0}
	var end int
	const target = 1000000000000

	for symbolCount < target {
		snapshot := Snapshot{
			SymbolCount: symbolCount,
			Height:      depth + int64(end),

			SymbolNum: uint8(symbolCount % 5),
			Column:    slices.Clone(column),
			JetIndex:  jetIndex,
		}
		snapshots = append(snapshots, snapshot)
		for _, ss := range snapshots {
			delta := snapshot.SymbolCount - ss.SymbolCount
			if symbolCount+delta > target {
				continue
			}

			if ss.SymbolCount != snapshot.SymbolCount && ss.Equals(snapshot) {
				log.Printf("Prior: %s", ss)
				log.Printf("Now:   %s", snapshot)
				heightDelta := snapshot.Height - ss.Height
				count := (target - symbolCount) / delta
				symbolCount += count * delta
				depth += count * heightDelta
			}
		}

		end += 4
		for cap(column) <= end+4 || len(column) <= end+4 {
			nc := make([]uint8, cap(column)*2)
			copy(nc, column)
			column = nc
		}

		sym := symbols[symbolCount%5]
	moving:
		for {
			dir := lines[0][jetIndex]
			jetIndex++
			if jetIndex >= len(lines[0]) {
				jetIndex = 0
			}

			var proposed symbol
			if dir == '<' {
				proposed = sym.Left()
			} else {
				proposed = sym.Right()
			}
			can := true
			for i, row := range proposed {
				if end+i < len(column) && row&column[end+i] > 0 {
					can = false
					break
				}
			}
			if can {
				sym = proposed
			}

			for i, row := range sym {
				if end-1+i < len(column) && row&column[end-1+i] > 0 {
					break moving
				}
			}
			end--
		}
		for i, row := range sym {
			column[end+i] |= row
		}

		chop := 0
		// scan the puzzle for two lines that form a 7-wide barrier
		for i, row := range column {
			if i == 0 {
				continue
			}
			if row|column[i-1] == 0x7F {
				chop = i - 1
			}
		}

		// remove the first `chop` rows
		for i := 0; i < len(column); i++ {
			if i+chop >= len(column) {
				column[i] = 0
			} else {
				column[i] = column[i+chop]
			}
		}

		// record that there are `chop` rows below us
		depth += int64(chop)

		symbolCount++

		end = 0
		for {
			if end >= len(column)-1 || column[end+1] == 0 {
				break
			}
			end++
		}
	}

	return depth + int64(end)
}

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02T15:04:05",
	})
	test, err := os.ReadFile("test")
	if err == nil {
		log.Printf("test solution: %d", solution("test", test))
	} else {
		log.Warningf("no test data present: %v", err)
	}

	input := aoc.Input(2022, 17)
	log.Printf("input solution: %d", solution("input", input))
}
