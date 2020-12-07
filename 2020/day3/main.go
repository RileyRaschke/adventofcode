package main

import (
    "os"
    "bufio"
    "io"
    "fmt"
    //"strconv"
    "strings"
)

func BobSled( sledMap [][]rune, right int, down  int) int {
    var x int
    var y int
    var treeCount int
    for {

        x = x+right
        y = y+down

        if y >= len(sledMap) {
            break
        }

        if x >= len(sledMap[y]) {
            x = x-len(sledMap[y])
        }

        if sledMap[y][x] == '#' {
            treeCount++
        }
    }
    return treeCount
}

type Route struct {
    Right int
    Down int
}

func main() {
    _, err := os.Stdin.Stat()
    if err != nil {
        panic(err)
    }

    reader := bufio.NewReader(os.Stdin)

    var mapHeight int
    sledMap := [][]rune{}
    for {
        line := readLine(reader)
        if line == "" {
            break
        }
        sledMap = append(sledMap, []rune{})
        for _, char := range line {
            sledMap[mapHeight] = append(sledMap[mapHeight], char)
        }
        mapHeight++
    }

    fmt.Printf("%v\n", sledMap );

    fmt.Printf("Part 1 - Hit %v trees\n", BobSled( sledMap, 3, 1) )

    routes := []Route{
        Route{1,1},
        Route{3,1},
        Route{5,1},
        Route{7,1},
        Route{1,2},
    }
    treeMultiple := 1
    for _, route := range routes {
        treeMultiple = treeMultiple * BobSled( sledMap, route.Right, route.Down )
    }
    fmt.Printf("Part 2 - Multiplied Tries: %v\n", treeMultiple )
}

func readLine(reader *bufio.Reader) string {
    str, _, err := reader.ReadLine()
    if err == io.EOF {
        return ""
    }
    return strings.TrimRight(string(str), "\r\n")
}
