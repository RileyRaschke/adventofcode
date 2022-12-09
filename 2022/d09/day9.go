package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"strings"
)

type Cord struct {
	X int
	Y int
}

type Rope struct {
	Visits map[Cord]int
	Knots  []*Cord
}

func NewRope(length int) *Rope {
	knots := make([]*Cord, length)
	for i := range knots {
		knots[i] = &Cord{0, 0}
	}
	r := &Rope{make(map[Cord]int), knots}
	r.Visits[*r.Knots[len(r.Knots)-1]] = 1
	return r
}

func main() {
	rope1 := NewRope(2)
	rope2 := NewRope(10)

	reader := bufio.NewReader(os.Stdin)
	for {
		str, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		line := strings.TrimSpace(string(str))
		rope1.Move(line)
		rope2.Move(line)
	}
	fmt.Printf("%s", rope1)
	fmt.Printf("%s", rope2)
	fmt.Printf("Part1: %d\n", len(rope1.Visits))
	fmt.Printf("Part2: %d\n", len(rope2.Visits))
}

func CordsAdjacent(lead *Cord, follow *Cord) (bool, string) {
	adjacent := true
	xDist := lead.X - follow.X
	yDist := lead.Y - follow.Y
	followMove := ""
	if math.Abs(float64(xDist))+math.Abs(float64(yDist)) == 3 {
		if yDist > 0 {
			yDist++
		} else {
			yDist--
		}
		if xDist > 0 {
			xDist++
		} else {
			xDist--
		}
	}
	if xDist > 1 {
		followMove += "R"
		adjacent = false
	} else if xDist < -1 {
		followMove += "L"
		adjacent = false
	}
	if yDist > 1 {
		followMove += "U"
		adjacent = false
	} else if yDist < -1 {
		followMove += "D"
		adjacent = false
	}
	return adjacent, followMove
}

func (c *Cord) MoveCord(dir string) {
	for _, m := range dir {
		switch string(m) {
		case "U":
			c.Y++
			break
		case "D":
			c.Y--
			break
		case "R":
			c.X++
			break
		case "L":
			c.X--
			break
		}
	}
}

func (r *Rope) Move(m string) {
	var dir string
	var moveCount int
	fmt.Sscanf(m, "%s %d", &dir, &moveCount)
	for i := 0; i < moveCount; i++ {
		r.Head().MoveCord(dir)
		for k := 0; k+1 < len(r.Knots); k++ {
			adjacent, followMove := CordsAdjacent(r.Knots[k], r.Knots[k+1])
			if !adjacent {
				r.Knots[k+1].MoveCord(followMove)
				if r.Knots[k+1] == r.Tail() {
					if _, ok := r.Visits[*r.Tail()]; ok {
						r.Visits[*r.Tail()]++
					} else {
						r.Visits[*r.Tail()] = 1
					}
				}
			}
		}
		//fmt.Printf("H%s T%s %s\n", r.Head, r.Tail, tailMove)
	}
}

func (c *Cord) String() string {
	return fmt.Sprintf("(%d,%d)", c.X, c.Y)
}

func (r *Rope) String() string {
	grid := ""
	var xMin, xMax, yMin, yMax int
	for k, _ := range r.Visits {
		if k.X > xMax {
			xMax = k.X
		}
		if k.X < xMin {
			xMin = k.X
		}
		if k.Y > yMax {
			yMax = k.Y
		}
		if k.Y < yMin {
			yMin = k.Y
		}
	}
	height := yMax - yMin + 2
	width := xMax - xMin + 2
	for y := height - 1; y >= 0; y-- {
		for x := 0; x < width; x++ {
			c := Cord{x + xMin, y + yMin}
			if val, ok := r.Visits[c]; ok {
				grid += fmt.Sprintf("%d", val)
			} else {
				grid += "."
			}
		}
		grid += "\n"
	}
	return grid
}

func (r *Rope) Head() *Cord {
	return r.Knots[0]
}

func (r *Rope) Tail() *Cord {
	return r.Knots[len(r.Knots)-1]
}
