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

    preamble := int64(25)
    if len(os.Args) > 1 {
        opt, err := strconv.Atoi(os.Args[1])
        if err == nil {
            preamble = int64(opt)
        }
    }

    errIndex, errVal, ok := CheckEncoding(intCode, preamble)
    if !ok {
        fmt.Printf("Part1: Encoding Error at index %v with val: %v\n", errIndex, errVal)
    }

    min, max := ContigMinMax(intCode[0:errIndex], errVal)
    fmt.Printf("Part2: Min/Max sum of contig range summing to %v is %v\n", errVal, min+max)
}

func CheckEncoding( code []int64, preamble int64) (int64, int64, bool) {
    for i := preamble; i <= int64(len(code)); i++ {
        _,_,ok := twoSum(code[i-preamble:i], code[i])
        if !ok {
            return i, code[i], false
        }
    }
    return 0, 0, true
}

func SumSlice(s []int64) (r int64) {
    for _, val := range s {
        r += val
    }
    return
}
// brute force.. tedius junk
func ContigMinMax(code []int64, tgt int64) (min, max int64) {
    var found bool
    var i, j int
    for {
        if i+2 >= len(code) { break }
        j = i+2
        for {
            if j >= len(code) { break }
            // heavy heavy double work here...
            t := SumSlice( code[i:j+1] )
            if tgt == t {
                found = true
                break
            }
            if tgt < t {break}
            if tgt > t { j++ }
        }
        if found {
            break
        } else { i++ }
    }
    r := code[i:j+1]
    sort.Slice(r, func(i,j int) bool { return r[i] < r[j] } )
    return r[0], r[len(r)-1]
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

