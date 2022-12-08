package main

import (
	"fmt"
	"strconv"
	"strings"
)

type ReverseShell struct {
	CmdHist []string
	OutHist []string
	fs      *FsNode
	cwd     *FsNode
	echo    bool
}

func NewReverseShell() *ReverseShell {
	rootFs := NewFs()
	return &ReverseShell{[]string{}, []string{}, rootFs, rootFs, false}
}

func (sh *ReverseShell) Cmd(cmd string) {
	sh.CmdHist = append(sh.CmdHist, cmd)
	if sh.echo {
		fmt.Println("$ " + cmd)
	}
	args := strings.Split(cmd, " ")
	switch args[0] {
	case "cd":
		switch args[1] {
		case "..":
			sh.PopCwd()
			return
		case "/":
			sh.cwd = sh.fs
			return
		default:
			sh.cwd = sh.cwd.MkDir(args[1])
			return
		}
	case "ls":
		if len(sh.cwd.Children) == 0 {
			return
		} else {
			sh.cwd.Lsr(0)
		}
		return
	case "pwd":
		fmt.Printf("%s\n", sh.cwd.Path())
		return
	case "du":
		u := sh.cwd.TotalSize()
		fmt.Printf("usage %.2f%% size %d - %s\n", (float32(u)/float32(TotalSpace))*100.0, u, sh.cwd.Path())
		return
	case "df":
		u := sh.fs.TotalSize()
		free := TotalSpace - u
		fmt.Printf("%.2f%% free %d available\n", (float32(free)/float32(TotalSpace))*100.0, free)
		return
	case "mkdir":
		sh.cwd.MkDir(args[1])
		return
	case "echo":
		if args[1] == "on" {
			sh.echo = true
		} else {
			sh.echo = false
		}
		return
	default:
		panic(fmt.Sprintf("command not found: %s", cmd))
	}
}

func (sh *ReverseShell) StdOut(out string) {
	sh.OutHist = append(sh.OutHist, out)
	if sh.echo {
		fmt.Println(out)
	}
	data := strings.Split(out, " ")
	if data[0] == "dir" {
		sh.cwd.MkDir(data[1])
	} else {
		size, _ := strconv.Atoi(data[0])
		sh.cwd.AddFile(data[1], size)
	}
}

func (sh *ReverseShell) PopCwd() {
	if sh.cwd.Parent == nil {
		sh.cwd = sh.fs
	} else {
		sh.cwd = sh.cwd.Parent
	}
}

func (sh *ReverseShell) DiskUsage() int {
	return sh.fs.TotalSize()
}

func (sh *ReverseShell) Cwd() string {
	return sh.cwd.Path()
}

func (sh *ReverseShell) LsAll() {
	sh.fs.Lsr(0)
}

func (sh *ReverseShell) FindPaths(filter func(*FsNode) bool) []*FsNode {
	return sh.fs.FilterDirs(filter)
}
