package main

import (
	"bytes"
	"os"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/sirupsen/logrus"
)

var log = logrus.StandardLogger()

type Swarm struct {
	Monkeys []*Monkey
}

func (s Swarm) Print() {
	for i, monkey := range s.Monkeys {
		log.Printf("Monkey %d inspected items %d times.", i, monkey.Inspects)
	}
}

func (s Swarm) Round() {
	for _, monkey := range s.Monkeys {
		for _, item := range monkey.Items {
			monkey.Inspects++

			// inspect item
			if monkey.WorryOp == "*" {
				if monkey.WorryOpOld {
					item *= item
				} else {
					item *= monkey.WorryOpInt
				}
			} else {
				if monkey.WorryOpOld {
					item += item
				} else {
					item += monkey.WorryOpInt
				}
			}

			// express relief
			item /= 3

			// throw item
			if item%monkey.Mod == 0 {
				s.Monkeys[monkey.IfTrue].Items = append(s.Monkeys[monkey.IfTrue].Items, item)
			} else {
				s.Monkeys[monkey.IfFalse].Items = append(s.Monkeys[monkey.IfFalse].Items, item)
			}
		}
		monkey.Items = monkey.Items[0:0]
	}
}

type Monkey struct {
	Items      []int
	WorryOp    string
	WorryOpInt int
	WorryOpOld bool
	Mod        int
	IfTrue     int
	IfFalse    int
	Inspects   int
}

func solution(name string, input []byte) int {
	// trim trailing space only
	input = bytes.Replace(input, []byte("\r"), []byte(""), -1)
	input = bytes.TrimRightFunc(input, unicode.IsSpace)
	lines := strings.Split(strings.TrimRightFunc(string(input), unicode.IsSpace), "\n")
	log.Printf("read %d %s lines", len(lines), name)

	var swarm Swarm
	//concern := map[int]int{}
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		monkey := &Monkey{}
		swarm.Monkeys = append(swarm.Monkeys, monkey)
		items := strings.Split(strings.Replace(lines[i+1][18:], " ", "", -1), ",")
		for _, item := range items {
			if itemI, err := strconv.Atoi(item); err != nil {
				log.Fatalf("%q: %v", item, err)
			} else {
				monkey.Items = append(monkey.Items, itemI)
			}
		}
		if strings.Contains(lines[i+2], "*") {
			monkey.WorryOp = "*"
		} else {
			monkey.WorryOp = "+"
		}

		opline := strings.Fields(lines[i+2])
		if len(opline) < 6 {
			log.Fatal(lines[i+2])
		}
		op := opline[5]
		if op == "old" {
			monkey.WorryOpOld = true
		} else {
			var err error
			monkey.WorryOpInt, err = strconv.Atoi(op)
			if err != nil {
				log.Fatal(line, err)
			}
		}

		testline := strings.Fields(lines[i+3])
		var err error
		monkey.Mod, err = strconv.Atoi(testline[3])
		if err != nil {
			log.Fatalf("line %q: %q", line, err)
		}

		trueline := strings.Fields(lines[i+4])
		monkey.IfTrue, err = strconv.Atoi(trueline[5])
		if err != nil {
			log.Fatalf("line %q: %q", line, err)
		}

		falseline := strings.Fields(lines[i+5])
		monkey.IfFalse, err = strconv.Atoi(falseline[5])
		if err != nil {
			log.Fatalf("line %q: %q", line, err)
		}

		i += 6
	}

	for i := 0; i < 20; i++ {
		swarm.Round()
	}
	swarm.Print()

	sort.Slice(swarm.Monkeys, func(i, j int) bool {
		return swarm.Monkeys[i].Inspects < swarm.Monkeys[j].Inspects
	})

	return swarm.Monkeys[len(swarm.Monkeys)-2].Inspects *
		swarm.Monkeys[len(swarm.Monkeys)-1].Inspects
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

	input, err := os.ReadFile("input")
	if err != nil {
		log.WithError(err).Fatal("could not read input")
	}
	log.Printf("input solution: %d", solution("input", input))
}
