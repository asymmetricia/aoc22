package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"time"
	"unicode"

	"github.com/sirupsen/logrus"
	"golang.org/x/exp/maps"

	"github.com/asymmetricia/aoc22/aoc"
)

var log = logrus.StandardLogger()

type Op int

const (
	Constant Op = iota
	Add
	Subtract
	Multiply
	Divide
)

type Monkey struct {
	Name        string
	Value       int
	Left, Right *Monkey
	Op          Op
}

func (m *Monkey) Compute() int {
	if m.Name == "root" {
		if m.Left.Compute() == m.Right.Compute() {
			return 1
		} else {
			return 0
		}
	}

	switch m.Op {
	case Constant:
		return m.Value
	case Add:
		return m.Left.Compute() + m.Right.Compute()
	case Subtract:
		return m.Left.Compute() - m.Right.Compute()
	case Multiply:
		return m.Left.Compute() * m.Right.Compute()
	case Divide:
		return m.Left.Compute() / m.Right.Compute()
	default:
		panic(m.Op)
	}
}

func (m *Monkey) Reverse(i int) int {
	if m.Name == "humn" {
		return i
	}

	left := m.Left
	right := m.Right

	if m.Op != Constant && left.Op != Constant && right.Op != Constant {
		panic("two non-constant branches")
	}

	reverseLeft := left.Name == "humn" || right.Op == Constant

	switch m.Op {
	case Constant:
		if m.Value != i {
			panic(fmt.Sprintf("at %s cannot reverse constant %d to get value %d", m.Name, m.Value, i))
		}
	case Add:
		if reverseLeft {
			return left.Reverse(i - right.Value)
		}
		return right.Reverse(i - left.Value)
	case Subtract:
		// i = left - right
		// right = left - i
		// left = i + right

		if reverseLeft {
			return left.Reverse(i + right.Value)
		}
		return right.Reverse(left.Value - i)
	case Multiply:
		// i = left * right
		// left = i / right
		// right = i / left
		if reverseLeft {
			return left.Reverse(i / right.Value)
		}
		return right.Reverse(i / left.Value)
	case Divide:
		// i = left / right
		// left = i * right
		// right = left / i
		if reverseLeft {
			return left.Reverse(i * right.Value)
		}
		return right.Reverse(left.Value / i)
	default:
		panic(m.Op)
	}
	return -1
}

func (m *Monkey) String() string {
	switch m.Op {
	case Constant:
		return fmt.Sprintf("%s: %d", m.Name, m.Value)
	case Add:
		return fmt.Sprintf("%s: %s + %s", m.Name, m.Left.Name, m.Right.Name)
	case Subtract:
		return fmt.Sprintf("%s: %s - %s", m.Name, m.Left.Name, m.Right.Name)
	case Multiply:
		return fmt.Sprintf("%s: %s * %s", m.Name, m.Left.Name, m.Right.Name)
	case Divide:
		return fmt.Sprintf("%s: %s / %s", m.Name, m.Left.Name, m.Right.Name)
	default:
		panic(m.Op)
	}
}

func (m *Monkey) Print(i int) {
	for j := 0; j < i; j++ {
		print(" ")
	}
	println(m.String())
	if m.Op != Constant {
		m.Left.Print(i + 2)
		m.Right.Print(i + 2)
	}
}

func ParseMonkey(values map[string]*Monkey, line string) bool {
	var monkey, op1, op2 string
	var op rune
	var value int

	monkey = line[0:4]
	if _, ok := values[monkey]; ok {
		return true
	}

	_, opErr := fmt.Sscanf(line, "%4s: %4s %c %4s", &monkey, &op1, &op, &op2)
	if opErr == nil {
		if _, ok := values[monkey]; ok {
			return true
		}

		left, ok := values[op1]
		if !ok {
			return false
		}

		right, ok := values[op2]
		if !ok {
			return false
		}

		values[monkey] = &Monkey{
			Name:  monkey,
			Left:  left,
			Right: right,
			Op: map[rune]Op{
				'+': Add,
				'-': Subtract,
				'*': Multiply,
				'/': Divide,
			}[op],
		}

		return true
	}

	_, constErr := fmt.Sscanf(line, "%4s: %d", &monkey, &value)
	if constErr == nil {
		if _, ok := values[monkey]; ok {
			return true
		}

		values[monkey] = &Monkey{
			Name:  monkey,
			Value: value,
			Op:    Constant,
		}
		return true
	}

	panic(line)
}

func solution(name string, input []byte) int {
	// trim trailing space only
	input = bytes.Replace(input, []byte("\r"), []byte(""), -1)
	input = bytes.TrimRightFunc(input, unicode.IsSpace)
	lines := strings.Split(strings.TrimRightFunc(string(input), unicode.IsSpace), "\n")
	uniq := map[string]bool{}
	for _, line := range lines {
		uniq[line] = true
	}
	log.Printf("read %d %s lines (%d unique)", len(lines), name, len(uniq))

	values := map[string]*Monkey{}
	for len(values) != len(lines) {
		for _, line := range lines {
			ParseMonkey(values, line)
		}
	}

	changed := true
	for changed {
		changed = false
		for _, name := range maps.Keys(values) {
			if name == "humn" || values[name].Op == Constant {
				continue
			}
			l, r := values[name].Left, values[name].Right
			if l.Name == "humn" || r.Name == "humn" {
				continue
			}
			if values[name].Op != Constant && values[name].Left.Op == Constant && values[name].Right.Op == Constant {
				changed = true
				values[name].Value = values[name].Compute()
				values[name].Op = Constant
			}
		}
	}

	values["root"].Print(0)

	left, right := values["root"].Left, values["root"].Right
	if left.Op == Constant {
		return right.Reverse(left.Value)
	} else {
		return left.Reverse(right.Value)
	}
}

func main() {
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02T15:04:05",
	})
	test, err := os.ReadFile("test")
	soln := solution("test", test)
	if soln != 301 {
		panic("nope")
	}
	if err == nil {
		log.Printf("test solution: %d", soln)
	} else {
		log.Warningf("no test data present")
	}

	start := time.Now()
	input := aoc.Input(2022, 21)
	log.Printf("input solution: %d in %dms", solution("input", input), time.Since(start).Microseconds())
}
