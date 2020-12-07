package main

import (
    "os"
    "bufio"
    "io"
    "fmt"
    //"strconv"
    "strings"
)

type FieldChecker struct {
    requiredFields []string
    fieldMap map[string]string
}

func NewFieldChecker( reqFields []string ) *FieldChecker {
    return &FieldChecker{
        requiredFields: reqFields,
        fieldMap: make(map[string]string),
    }
}

func (self *FieldChecker) AddField( x string ) {
    data := strings.Split(x,":")
    self.fieldMap[data[0]] = data[1]
}

func (self *FieldChecker) IsValid() bool {
    for _, field := range self.requiredFields {
        if _, ok := self.fieldMap[field]; !ok {
            return false
        }
    }
    return true
}

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
    var validCount int
    checker := NewFieldChecker( requiredFields )
    reader := bufio.NewReader(os.Stdin)
    for {
        str, _, err := reader.ReadLine()
        if err == io.EOF {
            if checker.IsValid() {
                validCount++
            }
            break;
        }
        line := string(str)
        if line == "" {
            if checker.IsValid() {
                validCount++
            }
            checker = NewFieldChecker( requiredFields )
        } else {
            for _, field := range strings.Split(line," ") {
                checker.AddField( field )
            }
        }
    }

    fmt.Printf("Found %v valid passports\n", validCount )
}
