package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"os"
	"sort"
	"time"

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
	Seen     bool
	LosDist  float64
}

type Path []*MapNode

type HeightMap struct {
	Grid         [][]*MapNode
	StartPos     *MapNode
	EndPos       *MapNode
	LineOfSight  float64
	Frontiers    map[Cord][]Path
	ShortestPath Path
	FoundPaths   []Path
}

var (
	DumpPaths    bool = false
	DumpFrontier bool = false
)

func main() {
	tStart := time.Now()
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
	path := g.FindShortestPath()
	runtime := time.Since(tStart)
	fmt.Printf("%s\n", g)
	fmt.Printf("Found %d paths\n", len(g.FoundPaths))
	fmt.Printf("Part1: %d - %v\n", len(path)-1, path)
	fmt.Printf("Total runtime: %v (excludes drawing)\n", runtime)
}

func NewMapNode(x int, y int, z rune) *MapNode {
	return &MapNode{Cord{x, y}, int(z - 'a'), string(z), "", false, 0, -1, false, -1}
}

func (m *HeightMap) FindShortestPath() Path {
	m.Frontiers = make(map[Cord][]Path)
	m.ShortestPath = nil
	m.CacheUtility()
	m.FindPaths()
	return m.ShortestPath
}

func (m *HeightMap) CacheUtility() {
	m.EndPos.Cost = math.MaxInt
	m.LineOfSight = m.Distance(m.StartPos.Position)
	for _, row := range m.Grid {
		for _, n := range row {
			if n == m.EndPos {
				n.Utility = float64(math.MaxInt64)
			} else if n == m.StartPos {
				n.Utility = 0
				n.LosDist = m.LineOfSight
			} else {
				n.Utility = m.CalcUtility(n)
			}
			n.Cost = math.MaxInt
		}
	}
	m.StartPos.Cost = 0
	//fmt.Printf("%s", m.UtilityString())
}

func (m *HeightMap) FindPaths() Path {
	m.AddFrontierPath(Path{m.StartPos})
	for {
		found := m.ExplorePath()
		if found {
			break
		}
	}
	return m.ShortestPath
}

func (m *HeightMap) DumpFrontier() {
	s, err := json.MarshalIndent(m.Frontiers, "", " ")
	if err != nil {
		fmt.Println("Error dumping frontier:", err)
	}
	fmt.Println(string(s))
}

func (m *HeightMap) AddFrontierPath(pS Path) {
	var p Path = make(Path, len(pS))
	copy(p, pS)
	last := p.Last()
	if _, ok := m.Frontiers[last.Position]; !ok {
		m.Frontiers[last.Position] = make([]Path, 0)
	}
	m.Frontiers[last.Position] = append(m.Frontiers[last.Position], p)
	if DumpFrontier {
		m.DumpFrontier()
	}
}

func (m *HeightMap) NextPath() Path {
	if len(m.Frontiers) == 0 {
		return nil
	}
	bestNode := m.StartPos
	for pos, _ := range m.Frontiers {
		node := m.At(pos)
		if node.Utility > bestNode.Utility {
			bestNode = node
		}
	}
	paths := m.Frontiers[bestNode.Position]
	if len(paths) == 0 {
		return nil
	}
	sort.Slice(paths, func(i, j int) bool {
		return len(paths[i]) < len(paths[j])
	})
	if DumpPaths {
		fmt.Println(m.DrawPath(paths[0]))
	}
	delete(m.Frontiers, bestNode.Position)
	return paths[0]
}

func (m *HeightMap) ExplorePath() bool {
	p := m.NextPath()
	if p == nil {
		return true
	}
	this := p.Last()
	cost := len(p) - 1

	if cost < this.Cost {
		this.Cost = cost
	} else if this.Seen {
		// if it doesn't cost less and we see it again, short circuit
		return false
	}
	this.Seen = true

	next := m.Neighbors(this)
	if len(next) == 0 {
		this.Dead = true
		return false
	}
	for _, n := range next {
		if n == m.EndPos {
			m.AddPath(append(p, n))
			return false
		} else if p.Contains(n) >= 0 {
			continue
		} else {
			m.AddFrontierPath(append(p, n))
		}
	}
	return false
}

func (m *HeightMap) Neighbors(curr *MapNode) []*MapNode {
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
		if n.Dead {
			continue
		}
		neighbors = append(neighbors, n)
	}
	return neighbors
}

func (m *HeightMap) AddPath(np Path) {
	p := make(Path, len(np))
	copy(p, np)
	m.FoundPaths = append(m.FoundPaths, p)
	if m.ShortestPath == nil || len(p) < len(m.ShortestPath) {
		m.ShortestPath = p
		//fmt.Printf("Current shortest: %d LoS: %0.3f\n", len(p)-1, m.LineOfSight)
		//fmt.Printf("%s\n", m)
		//fmt.Printf("%s\n", m.PathString())
	}
}

func (p Path) Contains(n *MapNode) int {
	for i, tn := range p {
		if tn == n {
			return i
		}
	}
	return -1
}
func (m *HeightMap) CalcUtility(n *MapNode) float64 {
	n.LosDist = m.Distance(n.Position)
	u := m.LineOfSight - n.LosDist
	//u = u / m.MinElevationDelta(n.Position)
	return u
}
func (m *HeightMap) Distance(c Cord) float64 {
	return m.Distance2d(c) // Much faster...
	//return m.Distance3d(c)
}

/**
 * Utilities/Drawing/Strings/etc
 */
func (m *HeightMap) Distance2d(c Cord) float64 {
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

func (p Path) Last() *MapNode {
	if len(p) == 0 {
		return nil
	}
	return p[len(p)-1]
}

func (c Cord) String() string {
	return fmt.Sprintf("(%d,%d)", c.X, c.Y)
}

func (n *MapNode) String() string {
	return fmt.Sprintf("%s", n.Position)
}
func (p Path) String() string {
	if len(p) == 0 {
		return "{}"
	}
	last := p.Last()
	return fmt.Sprintf("F:%s U:(%.3f) L:(%d)", last, last.Utility, len(p))
}
func (p Path) StringLong() string {
	if len(p) == 0 {
		return "{}"
	}
	last := p.Last()
	vals := ""
	for _, n := range p {
		vals += n.String()
	}
	return fmt.Sprintf("F:%s U:(%.3f) L:(%d) - %s", last, last.Utility, len(p), vals)
}

func (p Path) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.String())
}
func (c Cord) MarshalText() ([]byte, error) {
	return []byte(c.String()), nil
}
func (n *MapNode) MarshalText() ([]byte, error) {
	return []byte(n.String()), nil
}

func CordUp(c Cord) Cord    { return Cord{c.X, c.Y - 1} }
func CordDown(c Cord) Cord  { return Cord{c.X, c.Y + 1} }
func CordLeft(c Cord) Cord  { return Cord{c.X - 1, c.Y} }
func CordRight(c Cord) Cord { return Cord{c.X + 1, c.Y} }

func (m *HeightMap) String() string {
	return m.MarkPath(m.ShortestPath)
}
func (m *HeightMap) PathString() string {
	return m.DrawPath(m.ShortestPath)
}

func (m *HeightMap) DrawPath(sp Path) string {
	s := fmt.Sprintf("%s\n", sp)
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

func (m *HeightMap) MarkPath(sp Path) string {
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
func (m *HeightMap) UtilityString() string {
	s := ""
	for _, row := range m.Grid {
		for _, n := range row {
			s += fmt.Sprintf("%-4.1f ", n.Utility)
		}
		s += "\n"
	}
	return s
}
