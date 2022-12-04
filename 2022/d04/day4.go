package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

var (
	groupSize int = 3
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	var total int = 0
	var p2total int = 0

	for {
		str, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		line := strings.TrimSpace(string(str))
		z := strings.Split(line, ",")
		a := strings.Split(z[0], "-")
		b := strings.Split(z[1], "-")

		a0, err := strconv.Atoi(a[0])
		a1, err := strconv.Atoi(a[1])
		b0, err := strconv.Atoi(b[0])
		b1, err := strconv.Atoi(b[1])

		if a0 >= b0 && a1 <= b1 {
			total += 1
		} else if a0 <= b0 && a1 >= b1 {
			total += 1
		}

		if a0 >= b0 && (a0 <= b1 && a1 <= b1) {
			fmt.Println("a")
			p2total += 1
		} else if a0 <= b1 && (a0 >= b0 || a1 >= b0) {
			p2total += 1
			fmt.Println("b")
		}

	}

	fmt.Printf("Part1 Total is: %d\n", total)
	fmt.Printf("Part2 Total is: %d\n", p2total)
}
