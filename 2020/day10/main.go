package main

import (
    "os"
    "bufio"
    "io"
    "fmt"
    "strings"
    "strconv"
    "sort"
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
}

func p1(adapters []int) (map[int]int){
    tracker := make(map[int]int)
    // p1 - brute with mini optimization
    for i:=0; i < len(adapters); i++ {
        for j:=i+1; j < i+4 && j<len(adapters); j++ {
            diff := adapters[j] - adapters[i]
            if( diff == 1 || diff == 3 ){
                if _, ok := tracker[diff]; ok {
                    tracker[diff]++
                } else {
                    tracker[diff] = 1
                }
                break
            }
        }
    }
    return tracker
}

