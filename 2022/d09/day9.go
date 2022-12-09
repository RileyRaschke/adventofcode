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
	Head   *Cord
	Tail   *Cord
}

var (
	rope *Rope
)

func NewRope() *Rope {
	r := &Rope{make(map[Cord]int), &Cord{0, 0}, &Cord{0, 0}}
	r.Visits[*r.Tail] = 1
	return r
}

func main() {
	rope = NewRope()
	reader := bufio.NewReader(os.Stdin)
	for {
		str, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		line := strings.TrimSpace(string(str))
		rope.MoveHead(line)
	}
	fmt.Printf("%s", rope)
	fmt.Printf("Part1: %d\n", len(rope.Visits))
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

func (r *Rope) MoveHead(m string) {
	var dir string
	var moveCount int
	fmt.Sscanf(m, "%s %d", &dir, &moveCount)
	for i := 0; i < moveCount; i++ {
		r.Head.MoveCord(dir)
		adjacent, tailMove := r.HeadTailAdjacent()
		if !adjacent {
			r.Tail.MoveCord(tailMove)

			if _, ok := r.Visits[*r.Tail]; ok {
				r.Visits[*r.Tail]++
			} else {
				r.Visits[*r.Tail] = 1
			}
		}
		//fmt.Printf("H%s T%s %s\n", r.Head, r.Tail, tailMove)
	}
}

func (r *Rope) HeadTailAdjacent() (bool, string) {
	adjacent := true
	xDist := r.Head.X - r.Tail.X
	yDist := r.Head.Y - r.Tail.Y
	tailMove := ""
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
		tailMove += "R"
		adjacent = false
	} else if xDist < -1 {
		tailMove += "L"
		adjacent = false
	}
	if yDist > 1 {
		tailMove += "U"
		adjacent = false
	} else if yDist < -1 {
		tailMove += "D"
		adjacent = false
	}
	return adjacent, tailMove
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
