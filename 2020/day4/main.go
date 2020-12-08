package main

import (
    "os"
    "bufio"
    "io"
    "fmt"
    "strings"
)

func main() {
    requiredFields := []string{
        "byr",
        "iyr",
        "eyr",
        "hgt",
        "hcl",
        "ecl",
        "pid",
    }

    checker := NewFieldChecker( requiredFields )
    var fieldPresentCount int

    aocValidator := NewAocPassportValidator()
    var validPassportCount int

    reader := bufio.NewReader(os.Stdin)
    for {
        str, _, err := reader.ReadLine()
        if err == io.EOF {
            if checker.IsValid() {
                fieldPresentCount++
            }
            if aocValidator.IsValid() {
                validPassportCount++
            }
            break;
        }
        line := string(str)
        if line == "" {
            if checker.IsValid() {
                fieldPresentCount++
            }
            if aocValidator.IsValid() {
                validPassportCount++
            }
            checker = NewFieldChecker( requiredFields )
            aocValidator = NewAocPassportValidator()
        } else {
            for _, field := range strings.Split(line," ") {
                checker.AddField( field )
                aocValidator.AddField( field )
            }
        }
    }

    fmt.Printf("Found %v passports with all field present\n", fieldPresentCount )
    fmt.Printf("Found %v valid passports\n", validPassportCount )
}
