package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type ForestGrid [][]*Tree
type Cord struct {
	x int
	y int
}

type Tree struct {
	Height   int
	Explored bool
	ExposedX bool
	ExposedY bool
}

var (
	grid ForestGrid
)

func MakeRow(str string, row []*Tree) {
	for i, v := range str {
		h, _ := strconv.Atoi(string(v))
		row[i] = &Tree{h, false, false, false}
	}
}

func main() {
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
	grid.Search()
	fmt.Printf("%v", grid)
	exposed := grid.CountExposed()
	fmt.Printf("Part 1: %d trees are exposed\n", exposed)
}

func CordUp(c Cord) Cord {
	return Cord{c.x, c.y - 1}
}
func CordDown(c Cord) Cord {
	return Cord{c.x, c.y + 1}
}
func CordLeft(c Cord) Cord {
	return Cord{c.x - 1, c.y}
}
func CordRight(c Cord) Cord {
	return Cord{c.x + 1, c.y}
}

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
	return fmt.Sprintf("%d %b%v%v", t.Height, e, x, y)
}

func (t *Tree) IsExposed() bool {
	return t.ExposedX || t.ExposedY
}

func (g ForestGrid) CountExposed() int {
	var t int
	w := g.Width()
	h := g.Height()
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			if g[y][x].IsExposed() {
				t++
			}
		}
	}
	return t
}

func (g ForestGrid) Search() {
	w := g.Width()
	h := g.Height()
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			g.Explore(Cord{x, y})
		}
	}
}
func (g ForestGrid) Explore(c Cord) {
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
				this.Explored = true
				return
			}
			n := g.At(nc)
			if n.Height >= this.Height {
				break
			}
			if yPlane[idx] && n.ExposedY && this.Height > n.Height {
				this.ExposedY = true
				this.Explored = true
				return
			}
			if !yPlane[idx] && n.ExposedX && this.Height > n.Height {
				this.ExposedX = true
				this.Explored = true
				return
			}
			nc = f(nc)
		}
	}
	this.Explored = true
}
func (g ForestGrid) ValidCord(c Cord) bool {
	if c.x < 0 || c.y < 0 {
		return false
	}
	if c.x >= g.Width() {
		return false
	}
	if c.y >= g.Height() {
		return false
	}
	return true
}
func (g ForestGrid) At(c Cord) *Tree {
	return g[c.y][c.x]
}
func (g ForestGrid) Width() int {
	return len(g[0])
}
func (g ForestGrid) Height() int {
	return len(g)
}
func (g ForestGrid) String() string {
	var x string
	for _, row := range g {
		x += fmt.Sprintf("%v\n", row)
	}
	return x
}
