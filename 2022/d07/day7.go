package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

var (
	TotalSpace  int = 70000000
	NeededSpace int = 30000000
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	shell := NewReverseShell()
	for {
		str, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		line := strings.TrimSpace(string(str))
		if "$ " == line[0:2] {
			shell.Cmd(line[2:])
		} else {
			shell.StdOut(line)
		}
	}

	shell.LsAll()

	PartOne(shell)
	PartTwo(shell)
}

func PartOne(shell *ReverseShell) {
	sub100k := func(x *FsNode) bool {
		if x.TotalSize() <= 100000 {
			return true
		}
		return false
	}

	var p1_total int = 0
	for _, dir := range shell.FindPaths(sub100k) {
		p1_total += dir.TotalSize()
	}

	fmt.Printf("Part 1 total %d\n", p1_total)
}

func PartTwo(shell *ReverseShell) {
	usage := shell.DiskUsage()
	fmt.Printf("Total disk usage is: %d\n", usage)

	p2_min := usage
	freeSpace := TotalSpace - usage
	sizeToFree := NeededSpace - freeSpace

	deleteCandidates := func(x *FsNode) bool {
		if x.TotalSize() >= sizeToFree {
			return true
		}
		return false
	}

	for _, dir := range shell.FindPaths(deleteCandidates) {
		size := dir.TotalSize()
		if size < p2_min {
			p2_min = size
		}
	}
	fmt.Printf("Part 2 - best dir size %d\n", p2_min)
}
