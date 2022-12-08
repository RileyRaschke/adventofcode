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
	var p1_total int = 0
	sumQualifiers := func(x *FsNode) bool {
		size := x.TotalSize()
		if size <= 100000 {
			p1_total += size
			return true
		}
		return false
	}

	shell.FindPaths(sumQualifiers)

	fmt.Printf("Part 1 total %d\n", p1_total)
}

func PartTwo(shell *ReverseShell) {
	usage := shell.DiskUsage()
	fmt.Printf("Total disk usage is: %d\n", usage)
	freeSpace := TotalSpace - usage
	sizeToFree := NeededSpace - freeSpace

	var p2_min int = usage
	var delNode *FsNode

	deletionCandidates := func(x *FsNode) bool {
		size := x.TotalSize()
		if size >= sizeToFree {
			if size < p2_min {
				p2_min = size
				delNode = x
			}
			return true
		}
		return false
	}

	possiblePaths := len(shell.FindPaths(deletionCandidates))
	fmt.Printf("Part 2 - Delete %s with size %d (out of %d possible paths)\n", delNode.Path(), p2_min, possiblePaths)
}
