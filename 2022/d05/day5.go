package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	stacks := []string{}
	instructions := []string{}
	readingInstructions := false
	for {
		str, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		line := string(str)

		if line == "" {
			readingInstructions = true
			continue
		}

		if readingInstructions {
			instructions = append(instructions, line)
		} else {
			stacks = append(stacks, line)
		}
	}

	fmt.Println("*********** Part 1 ************")
	stackTracker := NewStackTracker(stacks)
	fmt.Printf("%s", stackTracker)

	p1_start := time.Now()
	for _, move := range instructions {
		stackTracker.RunInstruction(move, "9000")
	}
	p1_time := time.Since(p1_start)

	p1TopString := stackTracker.TopString()
	fmt.Printf("%s", stackTracker)

	fmt.Println("*********** Part 2 ************")
	stackTracker = NewStackTracker(stacks)
	fmt.Printf("%s", stackTracker)

	p2_start := time.Now()
	for _, move := range instructions {
		stackTracker.RunInstruction(move, "9001")
	}
	p2_time := time.Since(p2_start)

	p2TopString := stackTracker.TopString()
	fmt.Printf("%s", stackTracker)

	// Results
	fmt.Printf("Part one top string: %s in %s\n", p1TopString, p1_time)
	fmt.Printf("Part two top string: %s in %s\n", p2TopString, p2_time)
}

type StackTracker struct {
	Stacks [][]string
}

func NewStackTracker(rawStacks []string) *StackTracker {
	st := &StackTracker{[][]string{}}
	maxHeight := len(rawStacks) - 1
	var x, y int = 0, 0
	for i := maxHeight - 1; i >= 0; i-- {
		x = 0
		for c := 0; c+3 <= len(rawStacks[i]); c += 4 {
			if i == maxHeight-1 {
				st.Stacks = append(st.Stacks, make([]string, maxHeight))
			}
			st.Stacks[x][y] = rawStacks[i][c : c+3]
			x++
		}
		y++
	}
	st.TrimColumns()
	return st
}

func (s *StackTracker) RunInstruction(move string, model string) {
	var count, src, dst int
	fmt.Sscanf(move, "move %d from %d to %d", &count, &src, &dst)
	src--
	dst--
	if model == "9000" {
		for i := 0; i < count; i++ {
			srcCol := s.Stacks[src]
			s.Stacks[dst] = append(s.Stacks[dst], srcCol[len(srcCol)-1])
			s.Stacks[src] = srcCol[0 : len(srcCol)-1]
		}
	} else {
		srcCol := s.Stacks[src]
		s.Stacks[dst] = append(s.Stacks[dst], srcCol[len(srcCol)-count:len(srcCol)]...)
		s.Stacks[src] = srcCol[0 : len(srcCol)-count]
	}

}

func (s *StackTracker) String() string {
	var res string = ""
	for y := s.MaxHeight() - 1; y >= 0; y-- {
		row := []string{}
		for x := 0; x < len(s.Stacks); x++ {
			if y < len(s.Stacks[x]) {
				row = append(row, fmt.Sprintf("%-3s", s.Stacks[x][y]))
			} else {
				row = append(row, fmt.Sprintf("%-3s", " "))
			}
		}
		res += strings.Join(row, " ") + "\n"
	}
	for idx, _ := range s.Stacks {
		res += fmt.Sprintf(" %-2d ", idx+1)
	}
	res += "\n"
	return res
}

func (s *StackTracker) MaxHeight() int {
	var maxHeight int = 0
	for _, col := range s.Stacks {
		height := len(col)
		if height > maxHeight {
			maxHeight = height
		}
	}
	return maxHeight
}

func (s *StackTracker) TopString() string {
	tops := ""
	for _, col := range s.Stacks {
		tops += col[len(col)-1]
	}
	tops = strings.ReplaceAll(tops, "[", "")
	tops = strings.ReplaceAll(tops, "]", "")
	return tops
}

func (s *StackTracker) TrimColumns() {
	for x, col := range s.Stacks {
		for y := len(col) - 1; y >= 0; y-- {
			if strings.TrimSpace(col[y]) == "" {
				s.Stacks[x] = s.Stacks[x][0 : len(s.Stacks[x])-1]
			} else {
				break
			}
		}
	}
}
