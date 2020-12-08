package main

import "strings"

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
