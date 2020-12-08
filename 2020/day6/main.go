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
    ansMap := make(map[rune]int)
    var yesSumation int
    var groupSumation int
    var memberCount int
    for {
        str, _, err := reader.ReadLine()
        if err == io.EOF {
            break;
        }
        line := strings.TrimSpace(string(str))
        if line == "" {
            yesSumation = yesSumation + len(ansMap)
            for _, val := range ansMap {
                if val == memberCount {
                    groupSumation++
                }
            }
            ansMap = make(map[rune]int)
            memberCount = 0
        } else {
            memberCount++
            for _, chr := range line {
                if val, ok := ansMap[chr]; ok {
                    ansMap[chr] = val + 1
                } else {
                    ansMap[chr] = 1
                }
            }
        }
    }
    fmt.Printf("Yes answers in groups summed to: %v\n", yesSumation )
    fmt.Printf("Yes unions amongst groups summed to: %v\n", groupSumation )
}
