package main

import (
    "strings"
    "strconv"
    "regexp"
)

type validator func(string) bool

type FieldValidator struct {
    validators map[string]validator
    fieldMap map[string]string
}

func NewAocPassportValidator() *FieldValidator {
    return NewFieldValidator( map[string]validator {
        "byr": BirthYearValidator,
        "iyr": IssueYearValidator,
        "eyr": ExpirationYearValidator,
        "hgt": HeightValidator,
        "hcl": HairColorValidator,
        "ecl": EyeColorValidator,
        "pid": PassportIdValidator,
    })
}

func NewFieldValidator( validations map[string]validator ) *FieldValidator {
    return &FieldValidator{
        validators: validations,
        fieldMap: make(map[string]string),
    }
}

func (self *FieldValidator) AddField( x string ) {
    data := strings.Split(x,":")
    self.fieldMap[data[0]] = strings.TrimSpace(data[1])
}

func (self *FieldValidator) IsValid() bool {
    for key, validFunc := range self.validators {
        flVal, ok := self.fieldMap[key]
        if !ok {
            return false
        }
        if !validFunc(flVal) {
            return false
        }
    }
    return true
}

func BirthYearValidator(s string) bool {
    year, err := strconv.Atoi(s);
    if err != nil { return false }
    if year >= 1920 && year <= 2002 {
        return true
    }
    return false
}

func IssueYearValidator(s string) bool {
    year, err := strconv.Atoi(s);
    if err != nil { return false }
    if year >= 2010 && year <= 2020 {
        return true
    }
    return false
}

func ExpirationYearValidator(s string) bool {
    year, err := strconv.Atoi(s);
    if err != nil { return false }
    if year >= 2020 && year <= 2030 {
        return true
    }
    return false
}

func HeightValidator(s string) bool {
    hLen := len(s)
    if hLen <= 2 { return false }

    height, err := strconv.Atoi(s[:hLen-2])
    if err != nil { return false }

    if s[hLen-2:] == "in" {
        if height >= 59 && height <= 76 { return true } else { return false }
    }
    if s[hLen-2:] == "cm" {
        if height >= 150 && height <= 193 { return true } else { return false }
    }
    return false
}

func HairColorValidator(s string) bool {
    re := regexp.MustCompile(`^#[a-f0-9]{6}$`)
    return re.Match([]byte(s))
}

func EyeColorValidator(s string) bool {
    opts := []string{"amb","blu","brn","gry","grn","hzl","oth"};
    for _, clr := range opts {
        if clr == s {
            return true
        }
    }
    return false
}

func PassportIdValidator(s string) bool {
    re := regexp.MustCompile(`^[0-9]{9}$`)
    return re.Match([]byte(s))
}

