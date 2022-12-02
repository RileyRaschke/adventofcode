
package main

import(
    "testing"
    "fmt"
)

type TestCase struct {
    min int
    max int
    char rune
    pwd string
    valid bool
}

func Test_validPassword2( t *testing.T ){
    fmt.Print("\nTesting occurance based passwords\n")
    tests := []TestCase{
        TestCase{ 2, 4, 'a', "bbaba", true },
        TestCase{ 0, 1, 'a', "bbaba", false },
        TestCase{ 2, 4, 'a', "bbabaaaaaaa", false },
        TestCase{ 2, 4, 'a', "bbaba", true },
    }

    for _, tc := range tests {
        if tc.valid != validPassword( tc.min, tc.max, tc.char, tc.pwd ) {
            t.Errorf("%v-%v %v: %v should eval: %v, got: %v", tc.min, tc.max, string(tc.char), tc.pwd, tc.valid, !tc.valid )
        } else {
            fmt.Printf("%v-%v %v: %v  Valid: %v\n", tc.min, tc.max, string(tc.char), tc.pwd, tc.valid)
        }
    }
}

func Test_validPasswordPositions3( t *testing.T ){
    fmt.Print("\nTesting position based passwords\n")
    tests := []TestCase{
        TestCase{ 1, 3, 'a', "abcde", true },
        TestCase{ 1, 3, 'b', "cdefg", false },
        TestCase{ 2, 9, 'c', "cccccccccc", false },
    }

    for _, tc := range tests {
        if tc.valid != validPasswordPositions( tc.min, tc.max, tc.char, tc.pwd ) {
            t.Errorf("%v-%v %v: %v should eval: %v, got: %v", tc.min, tc.max, string(tc.char), tc.pwd, tc.valid, !tc.valid )
        } else {
            fmt.Printf("%v-%v %v: %v  Valid: %v\n", tc.min, tc.max, string(tc.char), tc.pwd, tc.valid)
        }
    }
}

