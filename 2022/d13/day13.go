package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type ItemType int

const (
	List ItemType = iota
	Number
)

type PacketPair struct {
	Left  string
	Right string
	level int
}

type PacketParser struct {
	Original  string `json:"original"`
	Current   string `json:"current"`
	Remainder string `json:"remainder"`
}

type PacketItem struct {
	v string
	t ItemType
}

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
		fmt.Printf("== Pair %d ==\n", i+1)
		if p.IsValid() {
			validPackets += (i + 1)
			fmt.Println("Y\n")
		} else {
			fmt.Println("N\n")
		}
	}

	fmt.Printf("Part1: %d\n", validPackets)
}

func (p PacketPair) IsValid() bool {
	fmt.Printf("- Compare %s vs %s\n", p.Left, p.Right)
	var valid bool = true
	var i int
	pL := NewPacketParser(p.Left)
	pR := NewPacketParser(p.Right)
	for {
		i++
		//fmt.Printf("Loop %d\n", i)
		/*
			if i > 10 {
				panic("too much")
			}
		*/
		if !valid {
			return false
		}
		lVal := pL.Next()
		rVal := pR.Next()

		fmt.Printf("%*s Compare %s vs %s\n", i, "-", lVal, rVal)

		if lVal == nil {
			fmt.Printf("%*s Left ran out of values\n", i, "-")
			return true
		}
		if rVal == nil {
			fmt.Printf("%*s Right ran out of items\n", i, "-")
			return false
		}
		if lVal.IsNumber() && rVal.IsNumber() {
			intL := lVal.IntVal()
			intR := rVal.IntVal()
			fmt.Printf("%*s Compare %d vs %d\n", i, "-", intL, intR)
			if intL > intR {
				fmt.Printf("%*s %d > %d\n", i, "-", intL, intR)
				fmt.Printf("%*s Right side smaller\n", i+1, "-")
				return false
			}
			if intL < intR {
				fmt.Printf("%*s Early qualify!!!\n", i, "-")
				return true
			}
			//newPair := PacketPair{pL.s, pR.s, p.level}
			//valid = newPair.IsValid()
			continue
		}

		if lVal.IsList() && rVal.IsList() {
			continue
		}
		if lVal.IsList() && rVal.IsNumber() {
			pL.Next()
			fmt.Println("rewrapped right")
			pR.Current = "[" + rVal.String() + "]"
			continue
		}

		if lVal.IsNumber() && rVal.IsList() {
			pR.Next()
			fmt.Println("rewrapped left")
			pL.Current = "[" + lVal.String() + "]"
			continue
		}

		if !pL.AnyRemain() {
			fmt.Println("%*s No items remain on left", i, "-")
			return true
		}
		reason = "Shouldn't have got here?"
		return false
	}
}

func NewPacketParser(s string) *PacketParser {
	return &PacketParser{Original: s, Current: s, Remainder: ""}
}

func (r *PacketParser) AnyRemain() bool {
	return !(len(r.Current) == 0)
}

func (r *PacketParser) String() string {
	s, err := json.MarshalIndent(*r, "", " ")
	if err != nil {
		fmt.Println("Error marshaling parser values:", err)
	}
	return string(s)
}

func (r *PacketParser) Next() *PacketItem {
	//fmt.Printf("%s\n", r)
	if len(r.Current) == 0 && len(r.Remainder) == 0 {
		return nil
	}
	if len(r.Current) == 0 && len(r.Remainder) > 0 {
		r.Current = r.Remainder
		r.Remainder = ""
	}
	if r.Current == "[]" {
		return nil
	}
	if r.Current[0] == ',' {
		// pop comma, try again
		r.Current = r.Current[1:]
		return r.Next()
	}
	if r.Current[0] >= '0' && r.Current[0] <= '9' {
		parts := strings.Split(r.Current, ",")
		sv := strings.ReplaceAll(parts[0], "]", "")
		r.Current = r.Current[len(sv):]
		if len(r.Current) > 0 && r.Current[0] == ',' {
			r.Current = r.Current[1:]
		}
		return &PacketItem{sv + "," + r.Current, Number}
	}
	var level, endIndex int
	if r.Current[0] == '[' {
		// Unwrap list
		for i, v := range r.Current {
			if i == 0 {
				continue
			}
			if v == '[' {
				level++
				continue
			}
			if v == ']' {
				if level == 0 {
					endIndex = i
					break
				} else {
					level--
				}
				continue
			}
		}
	}
	if len(r.Current) == 0 {
		return nil
	}
	if r.Current == "[]" {
		return &PacketItem{"[]", List}
	}
	if endIndex+1 != len(r.Current) {
		r.Remainder = r.Current[endIndex+1:]
		r.Current = r.Current[1:endIndex]
	} else {
		r.Current = r.Current[1:endIndex]
	}
	return &PacketItem{r.Current, List}
}

func (v *PacketItem) IsList() bool {
	return v.t == List
}
func (v *PacketItem) IsNumber() bool {
	return v.t == Number
}
func (v *PacketItem) IntVal() int {
	parts := strings.Split(v.v, ",")
	sv := strings.ReplaceAll(parts[0], "]", "")
	r, _ := strconv.Atoi(sv)
	return r
}
func (v *PacketItem) String() string {
	return v.v
}

/**
 * Might revisit...
 */
func (p PacketPair) IsValidFail() bool {
	fmt.Printf("Compare %s vs %s\n", p.Left, p.Right)
	if ItemsRemain(p.Left) && !ItemsRemain(p.Right) {
		reason = "Out of items on the right"
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
		n := PacketPair{l, r, 0}
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

	n := PacketPair{l, r, 0}
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
