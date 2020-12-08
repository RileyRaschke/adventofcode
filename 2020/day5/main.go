package main

import (
    "os"
    "bufio"
    "io"
    "strings"
    "fmt"
    "math"
    "sort"
)

func SeatDecode(s string) (row, col, id int){
    rowCode := s[:7]
    colCode := s[7:]
    row = BSP(rowCode,'F','B')
    col = BSP(colCode,'L','R')
    id = (row*8)+col
    return
}

func BSP(str string, down rune, up rune) int {
    high := IntPow(2,len(str))-1
    low  := 0
    for idx, char := range str {
        if char == down {
            high = high-((high-low)/2)-1
            if idx+1 == len(str){ return high }
        } else {
            low = low+((high-low)/2)+1
            if idx+1 == len(str){ return low }
        }
    }
    return 0
}

func IntPow(x,y int) int {
    return int(math.Pow(float64(x),float64(y)))
}

func main() {
    fmt.Println("Finding highest seat number...")
    reader := bufio.NewReader(os.Stdin)
    seatsTaken := []int{}
    var maxId int
    for {
        str, _, err := reader.ReadLine()
        if err == io.EOF {
            break;
        }
        line := strings.TrimSpace(string(str))
        _,_,id := SeatDecode(line)
        seatsTaken = append(seatsTaken, id)
        if id > maxId {
            maxId = id
        }
    }
    sort.Ints(seatsTaken)
    var mySeat int
    for i := 1; i < len(seatsTaken); i++ {
        if seatsTaken[i-1] != seatsTaken[i]-1{
            mySeat = seatsTaken[i]-1
        }
    }
    fmt.Printf("SeatList %v\n", seatsTaken)
    fmt.Printf("mySeat: %v\n", mySeat)
    fmt.Printf("Highest seat number: %v\n", maxId)
}


