package main

import (
	"bufio"
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
	o string
	s string
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
	pL := NewPacketParser(p.Left)
	pR := NewPacketParser(p.Right)

	lVal := pL.Next()
	rVal := pR.Next()

	if lVal == nil {
		fmt.Println("  - Left ran out of values")
		return true
	}
	if rVal == nil {
		fmt.Println("  - Right ran out of items")
		return false
	}
	if lVal.IsNumber() && rVal.IsNumber() {
		fmt.Printf("  - Compare %s vs %s\n", lVal, rVal)
		intL := lVal.IntVal()
		intR := rVal.IntVal()
		if intL > intR {
			fmt.Printf("   - %d > %d\n", intL, intR)
			return false
		}
		if intL < intR {
			fmt.Printf(" Early qualify!!!\n")
			return true
		}
		newPair := PacketPair{pL.s, pR.s, p.level}
		return newPair.IsValid()
	}

	if lVal.IsList() && rVal.IsList() {
		fmt.Println("Both are lists")
		newPair := PacketPair{lVal.String(), rVal.String(), p.level + 1}
		if newPair.IsValid() {
			if p.level == 0 {
				return true
			}
			if len(pL.s) > 0 && len(pL.s) == 0 {
				return false
			} else if strings.Index(lVal.String(), "[") < 0 && strings.Index(rVal.String(), "[") < 0 {
				// no more lists..
				return true
			} else {
				nextPair := PacketPair{pL.s, pR.s, p.level + 1}
				return nextPair.IsValid()
			}
		} else {
			return false
		}
	}

	if lVal.IsList() && rVal.IsNumber() {
		pL.Next()
		newPair := PacketPair{pL.s, pR.s, p.level}
		return newPair.IsValid()
	}

	if lVal.IsNumber() && rVal.IsList() {
		pR.Next()
		newPair := PacketPair{pL.s, pR.s, p.level}
		return newPair.IsValid()
	}

	if !pL.AnyRemain() {
		return true
	}
	reason = "Shouldn't have got here?"
	return false
}

func NewPacketParser(s string) *PacketParser {
	return &PacketParser{o: s, s: s}
}

func (r *PacketParser) AnyRemain() bool {
	return !(len(r.s) == 0)
}

func (r *PacketParser) Next() *PacketItem {
	if len(r.s) == 0 {
		return nil
	}
	if r.s == "[]" {
		return nil
	}
	if r.s[0] == ',' {
		// pop comma, try again
		r.s = r.s[1:]
		return r.Next()
	}
	if r.s[0] >= '0' && r.s[0] <= '9' {
		parts := strings.Split(r.s, ",")
		sv := strings.ReplaceAll(parts[0], "]", "")
		r.s = r.s[len(sv):]
		if len(r.s) > 0 && r.s[0] == ',' {
			r.s = r.s[1:]
		}
		return &PacketItem{sv, Number}
	}
	var level, endIndex int
	if r.s[0] == '[' {
		// Unwrap list
		for i, v := range r.s {
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
	if len(r.s) == 0 {
		return nil
	}
	if r.s == "[]" {
		return &PacketItem{"[]", List}
	}
	return &PacketItem{r.s[1:endIndex], List}
}

func (v *PacketItem) IsList() bool {
	return v.t == List
}
func (v *PacketItem) IsNumber() bool {
	return v.t == Number
}
func (v *PacketItem) IntVal() int {
	r, _ := strconv.Atoi(v.v)
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
