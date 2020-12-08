package main

import (
    "os"
    "bufio"
    "io"
    "fmt"
    "strings"
)

func main() {
    reader := bufio.NewReader(os.Stdin)
    ansMap := make(map[rune]bool)
    var sumation int
    for {
        str, _, err := reader.ReadLine()
        if err == io.EOF {
            break;
        }
        line := strings.TrimSpace(string(str))
        if line == "" {
            sumation = sumation + len(ansMap)
            ansMap = make(map[rune]bool)
        } else {
            for _, chr := range line {
                ansMap[chr] = true
            }
        }
    }
    fmt.Printf("Yes answers in groups summed to: %v\n", sumation )
}
