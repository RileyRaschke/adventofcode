package main

import (
    "os"
    "bufio"
    "io"
    "fmt"
    "strings"
    "strconv"
)

type BagRule struct  {
    color string
    contains map[*BagRule]int
}
func NewBagRule(color string) (*BagRule) {
    br := &BagRule{color: color, contains: make(map[*BagRule]int)}
    return br
}
func (self *BagRule) AddAllowedContent(br *BagRule, cnt int) {
    fmt.Printf("Adding %v with count: %v\n", br.color, cnt)
    self.contains[br] = cnt
}

func (self *BagRule) CanContain(color string) bool {
    for innerBr, _ := range self.contains {
        if innerBr.color == color || innerBr.CanContain( color ) {
            return true
        }
    }
    return false
}

type RuleSearch struct {
    bagMap map[string]*BagRule
}
func NewRuleSearch() *RuleSearch {
    return &RuleSearch{
        bagMap: make(map[string]*BagRule),
    }
}

func (self *RuleSearch) AddRuleText(str string) {
    fmt.Printf("\n%v\n",str)

    rplcr := strings.NewReplacer("bags","", "bag","", ".","")
    str = strings.TrimSpace(rplcr.Replace(str))

    ruleParts := strings.Split(str, " contain ")
    color := strings.TrimSpace(ruleParts[0])
    contents := strings.TrimSpace(ruleParts[1])

    fmt.Printf("Color: %v\n", color );
    fmt.Printf("Contents:\n")

    // If we haven't seen this color bag,
    // create a rule pointer to hold it's spec
    if _, ok := self.bagMap[color]; !ok {
        self.bagMap[color] = NewBagRule(color)
    }

    // For the bags it can contain, parse how many
    for _, canCon := range strings.Split(contents, ","){
        canCon = strings.TrimSpace(canCon)
        if( canCon == "no other" ){
            fmt.Printf("\t%v\n", "none")
            break
        }
        fmt.Printf("\t%v\n", canCon)
        canConParts := strings.SplitAfterN(canCon, " ", 2)
        canConColor := strings.TrimSpace(canConParts[1])
        canConCount, _ := strconv.Atoi(strings.TrimSpace(canConParts[0]))

        // If we haven't seen this color bag,
        // create a rule pointer to hold it's spec
        if _, ok := self.bagMap[canConColor]; !ok {
            self.bagMap[canConColor] = NewBagRule(canConColor)
        }
        // Add the content rule pointer to the bag color
        self.bagMap[color].AddAllowedContent( self.bagMap[canConColor], canConCount )
    }
}
func (self *RuleSearch) DumpRules() {
    for _, colorRule := range self.bagMap {
        fmt.Printf("%v\n\n", colorRule)
    }
}

func (self *RuleSearch) CountOptions(color string) (res int) {
    fmt.Println()
    for _, colorRule := range self.bagMap {
        if colorRule.CanContain( color ) {
            res++
        }
    }
    return
}

func main() {
    reader := bufio.NewReader(os.Stdin)
    rs := NewRuleSearch()
    for {
        str, _, err := reader.ReadLine()
        if err == io.EOF { break }
        line := strings.TrimSpace(string(str))
        if line == "" { break }
        rs.AddRuleText(line)
    }
    rs.DumpRules() // SUPER SLOW
    fmt.Printf("Part 1 - shiny gold bag combos: %v\n", rs.CountOptions("shiny gold"))
}

func (self *BagRule) String() string {
    return self.Padded("")
}
func (self *BagRule) Padded(pad string) string {
    if len(self.contains) == 0 {
        return self.color + " => ()"
    }
    str := self.color + " => (\n"
    for content, count := range( self.contains ) {
        if pad == "" {
            str += fmt.Sprintf("\t%v %v\n", count, content.Padded("\t\t"))
        } else {
            str += fmt.Sprintf("%v%v %v\n", pad, count, content.Padded(pad + "\t"))
        }
    }
    if pad != "" {
        str += fmt.Sprintf("%v)", pad[:len(pad)-1])
    } else {
        str += fmt.Sprintf(")")
    }
    return str
}
