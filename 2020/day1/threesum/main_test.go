package main

import (
    "testing"
    "fmt"
    "time"
    "math/rand"
    "sort"
)

type TestCase struct {
    Input []int
    Res [][]int
}

func Test_3Sum( t *testing.T ){

    testCases := []TestCase{
        BuildCase(4000,500),
        BuildCase(45,15),
        BuildCase(10,4),
        BuildCase(6,3),
        BuildCase(6,3),
        BuildCase(6,3),
        BuildCase(6,3),
        TestCase{ Input: []int{1,-1,-1,0}, Res: [][]int{{-1,0,1}} },
        TestCase{ Input: []int{3,0,-2,-1,1,2},  Res: [][]int{{-2,-1,3},{-2,0,2},{-1,0,1 }} },
        TestCase{ Input: []int{1,1,1,1,0,-1,0}, Res: [][]int{{-1,0,1}} },
        TestCase{ Input: []int{-1,0,1,2,-1,-4}, Res: [][]int{{-1,-1,2},{-1,0,1}} },
    }
    fmt.Println();
    for _, test := range testCases {
        fmt.Printf("   Input: %v\n", test.Input)
        fmt.Printf("     Got: %v\n", threeSum( test.Input ))
        fmt.Printf("Expected: %v\n\n", test.Res )
    }
}

func BuildCase(length, solutions int) TestCase {
    hits := make(map[string][]int)
    var hitSets [][]int
    var res []int
    rand.Seed(time.Now().UnixNano())

    for n := 0; n < solutions; {
        i, j, k := 0,0,0
        if n%2 == 0 {
            i, j, k = rand.Intn(length), rand.Intn(length)*-1, rand.Intn(length)
        } else {
            i, j, k = rand.Intn(length), rand.Intn(length)*-1, rand.Intn(length)*-1
        }
        sum := i + j + k
        if sum == 0 {
            hitSet := []int{i,j,k}
            sort.Ints(hitSet)
            strRep := fmt.Sprintf("%v",hitSet)
            if hits[strRep] == nil {
                n++
                hits[strRep] = hitSet
            }
            res = append( res, hitSet...)
        }
    }
    for _, v := range hits {
        hitSets = append( hitSets, v )
    }
    return TestCase{ Input: res, Res: hitSets }
}
