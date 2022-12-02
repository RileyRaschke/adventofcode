package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	i := 0
	currPack := 0
	largestIndex := -1
	elfPacks := []int{}

	for {
		str, _, err := reader.ReadLine()

		if err == io.EOF {
			// Don't forget the last guy!
			if elfPacks[largestIndex] < currPack {
				largestIndex = i
			}
			elfPacks = append(elfPacks, currPack)
			break
		}

		line := strings.TrimSpace(string(str))

		if line == "" {
			if largestIndex == -1 {
				largestIndex = i
			} else if elfPacks[largestIndex] < currPack {
				largestIndex = i
			}
			elfPacks = append(elfPacks, currPack)
			currPack = 0
			i++
		} else {
			cals, err := strconv.Atoi(line)
			if err != nil {
				panic(err)
			}
			currPack += cals
		}
	}

	fmt.Printf("There are %v elves\n", len(elfPacks)+1)
	fmt.Printf("Biggest pack is %v at index %v\n", elfPacks[largestIndex], largestIndex)

	sort.Ints(elfPacks)
	topCount := 3
	total := 0
	for _, val := range elfPacks[len(elfPacks)-topCount:] {
		fmt.Printf("%v\n", val)
		total += val
	}
	fmt.Printf("The top %d elves have a total of %d\n", topCount, total)
}
