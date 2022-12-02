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
	elfPacks := []int{}
	currPack := 0

	reader := bufio.NewReader(os.Stdin)
	for {
		str, _, err := reader.ReadLine()

		if err == io.EOF {
			break
		}

		line := strings.TrimSpace(string(str))

		if line == "" {
			elfPacks = append(elfPacks, currPack)
			currPack = 0
		} else {
			cals, err := strconv.Atoi(line)
			if err != nil {
				panic(err)
			}
			currPack += cals
		}
	}

	// Don't forget the last guy!
	if currPack != 0 {
		elfPacks = append(elfPacks, currPack)
	}

	sort.Ints(elfPacks)

	fmt.Printf("There are %v elves\n", len(elfPacks)+1)
	fmt.Printf("Biggest pack is %v\n", elfPacks[len(elfPacks)-1])

	topCount := 3
	total := 0
	for _, val := range elfPacks[len(elfPacks)-topCount:] {
		total += val
	}
	fmt.Printf("The top %d elves have a total of %d\n", topCount, total)
}
