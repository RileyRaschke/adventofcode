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

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%v - %d\n", string('a'), charValue('a'))
	fmt.Printf("%v - %d\n", string('z'), charValue('z'))
	fmt.Printf("%v - %d\n", string('A'), charValue('A'))
	fmt.Printf("%v - %d\n", string('Z'), charValue('Z'))

	var total int = 0
	for {
		str, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		line := strings.TrimSpace(string(str))
		c1 := line[:len(line)/2]
		c2 := line[len(line)/2:]
		for _, c := range []rune(c1) {
			if strings.IndexRune(c2, c) >= 0 {
				total += charValue(c)
				break
			}
		}
	}

	fmt.Printf("Total is: %d\n", total)
}
