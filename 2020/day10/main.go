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
    tracker := make(map[int]int)

    // brute with mini optimization
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
    // Add the computer
    tracker[3]++
    fmt.Printf("Part1:\n\t1-jolt diff: %v\n\t3-jolt diff: %v\n\tmultiple: %v\n", tracker[1], tracker[3], tracker[1]*tracker[3])
}

