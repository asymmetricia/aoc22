package main

import (
	"bytes"
	"os"
	"strconv"
	"strings"
	"unicode"

	"github.com/sirupsen/logrus"

	"github.com/asymmetricia/aoc22/aoc"
)

var log = logrus.StandardLogger()

// To mix the file, move each number forward or backward in the file a number of positions equal to the value of the number being moved.

type Position struct {
	Value       int
	Index       int
	Left, Right *Position
}

func (p *Position) FindValue(value int) *Position {
	if p.Value == value {
		return p
	}
	pp := p.Right
	for pp != p {
		if pp.Value == value {
			return pp
		}
		pp = pp.Right
	}
	return nil
}

func (p *Position) Find(index int) *Position {
	if p.Index == index {
		return p
	}
	pp := p.Right
	for pp != p {
		if pp.Index == index {
			return pp
		}
		pp = pp.Right
	}
	return nil
}

func (p *Position) MoveRight(n int) {
	pp := p
	for i := 0; i < n; i++ {
		pp.Value, pp.Right.Value = pp.Right.Value, pp.Value
		pp.Index, pp.Right.Index = pp.Right.Index, pp.Index
		pp = pp.Right
	}
}

func (p *Position) MoveLeft(n int) {
	pp := p
	for i := 0; i < n; i++ {
		pp.Value, pp.Left.Value = pp.Left.Value, pp.Value
		pp.Index, pp.Left.Index = pp.Left.Index, pp.Index
		pp = pp.Left
	}
}

func (p *Position) String() string {
	builder := strings.Builder{}
	builder.WriteString("[ ")
	builder.WriteString(strconv.Itoa(p.Value))
	pp := p.Right
	for pp != p {
		builder.WriteString(", ")
		builder.WriteString(strconv.Itoa(pp.Value))
		pp = pp.Right
	}
	builder.WriteString(" ]")
	return builder.String()
}

func (p *Position) Slice() []int {
	ret := []int{p.Value}
	pp := p.Right
	for pp != p {
		ret = append(ret, pp.Value)
		pp = pp.Right
	}
	return ret
}

func (p *Position) InsertAfter(index, value int) {
	// a <-> b
	// a <-> newpos <-> b
	a := p
	b := p.Right
	newpos := &Position{Value: value, Index: index}

	a.Right = newpos
	newpos.Left = a

	newpos.Right = b
	b.Left = newpos
}

func mix(name string, nums []int) *Position {
	head := &Position{Value: nums[0], Index: 0}
	head.Left = head
	head.Right = head
	tail := head
	for i, n := range nums[1:] {
		tail.InsertAfter(i+1, n)
		tail = tail.Right
	}

	for i := 0; i < 10; i++ {
		for j, toMix := range nums {
			pos := head.Find(j)
			if toMix < 0 {
				pos.MoveLeft((-toMix) % (len(nums) - 1))
			} else {
				pos.MoveRight(toMix % (len(nums) - 1))
			}
			if name == "test" {
				log.Print(toMix, " ", head)
			}
		}
	}

	return head
}

func solution(name string, input []byte) int {
	// trim trailing space only
	input = bytes.Replace(input, []byte("\r"), []byte(""), -1)
	input = bytes.TrimRightFunc(input, unicode.IsSpace)
	lines := strings.Split(strings.TrimRightFunc(string(input), unicode.IsSpace), "\n")
	log.Printf("read %d %s lines", len(lines), name)

	var nums []int
	for _, line := range lines {
		nums = append(nums, aoc.Int(line)*key)
	}

	result := mix(name, nums).FindValue(0)

	if result == nil {
		panic("could not find zero?!")
	}

	ans := 0
	for j := 0; j < 3; j++ {
		for i := 0; i < 1000; i++ {
			result = result.Right
		}
		log.Print((j+1)*1000, result.Value)
		ans += result.Value
	}

	log.Print(ans)
	if name == "test" && ans != 1623178306 {
		panic("nope")
	}
	return ans
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

	input := aoc.Input(2022, 20)
	log.Printf("input solution: %d", solution("input", input))
}

const key = 811589153
