package main

import (
	"bufio"
	"encoding/json"
	"flag"
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
	prior     *PacketParser
	Current   string `json:"current"`
	Remainder string `json:"remainder"`
}

type PacketItem struct {
	v string
	t ItemType
}

var (
	pairs      []PacketPair = make([]PacketPair, 0)
	PrintDebug bool         = false
)

func main() {

	flag.BoolVar(&PrintDebug, "d", false, "Print debug info")
	flag.Parse()

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
	var i int = 1
	pL := NewPacketParser(p.Left)
	pR := NewPacketParser(p.Right)
	for {
		pad := i*2 - 1
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

		//fmt.Printf("%*s Compare %s vs %s\n", pad, "-", lVal, rVal)

		if lVal == nil {
			fmt.Printf("%*s Left ran out of items\n", pad, "-")
			return true
		}
		if rVal == nil {
			fmt.Printf("%*s Right ran out of items\n", pad, "-")
			return false
		}
		if lVal.IsNumber() && rVal.IsNumber() {
			intL := lVal.IntVal()
			intR := rVal.IntVal()
			fmt.Printf("%*s Compare %d vs %d\n", pad, "-", intL, intR)
			if intL > intR {
				fmt.Printf("%*s %d > %d\n", i*2, "-", intL, intR)
				fmt.Printf("%*s Right side smaller\n", pad+2, "-")
				return false
			}
			if intL < intR {
				fmt.Printf("%*s Left side smaller\n", pad+2, "-")
				return true
			}
			continue
		}

		if lVal.IsList() && rVal.IsList() {
			fmt.Printf("%*s Compare %s vs %s\n", pad, "-", lVal, rVal)
			i++
			continue
		}
		if lVal.IsList() && rVal.IsNumber() {
			i--
			pR.Rollback()
			continue
		}

		if lVal.IsNumber() && rVal.IsList() {
			i--
			pL.Rollback()
			continue
		}

		if !pL.AnyRemain() {
			fmt.Println("%*s No items remain on left", pad, "-")
			return true
		}
		fmt.Println("Shouldn't have got here?")
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
	//s, err := json.MarshalIndent(r, "", " ")
	s, err := json.Marshal(r)
	if err != nil {
		fmt.Println("Error marshaling parser values:", err)
	}
	return string(s)
}
func (r *PacketParser) Rollback() {
	if PrintDebug {
		fmt.Printf("RBP     %s\n", r)
	}
	r.Current = r.prior.Current
	r.Remainder = r.prior.Remainder
	r.prior = r.prior.prior
	if PrintDebug {
		fmt.Printf("RBA     %s\n", r)
	}
}

func (r *PacketParser) Next() *PacketItem {
	if PrintDebug {
		fmt.Printf("\t%s\n", r)
	}
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
		// remove comma, try again
		r.Current = r.Current[1:]
		return r.Next()
	}
	if r.Current[0] >= '0' && r.Current[0] <= '9' {
		parts := strings.Split(r.Current, ",")
		sv := strings.ReplaceAll(parts[0], "]", "")
		r.prior = r.Clone()
		r.Current = r.Current[len(sv):]
		if len(r.Current) > 0 && r.Current[0] == ',' {
			r.Current = r.Current[1:]
		}
		if r.Current == "" {
			return &PacketItem{sv, Number}
		} else {
			return &PacketItem{sv + "," + r.Current, Number}
		}
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
	r.prior = r.Clone()
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

func (r *PacketParser) Clone() *PacketParser {
	return &PacketParser{r.Original, r.prior, r.Current, r.Remainder}
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
