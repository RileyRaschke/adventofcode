package main

import (
    "os"
    "bufio"
    "io"
    "fmt"
    "strings"
    "strconv"
)

func main() {
    fmt.Printf("Parsing Encoding\n")
    reader := bufio.NewReader(os.Stdin)
    intCode := []int64{}
    for {
        str, _, err := reader.ReadLine()
        if err == io.EOF { break }
        line := strings.TrimSpace(string(str))
        if line == "" { break }
        x, err := strconv.ParseInt(line, 10, 64)
        intCode = append(intCode, x )
    }
    errIndex, errVal, ok := CheckEncoding(intCode, 25)
    if !ok {
        fmt.Printf("Encoding Error at index %v with val: %v\n", errIndex, errVal)
    }

}

func CheckEncoding( code []int64, preambleLen int64) (int64, int64, bool) {
    for i := preambleLen; i <= int64(len(code)); i++ {
        _,_,ok := twoSum(code[i-preambleLen:i], code[i])
        if !ok {
            return i, code[i], false
        }
    }
    return 0, 0, true
}

func twoSum(nums []int64, target int64) (int64,int64,bool) {
    possibleHits := make(map[int64]int64)
    for i:= int64(0); i < int64(len(nums)); i++ {
        if possibleHits[nums[i]] != 0 {
            //return []int64{possibleHits[nums[i]]-1, i }
            return nums[possibleHits[nums[i]]-1], nums[i], true
        }
        chance := target-nums[i]
        possibleHits[chance] = i+1
    }
    return 0,0,false
}


