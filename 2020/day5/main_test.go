package main

import (
    "testing"
    "fmt"
)

type TestCase struct {
    SeatCode string
    Row int
    Col int
    Id  int
}

func Test_SeatDecode(t *testing.T){
    tests := []TestCase{
        TestCase{"FBFBBFFRLR", 44, 5, 357},
        TestCase{"BFFFBBFRRR", 70, 7, 567},
        TestCase{"FFFBBBFRRR", 14, 7, 119},
        TestCase{"BBFFBBFRLL", 102, 4, 820},
    }

    for _, test := range tests {
        row, col, id := SeatDecode( test.SeatCode )
        if row != test.Row || col != test.Col || id != test.Id {
            var errStr string
            errStr = fmt.Sprintf("Given: %v\n", test.SeatCode )
            errStr = errStr + fmt.Sprintf("Expected: Row: %v Col: %v ID: %v\n", test.Row, test.Col, test.Id )
            errStr = errStr + fmt.Sprintf("     Got: Row: %v Col: %v ID: %v\n", row, col, id )
            t.Error(errStr)
        }
    }
}
