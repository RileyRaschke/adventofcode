package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Monkey struct {
	Items       []int
	Operation   string
	Test        int
	OnSuccess   int
	OnFail      int
	Inspections int
}

type MonkeyGroup []*Monkey

var (
	monkeys     MonkeyGroup
	reduceWorry bool = true
)

func main() {
	m := &Monkey{}
	reader := bufio.NewReader(os.Stdin)
	for {
		str, _, err := reader.ReadLine()
		if err == io.EOF {
			monkeys = append(monkeys, m)
			break
		}
		line := strings.TrimSpace(string(str))
		if line == "" {
			monkeys = append(monkeys, m)
			m = &Monkey{}
			continue
		}
		parts := strings.Split(line, ":")
		switch parts[0] {
		case "Starting items":
			items := []int{}
			sitems := strings.Split(parts[1], ",")
			for _, v := range sitems {
				v = strings.TrimSpace(v)
				w, _ := strconv.Atoi(v)
				items = append(items, w)
			}
			m.Items = items
			break
		case "Operation":
			m.Operation = strings.TrimSpace(parts[1])
			break
		case "Test":
			fmt.Sscanf(strings.TrimSpace(parts[1]), "divisible by %d", &(m.Test))
			break
		case "If true":
			var dest int
			fmt.Sscanf(parts[1], " throw to monkey %d", &dest)
			m.OnSuccess = dest
			break
		case "If false":
			var dest int
			fmt.Sscanf(parts[1], " throw to monkey %d", &dest)
			m.OnFail = dest
			break
		default:
			continue
		}
	}
	fmt.Printf("%v\n", monkeys)
	monkeys_Part2 := monkeys.DeepCopy()
	// Part1
	monkeys.RunRounds(20)
	var mb int = 1
	for _, ti := range monkeys.TopInspectors(2) {
		mb *= ti.Inspections
	}
	fmt.Printf("%v\n", monkeys)
	fmt.Printf("Part1: %d\n", mb)

	// Part2
	reduceWorry = false
	monkeys_Part2.RunRounds(10000)
	mb = 1
	for _, ti := range monkeys_Part2.TopInspectors(2) {
		mb *= ti.Inspections
	}
	//fmt.Printf("%v\n", monkeys_Part2)
	fmt.Printf("Part2: %d\n", mb)
}

func (p MonkeyGroup) RunRounds(n int) {
	for i := 0; i < n; i++ {
		p.RunRound()
		if i+1 == 1 || i+1 == 20 || (i+1)%1000 == 0 {
			p.Dump(i + 1)
		}
	}
}

func (p MonkeyGroup) RunRound() {
	cm := p.CommonMultiple()
	for _, m := range p {
		for {
			if len(m.Items) == 0 {
				break
			}
			destMonkey, value := m.ProcessItem(cm)
			p[destMonkey].Items = append(p[destMonkey].Items, value)
		}
	}
}

func (p MonkeyGroup) TopInspectors(c int) MonkeyGroup {
	var gc MonkeyGroup = make([]*Monkey, len(p))
	for i, m := range p {
		gc[i] = m
	}
	sort.Slice(gc, func(i, j int) bool {
		return gc[i].Inspections < gc[j].Inspections
	})
	return gc[len(p)-c:]
}

func (p MonkeyGroup) CommonMultiple() int {
	var x int = 1
	for _, m := range p {
		x *= m.Test
	}
	return x
}

func (p MonkeyGroup) Dump(round int) {
	fmt.Printf("== After round %d ==\n", round)
	for i, m := range p {
		fmt.Printf("Monkey %d, inspected items %d times.\n", i, m.Inspections)
	}
	fmt.Println("")
}

func (p MonkeyGroup) DeepCopy() MonkeyGroup {
	var gc MonkeyGroup = make([]*Monkey, len(p))
	for i, m := range p {
		gc[i] = &Monkey{m.Items, m.Operation, m.Test, m.OnSuccess, m.OnFail, m.Inspections}
	}
	return gc
}

func (m *Monkey) ProcessItem(cm int) (dest, val int) {
	m.Inspections++
	item := m.Items[0]
	m.Items = m.Items[1:]
	opArgs := strings.Split(m.Operation, " ")
	var x, y, n int
	if opArgs[2] == "old" {
		x = item
	} else {
		v, _ := strconv.Atoi(opArgs[2])
		x = v
	}
	if opArgs[4] == "old" {
		y = item
	} else {
		v, _ := strconv.Atoi(opArgs[4])
		y = v
	}
	switch opArgs[3] {
	case "*":
		n = x * y
		break
	case "+":
		n = x + y
		break
	case "-":
		n = x - y
		break
	case "/":
		n = x / y
		break
	}
	var tv int = n
	if reduceWorry {
		tv = n / 3
	} else {
		tv = n % cm
	}

	if tv%m.Test == 0 {
		return m.OnSuccess, tv
	}
	return m.OnFail, tv
}

func (m *Monkey) String() string {
	s := "\n{\n"
	s += fmt.Sprintf("\tItems: %v\n", m.Items)
	s += fmt.Sprintf("\tOp: %s\n", m.Operation)
	s += fmt.Sprintf("\tTest: %d\n", m.Test)
	s += fmt.Sprintf("\tOnSuccess: Throw to monkey %d\n", m.OnSuccess)
	s += fmt.Sprintf("\tOnFail: Throw to monkey %d\n", m.OnFail)
	s += fmt.Sprintf("\tInspections: %d\n", m.Inspections)
	s += "}"
	return s
}
