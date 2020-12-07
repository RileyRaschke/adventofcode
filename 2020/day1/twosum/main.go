package main

import (
    "os"
    "bufio"
    "io"
    "fmt"
    "strconv"
    "strings"
    "sort"
)

func main() {
    _, err := os.Stdin.Stat()
    if err != nil {
        panic(err)
    }

    reader := bufio.NewReader(os.Stdin)
    input := []int{}

    for {
        line := readLine(reader)
        if line == "" {
            break
        }
        num, err := strconv.Atoi(line)
        if err != nil {
            panic(fmt.Sprintf("Non integer input received: %v", line))
        }
        input = append(input, num)
    }
    fmt.Printf("%v\n", input)
    sort.Ints(input)

    indexHits := twoSum(input, 2020)

    fmt.Printf("%v\n", input[indexHits[0]]*input[indexHits[1]] )
}

func readLine(reader *bufio.Reader) string {
    str, _, err := reader.ReadLine()
    if err == io.EOF {
        return ""
    }
    return strings.TrimRight(string(str), "\r\n")
}

func twoSum(nums []int, target int) []int {
    possibleHits := make(map[int]int)
    for i:= 0; i < len(nums); i++ {
        if possibleHits[nums[i]] != 0 {
            return []int{possibleHits[nums[i]]-1, i }
        }
        chance := target-nums[i]
        possibleHits[chance] = i+1
    }
    return []int{}
}

