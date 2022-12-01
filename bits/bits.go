package bits

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type Packet struct {
	Version int64
	Type    int64
	Packets []Packet
	Value   int64
}

func (p Packet) String() string {
	ret := ""
	outro := ")"

	if len(p.Packets) == 1 {
		return p.Packets[0].String()
	}

	switch p.Type {
	case TypeSum:
		ret = "(+ "
	case TypeProduct:
		ret = "(* "
	case TypeMin:
		ret = "(min "
	case TypeMax:
		ret = "(max "
	case TypeLiteral:
		return strconv.FormatInt(p.Value, 10)
	case TypeGT:
		ret = "(if (> "
		outro = ") 1 0)"
	case TypeLT:
		ret = "(if (< "
		outro = ") 1 0)"
	case TypeEq:
		ret = "(if (= "
		outro = ") 1 0)"
	}

	indent := ""
	for i := 0; i < len(ret); i++ {
		indent += " "
	}
	var ps []string
	for _, p := range p.Packets {
		ps = append(ps, strings.Join(strings.Split(p.String(), "\n"), "\n"+indent))
	}
	ret += strings.Join(ps, "\n"+indent)
	return ret + outro
}

var decode = map[rune]string{
	'0': "0000",
	'1': "0001",
	'2': "0010",
	'3': "0011",
	'4': "0100",
	'5': "0101",
	'6': "0110",
	'7': "0111",
	'8': "1000",
	'9': "1001",
	'A': "1010",
	'B': "1011",
	'C': "1100",
	'D': "1101",
	'E': "1110",
	'F': "1111",
}

func Decode(packet string) (string, error) {
	var ret string
	for _, c := range packet {
		ex, ok := decode[c]
		if !ok {
			return "", fmt.Errorf("could not expand %c", c)
		}
		ret += ex
	}
	return ret, nil
}

func Parse(packet string) (Packet, int64, error) {
	p := Packet{}
	var err error

	p.Version, err = strconv.ParseInt("0b"+packet[0:3], 0, 64)
	if err != nil {
		return Packet{}, 0, fmt.Errorf(
			"could not parse binary version number %q: %w",
			packet[0:3],
			err)
	}

	p.Type, err = strconv.ParseInt("0b"+packet[3:6], 0, 64)
	if err != nil {
		return Packet{}, 0, fmt.Errorf(
			"could not parse binary type number %q: %w",
			packet[3:6],
			err)
	}

	if p.Type == TypeLiteral {
		number := "0b"
		cursor := int64(6)
		for {
			number += packet[cursor+1 : cursor+5]
			if packet[cursor] == '1' {
				cursor += 5
				continue
			}

			p.Value, err = strconv.ParseInt(number, 0, 64)
			if err != nil {
				return Packet{}, 0, fmt.Errorf(
					"could not parse literal number %q: %w",
					number, err)
			}
			return p, cursor + 5, nil
		}
	}

	if packet[6] == '0' {
		// length is in bits
		packetBytes, err := strconv.ParseInt(packet[7:22], 2, 64)
		if err != nil {
			return Packet{}, 0, fmt.Errorf(
				"could not parse 15-bit length %q: %w",
				packet[7:22], err)
		}
		cursor := int64(22)
		for cursor < 22+packetBytes {
			subp, consumed, err := Parse(packet[cursor : 22+packetBytes])
			if err != nil {
				return Packet{}, 0, err
			}
			cursor += consumed
			p.Packets = append(p.Packets, subp)
		}
		return p, cursor, nil
	}

	// length is in # of packets
	packets, err := strconv.ParseInt(packet[7:18], 2, 64)
	if err != nil {
		return Packet{}, 0, fmt.Errorf(
			"could not parse 11-bit length %q: %w",
			packet[7:18], err)
	}
	cursor := int64(18)
	for len(p.Packets) < int(packets) {
		subp, consumed, err := Parse(packet[cursor:])
		if err != nil {
			return Packet{}, 0, err
		}
		cursor += consumed
		p.Packets = append(p.Packets, subp)
	}
	return p, cursor, nil
}

const (
	TypeSum     int64 = 0
	TypeProduct       = 1
	TypeMin           = 2
	TypeMax           = 3
	TypeLiteral       = 4
	TypeGT            = 5
	TypeLT            = 6
	TypeEq            = 7
)

func Execute(p Packet) (int64, error) {
	var op func(inputs []int64) int64

	switch p.Type {
	case TypeSum:
		op = func(inputs []int64) int64 {
			sum := int64(0)
			for _, v := range inputs {
				sum += v
			}
			return sum
		}
	case TypeProduct:
		op = func(inputs []int64) int64 {
			product := int64(1)
			for _, v := range inputs {
				product *= v
			}
			return product
		}
	case TypeMin:
		op = func(inputs []int64) int64 {
			min := int64(math.MaxInt64)
			for _, v := range inputs {
				if v < min {
					min = v
				}
			}
			return min
		}
	case TypeMax:
		op = func(inputs []int64) int64 {
			max := int64(math.MinInt64)
			for _, v := range inputs {
				if v > max {
					max = v
				}
			}
			return max
		}
	case TypeLiteral:
		return p.Value, nil
	case TypeGT:
		op = func(inputs []int64) int64 {
			if inputs[0] > inputs[1] {
				return 1
			}
			return 0
		}
	case TypeLT:
		op = func(inputs []int64) int64 {
			if inputs[0] < inputs[1] {
				return 1
			}
			return 0
		}
	case TypeEq:
		op = func(inputs []int64) int64 {
			if inputs[0] == inputs[1] {
				return 1
			}
			return 0
		}
	default:
		return 0, fmt.Errorf("unknown packet type %d", p.Type)
	}

	inputs := make([]int64, 0, len(p.Packets))
	for i, p := range p.Packets {
		pv, err := Execute(p)
		if err != nil {
			return 0, fmt.Errorf("packet #%d: %w", i, err)
		}
		inputs = append(inputs, pv)
	}

	return op(inputs), nil
}
