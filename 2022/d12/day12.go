package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
)

type Cord struct {
	X int
	Y int
}

type MapNode struct {
	Position Cord
	Height   int
	Val      string
	Dir      string
	Dead     bool
	Step     int
}

type HeightMap struct {
	Grid     [][]*MapNode
	StartPos *MapNode
	EndPos   *MapNode
	Paths    []Path
}

type Path []*MapNode

func main() {
	var x, y int
	g := &HeightMap{}
	row := []*MapNode{}
	reader := bufio.NewReader(os.Stdin)
	for {
		z, _, err := reader.ReadRune()
		if err == io.EOF {
			break
		}
		if z == '\n' {
			g.Grid = append(g.Grid, row)
			row = []*MapNode{}
			y++
			x = 0
			continue
		}
		var n *MapNode
		if z == 'S' {
			n = NewMapNode(x, y, 'a')
			g.StartPos = n
		} else if z == 'E' {
			n = NewMapNode(x, y, 'z')
			g.EndPos = n
		} else {
			n = NewMapNode(x, y, z)
		}
		row = append(row, n)
		x++
	}
	fmt.Printf("%s\n", g)

	path := g.ShortestPath()
	fmt.Printf("\nPart1: %d - %v\n", len(path)-1, path)
}

func NewMapNode(x int, y int, z rune) *MapNode {
	return &MapNode{Cord{x, y}, int(z - 'a'), string(z), "", false, 0}
}

func (m *HeightMap) ShortestPath() Path {
	paths := m.FindPaths()
	sort.Slice(paths, func(i, j int) bool {
		return len(paths[i]) < len(paths[j])
	})
	fmt.Println("")
	for _, p := range paths {
		fmt.Printf("(%d) %v\n", len(p), p)
	}
	return paths[0]
}

func (m *HeightMap) FindPaths() []Path {
	m.ExplorePath(Path{m.StartPos})
	return m.Paths
}

func (m *HeightMap) ExplorePath(p Path) {
	var this *MapNode
	pl := len(p)
	if pl == 1 {
		this = p[0]
	} else {
		this = p[pl-1]
	}

	if this.Step == 0 {
		this.Step = pl
	} else if this.Step > 0 && this.Step <= pl {
		return
	}

	//fmt.Printf("%v\n", p)
	for _, n := range m.Neighbors(this) {
		if n.Step > 0 && n.Step <= this.Step {
			continue
		}
		if p.Contains(n) {
			continue
		}
		if n.Height-this.Height > 1 {
			continue
		}
		if n.Dead {
			continue
		}
		if n == m.EndPos {
			m.Paths = append(m.Paths, append(p, n))
			continue
		}
		m.ExplorePath(append(p, n))
	}
}

func (p Path) Contains(n *MapNode) bool {
	for _, tn := range p {
		if tn == n {
			return true
		}
	}
	return false
}

func (m *HeightMap) Neighbors(n *MapNode) []*MapNode {
	dirs := []func(Cord) Cord{CordUp, CordDown, CordLeft, CordRight}
	neighbors := []*MapNode{}
	for _, f := range dirs {
		x := f(n.Position)
		if m.ValidCord(x) {
			neighbors = append(neighbors, m.At(x))
		}
	}
	return neighbors
}

func (m *HeightMap) At(c Cord) *MapNode {
	return m.Grid[c.Y][c.X]
}

func (m *HeightMap) ValidCord(c Cord) bool {
	if c.X < 0 || c.Y < 0 || c.X >= len(m.Grid[0]) || c.Y >= len(m.Grid) {
		return false
	}
	return true
}

func CordUp(c Cord) Cord    { return Cord{c.X, c.Y - 1} }
func CordDown(c Cord) Cord  { return Cord{c.X, c.Y + 1} }
func CordLeft(c Cord) Cord  { return Cord{c.X - 1, c.Y} }
func CordRight(c Cord) Cord { return Cord{c.X + 1, c.Y} }

func (m *HeightMap) String() string {
	s := ""
	for _, row := range m.Grid {
		for _, n := range row {
			s += fmt.Sprintf("%2d ", n.Height)
		}
		s += "\n"
	}
	return s
}

func (n *MapNode) String() string {
	return fmt.Sprintf("(%d,%d)", n.Position.X, n.Position.Y)
}
