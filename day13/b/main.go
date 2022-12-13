package main

import (
	"bytes"
	"encoding/json"
	"os"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/sirupsen/logrus"

	"github.com/asymmetricia/aoc22/aoc"
)

var log = logrus.StandardLogger()

type PacketValue struct {
	List bool
	int
	Packet
}

type Packet []PacketValue

func (p *Packet) String() string {
	ret := strings.Builder{}
	ret.WriteRune('[')
	for i, pv := range *p {
		if i != 0 {
			ret.WriteRune(',')
		}
		if pv.Packet != nil {
			ret.WriteString(pv.Packet.String())
		} else {
			ret.WriteString(strconv.Itoa(pv.int))
		}
	}
	ret.WriteRune(']')
	return ret.String()
}

func (p *Packet) UnmarshalJSON(data []byte) error {
	var s []json.RawMessage
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	for _, elem := range s {
		var i int
		err := json.Unmarshal(elem, &i)
		if err == nil {
			*p = append(*p, PacketValue{int: i})
			continue
		}

		var pk Packet
		if err = json.Unmarshal(elem, &pk); err != nil {
			return err
		}
		*p = append(*p, PacketValue{List: true, Packet: pk})
	}

	return nil
}

// Compare returns -1 if p less than b, 0 if equal/indeterminate, 1 otherwise
func (p Packet) Compare(b Packet) int {
	// If both values are lists, compare the first value of each list, then the
	// second value, and so on.
	for i, pv := range p {
		// If the right list runs out of items first, the inputs are not in the right
		// order.
		if i >= len(b) {
			return 1
		}
		bv := b[i]

		if !pv.List && !bv.List {
			// If both values are integers, the lower integer should come first. If the left
			// integer is lower than the right integer, the inputs are in the right order. If
			// the left integer is higher than the right integer, the inputs are not in the
			// right order. Otherwise, the inputs are the same integer; continue checking the
			// next part of the input.
			if pv.int < bv.int {
				return -1
			}
			if pv.int > bv.int {
				return 1
			}
		} else {
			// If exactly one value is an integer, convert the integer to a list which
			// contains that integer as its only value, then retry the comparison. For
			// example, if comparing [0,0,0] and 2, convert the right value to [2] (a list
			// containing 2); the result is then found by instead comparing [0,0,0] and [2].

			ppkt := pv.Packet
			if !pv.List {
				ppkt = Packet{pv}
			}
			bpkt := bv.Packet
			if !bv.List {
				bpkt = Packet{bv}
			}
			if result := ppkt.Compare(bpkt); result != 0 {
				return result
			}
		}
	}
	// If the left list runs out of items first, the inputs are in the right order.
	if len(p) < len(b) {
		return -1
	}

	// If the lists are the same length and no comparison makes a decision about the
	// order, continue checking the next part of the input.
	return 0
}

func solution(name string, input []byte) int {
	// trim trailing space only
	input = bytes.Replace(input, []byte("\r"), []byte(""), -1)
	input = bytes.TrimRightFunc(input, unicode.IsSpace)
	lines := strings.Split(strings.TrimRightFunc(string(input), unicode.IsSpace), "\n")
	log.Printf("read %d %s lines", len(lines), name)

	var packets []Packet
	for _, line := range lines {
		if line == "" {
			continue
		}

		var p Packet
		if err := json.Unmarshal([]byte(line), &p); err != nil {
			log.Fatal(err)
		}
		packets = append(packets, p)
	}

	packets = append(packets,
		Packet{{List: true, Packet: Packet{{int: 2}}}},
		Packet{{List: true, Packet: Packet{{int: 6}}}})

	sort.Slice(packets, func(i, j int) bool {
		return packets[i].Compare(packets[j]) == -1
	})

	var ret = 1
	for i, packet := range packets {
		if packet.String() == "[[2]]" || packet.String() == "[[6]]" {
			ret *= i + 1
		}
	}

	return ret
}

func main() {
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02T15:04:05",
	})
	test, err := os.ReadFile("test")
	if err == nil {
		log.Printf("test solution: %d", solution("test", test))
	} else {
		log.Warningf("no test data present")
	}

	input := aoc.Input(2022, 13)
	log.Printf("input solution: %d", solution("input", input))
}
