package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"sort"

	"github.com/gookit/color"
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
	Cost     int
	Utility  float64
}

type HeightMap struct {
	Grid         [][]*MapNode
	StartPos     *MapNode
	EndPos       *MapNode
	LineOfSight  float64
	Frontiers    []Path
	ShortestPath Path
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
	path := g.FindShortestPath()
	fmt.Printf("\nPart1: %d - %v\n", len(path)-1, path)
}

func NewMapNode(x int, y int, z rune) *MapNode {
	return &MapNode{Cord{x, y}, int(z - 'a'), string(z), "", false, 0, -1}
}

func (m *HeightMap) FindShortestPath() Path {
	m.ShortestPath = nil
	m.CacheUtility()
	m.FindPaths()
	return m.ShortestPath
}

func (m *HeightMap) CacheUtility() {
	m.LineOfSight = m.Distance3d(m.StartPos.Position)
	for _, row := range m.Grid {
		for _, n := range row {
			if n == m.EndPos {
				n.Utility = float64(math.MaxInt64)
			} else {
				n.Utility = m.CalcUtility(n)
			}
		}
	}
}

func (m *HeightMap) FindPaths() Path {
	neighbors := m.SortedNeighbors(m.StartPos)
	for _, n := range neighbors {
		n.Cost = 2
		p := Path{m.StartPos, n}
		m.AddFrontierPath(p)
	}
	for {
		found := m.ExplorePath()
		if found {
			break
		}
	}
	return m.ShortestPath
}

func (m *HeightMap) AddFrontierPath(p Path) {
	last := p[len(p)-1]
	for i, q := range m.Frontiers {
		if pos := q.Contains(last); pos >= 0 {
			if pos > len(p) {
				m.Frontiers[i] = p
				last.Cost = len(p)
				return
			}
			return
		}

	}
	m.Frontiers = append(m.Frontiers, p)
}

func (m *HeightMap) ExplorePath() bool {
	var this *MapNode
	p := m.NextPath()
	if p == nil {
		return true // shouldn't happen?
	}
	pl := len(p)
	if pl == 1 {
		this = p[0]
	} else {
		this = p[pl-1]
	}
	if this == m.EndPos {
		m.AddPath(p)
		return true
	}
	if this.Cost == 0 {
		this.Cost = pl
	}
	if this.Cost < pl {
		return false
	}
	if len(m.Grid) < 20 {
		fmt.Printf("%v\n", p)
	}
	next := m.SortedNeighbors(this)
	if len(next) == 0 {
		//this.Dead = true
		return false
	}
	for _, n := range next {
		if n == m.EndPos {
			m.AddPath(append(p, n))
			return true
		} else if p.Contains(n) >= 0 {
			continue
		} else {
			m.AddFrontierPath(append(p, n))
		}
	}
	return false
}

func (m *HeightMap) SortedNeighbors(curr *MapNode) []*MapNode {
	dirs := []func(Cord) Cord{CordRight, CordUp, CordLeft, CordDown}
	neighbors := []*MapNode{}
	for _, f := range dirs {
		x := f(curr.Position)
		if !m.ValidCord(x) {
			continue
		}
		n := m.At(x)
		if n.Height > curr.Height+1 {
			continue
		}
		if n.Cost > 0 && n.Cost < curr.Cost {
			continue
		}
		if !n.Dead {
			neighbors = append(neighbors, n)
		}
	}
	sort.Slice(neighbors, func(i, j int) bool {
		return neighbors[i].Utility > neighbors[j].Utility
	})
	//fmt.Printf("%v\n", neighbors)
	return neighbors
	/*
		if len(neighbors) > 0 {
			return []*MapNode{neighbors[0]}
		}
		return neighbors
	*/
}

func (m *HeightMap) NextPath() Path {
	var p Path
	if len(m.Frontiers) == 0 {
		return nil
	} else if len(m.Frontiers) == 1 {
		p = m.Frontiers[0]
		m.Frontiers = []Path{}
		return p
	}

	sort.Slice(m.Frontiers, func(i, j int) bool {
		x := m.Frontiers[i]
		y := m.Frontiers[j]
		iLast := x[len(x)-1]
		jLast := y[len(y)-1]
		//return iLast.Utility/float64(len(x)) > jLast.Utility/float64(len(y))
		return iLast.Cost < jLast.Cost
	})
	fmt.Println("")
	for _, f := range m.Frontiers {
		fmt.Printf("%s\n", m.DrawPath(f))
	}
	// pop front
	p, m.Frontiers = m.Frontiers[0], m.Frontiers[1:len(m.Frontiers)]
	return p
}

func (m *HeightMap) AddPath(p Path) {
	//m.Paths = append(m.Paths, p)
	if m.ShortestPath == nil || len(p) < len(m.ShortestPath) {
		m.ShortestPath = p
		fmt.Printf("%s\n", m)
		fmt.Printf("%s\n", m.PathString())
		fmt.Printf("Current shortest: %d LoS: %0.3f\n", len(p)-1, m.LineOfSight)
	}
}

func (p Path) Contains(n *MapNode) int {
	for i, tn := range p {
		if tn == n {
			return i + 1
		}
	}
	return -1
}
func (m *HeightMap) CalcUtility(n *MapNode) float64 {
	/*
		u := m.LineOfSight / m.Distance3d(n.Position)
		//u = u / m.MinElevationDelta(n.Position)
		return u
	*/
	return 1
}

func (m *HeightMap) Distance2d(c Cord) float64 {
	//s := m.At(c)
	d := m.EndPos
	maxX := math.Max(float64(c.X), float64(d.Position.X))
	minX := math.Min(float64(c.X), float64(d.Position.X))
	maxY := math.Max(float64(c.Y), float64(d.Position.Y))
	minY := math.Min(float64(c.Y), float64(d.Position.Y))
	xd := maxX - minX
	yd := maxY - minY
	return math.Sqrt(math.Pow(xd, 2) + math.Pow(yd, 2))
}

func (m *HeightMap) Distance3d(c Cord) float64 {
	s := m.At(c)
	d := m.EndPos
	maxX := math.Max(float64(c.X), float64(d.Position.X))
	minX := math.Min(float64(c.X), float64(d.Position.X))
	maxY := math.Max(float64(c.Y), float64(d.Position.Y))
	minY := math.Min(float64(c.Y), float64(d.Position.Y))
	xd := maxX - minX
	yd := maxY - minY
	zd := float64(d.Height - s.Height)
	return math.Sqrt(math.Pow(xd, 2) + math.Pow(yd, 2) + math.Pow(zd, 2))
}

func (m *HeightMap) MinElevationDelta(c Cord) float64 {
	//s := m.At(c)
	d := m.EndPos
	maxX := int(math.Max(float64(c.X), float64(d.Position.X)))
	minX := int(math.Min(float64(c.X), float64(d.Position.X)))
	maxY := int(math.Max(float64(c.Y), float64(d.Position.Y)))
	minY := int(math.Min(float64(c.Y), float64(d.Position.Y)))
	var xFirst, yFirst float64
	for x := minX + 1; x <= maxX; x++ {
		xFirst += math.Abs(float64(m.Grid[c.Y][x].Height - m.Grid[c.Y][x-1].Height))
	}
	for y := minY + 1; y <= maxY; y++ {
		xFirst += math.Abs(float64(m.Grid[y][d.Position.X].Height - m.Grid[y-1][d.Position.X].Height))
	}
	for y := minY + 1; y <= maxY; y++ {
		yFirst += math.Abs(float64(m.Grid[y][c.X].Height - m.Grid[y-1][c.X].Height))
	}
	for x := minX + 1; x <= maxX; x++ {
		yFirst += math.Abs(float64(m.Grid[d.Position.Y][x].Height - m.Grid[d.Position.Y][x-1].Height))
	}
	return math.Min(xFirst, yFirst)
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

func (m *HeightMap) PathString() string {
	return m.DrawPath(m.ShortestPath)
}

func (m *HeightMap) String() string {
	return m.MarkPathString(m.ShortestPath)
}
func (m *HeightMap) DrawPath(sp Path) string {
	s := fmt.Sprintf("Len=%d\n", len(sp))
	pathMap := make(map[*MapNode]string, len(sp))
	if len(sp) > 0 {
		for i, n := range sp {
			if i+1 == len(sp) || i == 0 {
				continue
			}
			next := sp[i+1]
			if next.Position.X > n.Position.X {
				pathMap[n] = ">"
			} else if next.Position.X < n.Position.X {
				pathMap[n] = "<"
			} else if next.Position.Y > n.Position.Y {
				pathMap[n] = "v"
			} else {
				pathMap[n] = "^"
			}
		}
	}
	nfmt := "%2d "
	for _, row := range m.Grid {
		for _, n := range row {
			if n == m.StartPos {
				s += color.White.Render(color.Bold.Render(fmt.Sprintf(nfmt, n.Height)))
			} else if n == m.EndPos {
				s += color.Green.Render(color.Bold.Render(fmt.Sprintf(nfmt, n.Height)))
			} else if val, ok := pathMap[n]; ok {
				s += color.Yellow.Render(color.Bold.Render(fmt.Sprintf("%2s ", val)))
			} else {
				s += fmt.Sprintf("%2d ", n.Height)
			}
		}
		s += "\n"
	}
	return s
}

func (m *HeightMap) MarkPathString(sp Path) string {
	s := ""
	pathMap := make(map[*MapNode]bool, len(sp))
	if len(sp) > 0 {
		for _, n := range sp {
			pathMap[n] = true
		}
	}
	nfmt := "%2d "
	for _, row := range m.Grid {
		for _, n := range row {
			if n == m.StartPos {
				s += color.White.Render(color.Bold.Render(fmt.Sprintf(nfmt, n.Height)))
			} else if n == m.EndPos {
				s += color.Green.Render(color.Bold.Render(fmt.Sprintf(nfmt, n.Height)))
			} else if _, ok := pathMap[n]; ok {
				s += color.Yellow.Render(color.Bold.Render(fmt.Sprintf(nfmt, n.Height)))
			} else {
				s += fmt.Sprintf("%2d ", n.Height)
			}
		}
		s += "\n"
	}
	return s
}

func (n *MapNode) String() string {
	//return fmt.Sprintf("P(%d,%d) U(%0.5f)", n.Position.X, n.Position.Y, n.Utility)
	return fmt.Sprintf("(%d,%d)", n.Position.X, n.Position.Y)
}
