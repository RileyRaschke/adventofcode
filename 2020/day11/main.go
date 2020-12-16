package main

import (
    "bufio"
    "fmt"
    "io"
    "os"
    //"time"
)

func main() {
    reader := bufio.NewReader(os.Stdin)
    seatMap := [][]rune{{}}
    var i int
    for {
        char, _, err := reader.ReadRune()
        if err == io.EOF {
            break
        }
        if char == '\n' || char == '\r' {
            i++
            seatMap = append(seatMap, []rune{})
            continue
        }
        seatMap[i] = append(seatMap[i], char)
    }
    fmt.Printf("%v\b", seatMap)
}
