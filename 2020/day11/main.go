package main

import (
    "bufio"
    "fmt"
    "io"
    "os"
    "time"
)

func main() {
    reader := bufio.NewReader(os.Stdin)
    seatMap := [][]rune{{}}
    var i int
    var priorI int
    for {
        char, _, err := reader.ReadRune()
        if err == io.EOF {
            break
        }
        if char == '\n' || char == '\r' {
            i++
            continue
        }
        if i != priorI {
            seatMap = append(seatMap, []rune{})
            priorI = i
        }
        seatMap[i] = append(seatMap[i], char)
    }
    //fmt.Printf("%v\b", seatMap)

    start := time.Now()
    p1 := Part1( seatMap )
    rt := time.Since(start)
    fmt.Printf("Part1: %v seats used. (%v)\n", p1, rt)
}

type SeatNumber struct {
    i int
    j int
}

func Part1( seats [][]rune ) int {
    rSeats := make([][]rune, len(seats))
    fmt.Printf("%v\n", seats)
    fmt.Printf("\n%v\n", rSeats)
    for idx, row := range seats {
        rSeats[idx] = make([]rune, len(row))
        copy(rSeats[idx], row)
    }
    fmt.Printf("\n%v\n", rSeats)
    for i:=0; i < len(seats); i++ {
        for j:=0; j < len(seats[i]); j++ {
            if len(OccupiedSeats( seats, SeatNumber{i,j} )) == 0 {
                rSeats[i][j] = '#'
            }
        }
    }
    return 0
}

func OccupiedSeats( seats [][]rune, x SeatNumber ) []SeatNumber {
    occupiedSeats := []SeatNumber{}

    return occupiedSeats
}

func AllOccupiedSeats( seats [][]rune ) []SeatNumber {
    occupiedSeats := []SeatNumber{}
    return occupiedSeats
}

