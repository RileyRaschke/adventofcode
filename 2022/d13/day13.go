package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type PacketPair struct {
	Left  string
	Right string
}

type Item rune
type List string

var (
	pairs  []PacketPair = make([]PacketPair, 0)
	reason string
)

func main() {
	var i, validPackets int
	reader := bufio.NewReader(os.Stdin)
	p := PacketPair{}
	for {
		str, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		line := strings.TrimSpace(string(str))

		if line == "" {
			pairs = append(pairs, p)
			p = PacketPair{}
			continue
		}
		if i%2 == 0 {
			p.Left = line
		} else {
			p.Right = line
		}
		i++
	}
	pairs = append(pairs, p)

	for i, p := range pairs {
		fmt.Printf("\n%s\n", p.Left)
		if p.IsValid() {
			validPackets += (i + 1)
			fmt.Println("Y")
		} else {
			fmt.Println("N - ", reason)
		}
		fmt.Printf("%s\n\n", p.Right)
		reason = ""
	}

	fmt.Printf("Part1: %d\n", validPackets)
}

func (p PacketPair) IsValid() bool {
	if ItemsRemain(p.Left) && !ItemsRemain(p.Right) {
		return false
	}
	if len(p.Left) > 0 && len(p.Right) == 0 {
		reason = "Right side ran out of items"
		return false
	}

	if !ItemsRemain(p.Left) {
		return true
	}

	var l, r string
	var unwrapped bool = false
	if p.Right[0] == '[' {
		r = Unwrap(p.Right)
		unwrapped = true
	} else {
		r = p.Right
	}
	if p.Left[0] == '[' {
		l = Unwrap(p.Left)
		unwrapped = true
	} else {
		l = p.Left
	}
	if unwrapped {
		n := PacketPair{l, r}
		return n.IsValid()
	}

	tL := NextItem(p.Left, &l)
	tR := NextItem(p.Right, &r)
	if tL > tR {
		reason = fmt.Sprintf("%d > %d", tL, tR)
		return false
	}

	l = PruneItem(l)
	r = PruneItem(r)

	if ItemsRemain(l) && !ItemsRemain(r) {
		return false
	}

	if !ItemsRemain(l) {
		return true
	}

	n := PacketPair{l, r}
	return n.IsValid()
}

func Unwrap(s string) string {
	x := s[1:]
	x = x[:len(x)-1]
	return x
}

func NextItem(s string, out *string) int {
	if s[0] == ']' {
		return -1
	}
	parts := strings.Split(s, ",")
	sv := strings.ReplaceAll(parts[0], "]", "")
	i, err := strconv.Atoi(sv)
	if err != nil {
		panic(err)
	}
	z := s[len(sv):]
	*out = z
	return i
}

func PruneItem(s string) string {
	if len(s) == 0 {
		return ""
	}
	switch s[0] {
	case ',':
		return PruneItem(s[1:])
	case ']':
		return PruneItem(s[1:])
	default:
		return s
	}
}

func ItemsRemain(s string) bool {
	for _, r := range s {
		if r >= '0' && r <= '9' {
			return true
		}
	}
	return false
}
