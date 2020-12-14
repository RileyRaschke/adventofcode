package main

import (
    "os"
    "bufio"
    "io"
    "fmt"
    "strings"
    "strconv"
    "sort"
    //"math"
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

    p1_tracker := p1(adapters)
    fmt.Printf("Part1:\n\t1-jolt diff: %v\n\t3-jolt diff: %v\n\tmultiple: %v\n", p1_tracker[1], p1_tracker[3], p1_tracker[1]*p1_tracker[3])

    for key, val := range p1_tracker {
        fmt.Printf("%v - %v\n", key, val)
    }

    fmt.Print("\nPart2:\n")
    permutations := p2_permutations(adapters)
    fmt.Printf("\n%v permutations\n", permutations)
}

func p1(adapters []int) (map[int]int){
    tracker := make(map[int]int)
    // p1 - brute with mini optimization
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

func p2_permutations(adapters []int) (perms int64){
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
    fmt.Printf("len=%v - %v\n", len(aSets), aSets)
    fmt.Printf("Found %v sets; %v sets are required out of %v total adapters\n", sets, requiredSets, len(adapters))

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
    perms = 1
    keys := []int{}
    for k, _ := range adaptOpts {
        keys = append(keys, k)
    }
    sort.Ints(keys)
    for _, key := range keys {
        setLen := len(adaptOpts[key])
        fmt.Printf("(%d) %2d - %v\n", setLen, key, adaptOpts[key])
        if setLen > 1 {
            perms += int64(setLen)
        }
    }
    //diffMap := make(map[int]int)
    //for _, key := range keys {
        //setLen := len(adaptOpts[key])
    //}
    perms = CalcPerms( 0, keys, adaptOpts )
    return perms
}

func CalcPerms( start int, sortedKeys []int, setMap map[int][]adptSet) (perms int64){
    perms = 1
    for _, key := range sortedKeys {
        setLen := len(setMap[key])
        for _, set := range setMap[key] {
            if next, ok := setMap[set.Second()]; ok {
                if len(next) == 1 {
                    perms += int64(setLen)
                }
                if( len(next) == 3 ){
                    //perms *= int64(len(next)-1)
                    perms *= int64(len(next)-setLen)
                }
                if( setLen == 3 && len(next) == 3 ){
                    perms *= int64(len(next))
                }
                /*
                if setLen > len(next) {
                    //perms *= (factorial(setLen)-factorial(len(next)))
                    perms += (factorial(setLen)-factorial(len(next)))/(factorial(len(next))*int64(setLen-len(next)))
                }
                */
                /*
                if( setLen == 3  && len(next) == 3 ){
                    perms *= 3
                }
                if( setLen == 3  && len(next) == 2 ){
                    perms *= 2
                }
                */
            }
        }
    }
    return
}

func factorial( x int ) (r int64) {
    r=1
    for {
        if x < 2 { break }
        r *= int64(x)
        x--
    }
    return
}

