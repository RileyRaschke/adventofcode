package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func charValue(x rune) int {
	i := int(x - '0')
	if i > 42 {
		return i - 48
	}
	return i - 16 + 26
}

func findBadge(group []string) int {
	fmt.Printf("%v\n", group)
	var smallestLength int = 0
	var smallestBagIndex int = 0
	for idx, val := range group {
		if idx == 0 {
			smallestLength = len(val)
		} else {
			if len(val) < smallestLength {
				smallestLength = len(val)
				smallestBagIndex = idx
			}
		}
	}
	fmt.Printf("Smallest bag index: %d\n", smallestBagIndex)
	for _, c := range []rune(group[smallestBagIndex]) {
		var foundAll = true
		for idx, bag := range group {
			if idx == smallestBagIndex {
				continue
			} else {
				if strings.IndexRune(bag, c) == -1 {
					foundAll = false
					break
				}
			}
		}
		if foundAll {
			fmt.Printf("Found %v\n", string(c))
			return charValue(c)
		}
	}
	return 0
}

var (
	groupSize int = 3
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	var total int = 0
	var badgeTotal int = 0

	groupBags := []string{}

	for {
		str, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		line := strings.TrimSpace(string(str))
		// part 1
		c1 := line[:len(line)/2]
		c2 := line[len(line)/2:]
		for _, c := range []rune(c1) {
			if strings.IndexRune(c2, c) >= 0 {
				total += charValue(c)
				break
			}
		}

		// part 2
		groupBags = append(groupBags, line)

		if len(groupBags) == groupSize {
			badgeTotal += findBadge(groupBags)
			groupBags = []string{}
		}
	}

	fmt.Printf("Part1 Total is: %d\n", total)
	fmt.Printf("Part2 Total is: %d\n", badgeTotal)
}
