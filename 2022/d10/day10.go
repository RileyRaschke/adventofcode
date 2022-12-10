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
	screen           [][]bool
	sHeight, sWidth  int = 6, 40
)

func main() {
	rx = 1
	screen = NewScreen()
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
	fmt.Printf("Part2:\n%s\n", DrawScreen())
}

func runCycle() {
	cycle++
	if isSignal(cycle) {
		total += cycle * rx
	}
	processPixel(cycle)
}

func isSignal(c int) bool {
	if (c+20)%40 == 0 {
		return true
	}
	return false
}

func processPixel(c int) {
	row := (c - 1) / 40
	col := (c - 1) % 40
	if col == rx || col == rx-1 || col == rx+1 {
		screen[row][col] = true
	}
}

func NewScreen() [][]bool {
	s := make([][]bool, sHeight)
	for i := range s {
		s[i] = make([]bool, sWidth)
	}
	return s
}

func DrawScreen() string {
	s := ""
	for _, row := range screen {
		for _, px := range row {
			if px {
				s += "#"
			} else {
				s += "."
			}
		}
		s += "\n"
	}
	return s
}
