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
*/
type adptSet struct {
    set []int
}

func (x adptSet) String() string {
    return fmt.Sprintf("(%v,%v)", x.set[0],x.set[1])
}
func (x adptSet) First() int {
    return x.set[0]
}
func (x adptSet) Second() int {
    return x.set[1]
}

func Part2_v1(adapters []int) (perms int64){
    var sets int64
    var requiredSets int64
    aSets := []adptSet{};
    for i:=0; i < len(adapters)-1; i++ {
        var forwardSets = 0
        for j:=i+1; j < i+4 && j<len(adapters); j++ {
            diff := adapters[j] - adapters[i]
            if diff < 4 {
                aSets = append(aSets, adptSet{[]int{adapters[i],adapters[j]}})
                sets++
                forwardSets++
            }
        }
        if forwardSets == 1 {
            requiredSets++
        }
    }
    //fmt.Printf("len=%v - %v\n", len(aSets), aSets)
    //fmt.Printf("Found %v sets; %v sets are required out of %v total adapters\n", sets, requiredSets, len(adapters))

    adaptOpts := make(map[int][]adptSet)
    var tmp int
    for i := 0; i < len(aSets)-1; i++{
        if tmp != aSets[i].First() {
            tmp = aSets[i].First()
        }
        if _,ok := adaptOpts[tmp]; ok {
            adaptOpts[tmp] = append(adaptOpts[tmp], aSets[i])
        } else {
            adaptOpts[tmp] = []adptSet{aSets[i]}
        }
    }
    fmt.Println()
    keys := []int{}
    for k, _ := range adaptOpts {
        keys = append(keys, k)
    }
    sort.Ints(keys)
    /*
    for _, key := range keys {
        setLen := len(adaptOpts[key])
        fmt.Printf("(%d) %2d - %v\n", setLen, key, adaptOpts[key])
        if setLen > 1 {
            perms += int64(setLen)
        }
    }
    */
    perms = CalcPerms( keys, adaptOpts )
    return perms
}

func CalcPerms( sortedKeys []int, setMap map[int][]adptSet) (perms int64){
    perms = 1
    setPerms := int64(0)
    for _, key := range sortedKeys {
        setLen := len(setMap[key])
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

