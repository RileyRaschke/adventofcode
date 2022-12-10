package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

type Cord struct {
	x int
	y int
}

type Tree struct {
	Height      int
	Explored    bool
	ExposedX    bool
	ExposedY    bool
	ScenicScore int
}

type ForestGrid [][]*Tree

var (
	grid ForestGrid
)

func MakeRow(str string, row []*Tree) {
	for i, v := range str {
		h, _ := strconv.Atoi(string(v))
		row[i] = &Tree{h, false, false, false, 0}
	}
}

func main() {
	start := time.Now()
	grid = ForestGrid{}
	reader := bufio.NewReader(os.Stdin)
	for {
		str, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		line := strings.TrimSpace(string(str))
		width := len(line)
		row := make([]*Tree, width)
		MakeRow(line, row)
		grid = append(grid, row)
	}

	exposed, maxScore := grid.Explore()

	fmt.Printf("Part 1: %d trees are exposed\n", exposed)
	fmt.Printf("Part 2: Highest senic score is %d\n", maxScore)
	fmt.Printf("Total time: %s\n", time.Since(start))
}

func CordUp(c Cord) Cord    { return Cord{c.x, c.y - 1} }
func CordDown(c Cord) Cord  { return Cord{c.x, c.y + 1} }
func CordLeft(c Cord) Cord  { return Cord{c.x - 1, c.y} }
func CordRight(c Cord) Cord { return Cord{c.x + 1, c.y} }

func (t *Tree) IsExposed() bool { return t.ExposedX || t.ExposedY }

func (t *Tree) String() string {
	var e, x, y int = 0, 0, 0
	if t.Explored {
		e = 1
	}
	if t.ExposedX {
		x = 1
	}
	if t.ExposedY {
		y = 1
	}
	return fmt.Sprintf("%d-%d-%d%d%d", t.Height, t.ScenicScore, e, x, y)
}

func (g ForestGrid) Explore() (int, int) {
	var exposed int
	var maxScore int
	w := g.Width()
	h := g.Height()
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			c := Cord{x, y}
			//go func(cord Cord) { // too many threads.

			g.FindExposed(c)
			g.ScenicScores(c)

			t := g.At(c)

			if t.IsExposed() {
				exposed++
			}
			if t.ScenicScore > maxScore {
				maxScore = t.ScenicScore
			}
			//}(c)
		}
	}
	return exposed, maxScore
}

func (g ForestGrid) FindExposed(c Cord) {
	funcs := []func(Cord) Cord{CordUp, CordDown, CordLeft, CordRight}
	yPlane := []bool{true, true, false, false}
	this := g.At(c)
	for idx, f := range funcs {
		nc := f(c)
		for {
			if !g.ValidCord(nc) {
				if yPlane[idx] {
					this.ExposedY = true
				} else {
					this.ExposedX = true
				}
				break
			}
			n := g.At(nc)
			if n.Height >= this.Height {
				break
			}
			if yPlane[idx] && n.ExposedY && this.Height > n.Height {
				this.ExposedY = true
				this.Explored = true
				break
			}
			if !yPlane[idx] && n.ExposedX && this.Height > n.Height {
				this.ExposedX = true
				break
			}
			nc = f(nc)
		}
	}
	this.Explored = true
}

func (g ForestGrid) ScenicScores(c Cord) {
	funcs := []func(Cord) Cord{CordUp, CordDown, CordLeft, CordRight}
	this := g.At(c)
	for idx, f := range funcs {
		dist := 0
		nc := f(c)
		for {
			if !g.ValidCord(nc) {
				break
			}
			dist++
			n := g.At(nc)
			if n.Height >= this.Height {
				break
			}
			nc = f(nc)
		}
		if idx == 0 {
			this.ScenicScore = dist
		} else {
			this.ScenicScore *= dist
		}
	}
}

func (g ForestGrid) At(c Cord) *Tree { return g[c.y][c.x] }
func (g ForestGrid) Width() int      { return len(g[0]) }
func (g ForestGrid) Height() int     { return len(g) }

func (g ForestGrid) ValidCord(c Cord) bool {
	if c.x < 0 || c.y < 0 || c.x >= g.Width() || c.y >= g.Height() {
		return false
	}
	return true
}

func (g ForestGrid) String() string {
	var x string
	for _, row := range g {
		x += fmt.Sprintf("%v\n", row)
	}
	return x
}
