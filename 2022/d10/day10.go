package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

var (
	total, cycle, rx int
)

func main() {
	rx = 1
	reader := bufio.NewReader(os.Stdin)
	for {
		str, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		line := strings.TrimSpace(string(str))
		var op string
		var v int
		n, _ := fmt.Sscanf(line, "%s %d", &op, &v)
		if n == 1 {
			runCycle()
		} else {
			runCycle()
			runCycle()
			rx += v
		}
	}

	fmt.Printf("Part1: %d\n", total)
}

func runCycle() {
	cycle++
	if isSignal(cycle) {
		total += cycle * rx
	}
}

func isSignal(c int) bool {
	if (c+20)%40 == 0 {
		return true
	}
	return false
}
