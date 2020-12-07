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

    hits := threeSumTarget(input, 2020)
    fmt.Printf("%v\n", hits )

    fmt.Printf("%v\n", hits[0][0] * hits[0][1] * hits[0][2] )
}

func readLine(reader *bufio.Reader) string {
    str, _, err := reader.ReadLine()
    if err == io.EOF {
        return ""
    }
    return strings.TrimRight(string(str), "\r\n")
}

func threeSum(nums []int) [][]int {
	return threeSumTarget(nums, 0)
}

func threeSumTarget(nums []int, target int) [][]int {

    res := make([][]int, 0)
    hits := make(map[string]int)
    length := len(nums)
    sort.Ints(nums)

    for i := 0; i < length; i++ {
        if i > 0 && nums[i] == nums[i-1] { continue }

        k := length-1
        for j := i+1; j < k; {
            sum := nums[i] + nums[j] + nums[k]

            if sum == target {
                hitSet := []int{nums[i],nums[j],nums[k]}
                strRep := fmt.Sprintf("%v",hitSet)

                if hits[strRep] == 0 {
                    hits[strRep] = 1
                    res = append(res, hitSet)
                } else {
                    hits[strRep]++
                }
                j++
                k--
                for { if j > k || nums[k] != nums[k+1] { break } else {k--}}
                for { if j > k || nums[j] != nums[j-1] { break } else {j++}}

            } else {
                if sum > target {
                    k--
                } else {
                    j++
                }
            }

        }
    }

    return res
}

