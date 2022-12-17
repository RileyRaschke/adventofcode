package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
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
	Depth     int    `json:"depth"`
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
	log.SetFlags(0)
	flag.BoolVar(&PrintDebug, "d", false, "Print debug info")
	flag.Parse()

	var i, validPackets int
	reader := bufio.NewReader(os.Stdin)
	p := PacketPair{}
	for {
		str, _, err := reader.ReadLine()
		if err == io.EOF {
			pairs = append(pairs, p)
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

	for i, p := range pairs {
		log.Printf("== Pair %d ==\n", i+1)
		if p.IsValid() {
			validPackets += (i + 1)
			log.Println("Y\n")
		} else {
			fmt.Printf("%s\n%s\n\n", p.Left, p.Right)
			log.Println("N\n")
		}
	}

	log.Printf("Part1: %d\n", validPackets)
}

func NewPacketParser(s string) *PacketParser {
	return &PacketParser{Original: s, prior: nil, Current: s, Remainder: ""}
}

func (p PacketPair) IsValid() bool {
	log.Printf("- Compare %s vs %s\n", p.Left, p.Right)
	var valid bool = true
	var i int = 1
	pL := NewPacketParser(p.Left)
	pR := NewPacketParser(p.Right)
	for {
		pad := i*2 - 1
		if !valid {
			return false
		}
		lVal := pL.Next()
		rVal := pR.Next()

		if lVal == nil {
			log.Printf("%*s Left ran out of items\n", pad, "-")
			return true
		}
		if rVal == nil {
			log.Printf("%*s Right ran out of items\n", pad, "-")
			return false
		}
		if lVal.IsNumber() && rVal.IsNumber() {
			/*
				if pL.Depth > pR.Depth {
					log.Printf("%*s Left ran out of items\n", pad, "-")
					return true
				}
			*/
			intL := lVal.IntVal()
			intR := rVal.IntVal()
			log.Printf("%*s Int Compare %d vs %d\n", pad, "-", intL, intR)
			if intL > intR {
				log.Printf("%*s %d > %d\n", i*2, "-", intL, intR)
				log.Printf("%*s Right side smaller\n", pad+2, "-")
				return false
			}
			if intL < intR {
				log.Printf("%*s Left side smaller\n", pad+2, "-")
				return true
			}
			continue
		}

		if lVal.IsList() && rVal.IsList() {
			log.Printf("%*s List Compare %s vs %s\n", pad, "-", lVal, rVal)
			/*
				if len(lVal.String()) >= 2 && lVal.String()[:2] == "[]" && pL.Depth == pR.Depth {
					log.Printf("%*s Left ran out of items\n", pad, "-")
					return true
				}
			*/
			/*
				if lVal.String() == "" && len(rVal.String()) > 0 {
					log.Printf("%*s Left ran out of items\n", pad, "-")
					return true
				}
			*/
			/*
				if lVal.String()[0] != '[' && rVal.String()[0] == '[' {
					log.Printf("%*s Left ran out of items\n", pad, "-")
					return true
				}
			*/
			i++
			continue
		}
		if lVal.IsList() && rVal.IsNumber() {
			/*
				if lVal.String() == "" {
					log.Printf("%*s Left ran out of items\n", pad, "-")
					return true
				}
			*/
			i--
			pR.Rollback()
			continue
		}

		if lVal.IsNumber() && rVal.IsList() {
			i--
			pL.Rollback()
			continue
		}

		log.Println("Shouldn't have got here?")
		return false
	}
}

func (r *PacketParser) Next() *PacketItem {
	if PrintDebug {
		log.Printf("\t%s\n", r)
	}
	if len(r.Current) == 0 && len(r.Remainder) == 0 {
		return nil
	}
	if len(r.Current) == 0 && len(r.Remainder) > 0 {
		r.Depth--
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
		r.Depth++
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

func (r *PacketParser) Rollback() {
	if PrintDebug {
		log.Printf("RBP     %s\n", r)
	}
	r.Current = r.prior.Current
	r.Remainder = r.prior.Remainder
	r.prior = r.prior.prior
	if PrintDebug {
		log.Printf("RBA     %s\n", r)
	}
}

func (r *PacketParser) String() string {
	s, err := json.Marshal(r)
	if err != nil {
		log.Println("Error marshaling parser values:", err)
	}
	return string(s)
}

func (v *PacketItem) IntVal() int {
	parts := strings.Split(v.v, ",")
	sv := strings.ReplaceAll(parts[0], "]", "")
	r, _ := strconv.Atoi(sv)
	return r
}

func (r *PacketParser) Clone() *PacketParser {
	return &PacketParser{r.Depth, r.Original, r.prior, r.Current, r.Remainder}
}

func (v *PacketItem) IsList() bool {
	return v.t == List
}

func (v *PacketItem) IsNumber() bool {
	return v.t == Number
}

func (v *PacketItem) String() string {
	return v.v
}
