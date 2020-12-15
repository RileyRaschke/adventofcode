package main

import (
    "os"
    "bufio"
    "io"
    "fmt"
    "strings"
    "strconv"
    "sort"
    "gonum.org/v1/gonum/graph"
    "gonum.org/v1/gonum/graph/simple"
    "time"
    //"gonum.org/v1/gonum/graph/path"
)

func main() {
    reader := bufio.NewReader(os.Stdin)
    adapters := []int{}
    for {
        str, _, err := reader.ReadLine()
        if err == io.EOF { break }
        line := strings.TrimSpace(string(str))
        if line == "" { break }
        x, err := strconv.Atoi(line)
        adapters = append(adapters, x )
    }
    // Add the outlet
    adapters = append(adapters,0)
    sort.Ints(adapters)
    // Add the device
    adapters = append(adapters, adapters[len(adapters)-1]+3)

    p1_tracker := Part1(adapters)
    fmt.Printf("Part1:\n\t1-jolt diff: %v\n\t3-jolt diff: %v\n\tmultiple: %v\n", p1_tracker[1], p1_tracker[3], p1_tracker[1]*p1_tracker[3])


    fmt.Print("\nPart2:\n")

    start := time.Now()
    permutations := Part2_v1(adapters)
    rt := time.Since(start)

    fmt.Printf("\n\tv1 Set   %15d permutations in %v\n", permutations, rt)

    start = time.Now()
    permutations_graph := Part2_v2(adapters)
    rt = time.Since(start)

    fmt.Printf("\n\tv2 Graph %15d permutations in %v\n\n", permutations_graph, rt)
}

func Part1(adapters []int) (map[int]int) {
    tracker := make(map[int]int)
    for i:=0; i < len(adapters); i++ {
        for j:=i+1; j < i+4 && j<len(adapters); j++ {
            diff := adapters[j] - adapters[i]
            if( diff < 4 ){
                if _, ok := tracker[diff]; ok {
                    tracker[diff]++
                } else {
                    tracker[diff] = 1
                }
                if( diff == 1 || diff == 3 ){
                    break
                }
            }
        }
    }
    return tracker
}

func Part2_v2(adapters []int) (int64) {
    g := simple.NewDirectedGraph()
    BuildGraph(g, adapters)
    return EvalPossibleConnections( g )
}

func BuildGraph(g *simple.DirectedGraph, adapters []int) {
    for i:=0; i < len(adapters)-1; i++ {
        for j:=i+1; j < i+4 && j<len(adapters); j++ {
            from := int64(adapters[i])
            to   := int64(adapters[j])
            diff := to-from
            if( diff < 4 ){
                frNode := g.Node(from)
                toNode := g.Node(to)
                if frNode == nil {
                    frNode = simple.Node(from)
                    g.AddNode(frNode)
                }
                if toNode == nil {
                    toNode = simple.Node(to)
                    g.AddNode(toNode)
                }
                g.SetEdge( g.NewEdge(frNode, toNode) )
            }
        }
    }
}

func EvalPossibleConnections( g *simple.DirectedGraph ) int64 {

    sortedNodeIds := SortedNodeIds( g )
    perms   := int64(1)
    edgeSum := int64(0)

    for _, nodeId := range sortedNodeIds {

        childNodes := graph.NodesOf(g.From(nodeId))
        edgedgeSum := len(childNodes)

        if edgedgeSum != 1 {
            edgeSum += int64(edgedgeSum)
        } else if edgeSum > 1 {

            if edgeSum == 2 {
                perms *= edgeSum
            } else {
                perms *= (edgeSum-1)
            }
            edgeSum = 0

        }
    }
    return perms
}

func SortedNodeIds( g graph.Graph ) []int64 {
    sortedNodeIds := []int64{}
    for _, node := range graph.NodesOf( g.Nodes() ) {
        sortedNodeIds = append(sortedNodeIds, node.ID())
    }
    sort.Slice(sortedNodeIds, func(i,j int) bool { return sortedNodeIds[i] < sortedNodeIds[j] } )
    return sortedNodeIds
}

/**
* v1 stuff.. I was getting super close on my own..
* v1 was nasty.. this is pretty clean after learning solution with graphs and a hint
*/
type adapterSet struct {
    from int
    to int
}
func (x adapterSet) String() string {
    return fmt.Sprintf("(%v,%v)", x.from, x.to)
}

func Part2_v1(adapters []int) (perms int64){

    setMap := make(map[int][]adapterSet)

    for i:=0; i < len(adapters)-1; i++ {
        for j:=i+1; j < i+4 && j<len(adapters); j++ {
            from := adapters[i]
            to   := adapters[j]
            diff := to-from
            if( diff < 4 ){
                if _, ok := setMap[from]; !ok {
                    setMap[from] = []adapterSet{ adapterSet{from, to} }
                } else {
                    setMap[from] = append( setMap[from], adapterSet{from, to} )
                }
            }
        }
    }
    sortedKeys := []int{}
    for k, _ := range setMap {
        sortedKeys = append(sortedKeys, k)
    }
    sort.Ints(sortedKeys)

    perms = CalcPerms( sortedKeys, setMap )
    return perms
}

func CalcPerms( sortedKeys []int, setMap map[int][]adapterSet) (perms int64){
    perms = 1
    setPerms := int64(0)
    for _, key := range sortedKeys {
        setLen := len(setMap[key])
        fmt.Printf("(%d) %2d - %v\n", setLen, key, setMap[key])
        if setLen != 1 {
            setPerms += int64(setLen)
        } else {
            if setPerms > 1 {
                if setPerms == 2 {
                    perms *= 2
                } else {
                    perms *= (setPerms-1)
                }
                setPerms = 0
            }
        }
    }
    return
}

