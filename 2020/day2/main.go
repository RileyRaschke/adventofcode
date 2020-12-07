package main

import (
    "os"
    "bufio"
    "io"
    "fmt"
    "strconv"
    "strings"
)

func validPassword( min int, max int, char rune, pwd string) bool {
    if len(pwd) == 0 {
        return false
    }
    var occuranceCount int
    for _, letter := range pwd {
        if letter == char {
            occuranceCount++
            if occuranceCount > max {
                return false
            }
        }
    }
    if occuranceCount < min {
        return false
    }
    return true
}

func validPasswordPositions( min int, max int, char rune, pwd string) bool {
    if len(pwd) == 0 {
        return false
    }
    a := rune(pwd[min-1])
    b := rune(pwd[max-1])
    if( (a==b) || (a!=b && a!=char && b!= char) ){
        return false
    }
    return true
}

func parseInput( inStr string ) (int, int, rune, string){
    strParams := strings.Fields(inStr)
    if len(strParams) != 3 {
        panic(fmt.Sprintf("Invalid input: %v", inStr))
    }

    counts := strings.Split( strParams[0], "-")
    min, err := strconv.Atoi(counts[0])
    if err != nil {
        panic(fmt.Sprintf("Invalid min input: %v", inStr))
    }
    max, err := strconv.Atoi(counts[1])
    if err != nil {
        panic(fmt.Sprintf("Invalid max input: %v", inStr))
    }
    //replacer := strings.NewReplacer(":","")
    //char := rune(replacer.Replace(strParams[1])[0])
    char := rune(strParams[1][0])

    return min, max, char, strParams[2]
}


func main() {
    fmt.Printf("Testing passwords\n")
    _, err := os.Stdin.Stat()
    if err != nil {
        panic(err)
    }

    reader := bufio.NewReader(os.Stdin)
    var validCount int
    var validPositionCount int

    for {
        line := readLine(reader)
        if line == "" {
            break
        }
        if validPassword( parseInput( line ) ) {
            validCount++
        }
        if validPasswordPositions( parseInput( line ) ) {
            validPositionCount++
        }
    }
    fmt.Printf("Found %v valid passwords\n", validCount )
    fmt.Printf("Found %v valid password positions\n", validPositionCount )
}

func readLine(reader *bufio.Reader) string {
    str, _, err := reader.ReadLine()
    if err == io.EOF {
        return ""
    }
    return strings.TrimRight(string(str), "\r\n")
}

