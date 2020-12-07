package main

import (
    "testing"
    "fmt"
)

type TestCase struct {
    Input []int
    Target int
    Res []int
}

func Test_twoSum( t *testing.T ){
    tests := []TestCase {
        TestCase{ Input: []int{3,2,4}, Target: 6, Res: []int{1,2} },
        TestCase{ Input: []int{2,7,11,15}, Target: 9, Res: []int{0,1} },
        TestCase{ Input: []int{2,9,11,15}, Target: 9, Res: []int{} },
        TestCase{ Input: []int{0,2,9,11,200}, Target: 200, Res: []int{0,4} },
        TestCase{ Input: []int{2,9,11,200}, Target: 200, Res: []int{} },
    }
    for _, c := range tests {
        got := twoSum( c.Input, c.Target )
        if fmt.Sprintf("%v",got) != fmt.Sprintf("%v",c.Res) {
            t.Errorf("Got %v, expected %v", got, c.Res)
        }
    }
}

