package main

import (
	"fmt"
	"sort"
)

type FsNode struct {
	Name     string
	Parent   *FsNode
	Children map[string]*FsNode
	Size     int
}

func NewFs() *FsNode {
	return &FsNode{"/", nil, map[string]*FsNode{}, 0}
}

func NewDirectory(parent *FsNode, name string) *FsNode {
	return &FsNode{name, parent, map[string]*FsNode{}, 0}
}

func (fs *FsNode) String() string {
	if fs.Size == 0 {
		return fmt.Sprintf("%s (dir, size=%d)", fs.Name, fs.TotalSize())
	}
	return fmt.Sprintf("%s (file, size=%d)", fs.Name, fs.Size)
}

func (fs *FsNode) Path() string {
	path := ""
	x := fs
	if x.Parent == nil {
		return "/"
	}
	for {
		if x.Parent == nil {
			break
		} else {
			path = "/" + x.Name + path
		}
		x = x.Parent
	}
	return path
}

func (fs *FsNode) FilterDirs(filter func(*FsNode) bool) []*FsNode {
	found := []*FsNode{}
	if filter(fs) {
		found = append(found, fs)
	}
	for _, node := range fs.Children {
		if node.Size > 0 {
			continue
		}
		found = append(found, node.FilterDirs(filter)...)
	}
	return found
}

func (fs *FsNode) ChildNames() []string {
	names := make([]string, len(fs.Children))
	i := 0
	for name := range fs.Children {
		names[i] = name
		i++
	}
	sort.Strings(names)
	return names
}

func (fs *FsNode) AddFile(name string, size int) {
	if _, ok := fs.Children[name]; !ok {
		fs.Children[name] = &FsNode{name, fs, nil, size}
	}
}

func (fs *FsNode) MkDir(name string) *FsNode {
	exists, ok := fs.Children[name]
	if ok {
		return exists
	}
	nd := &FsNode{name, fs, map[string]*FsNode{}, 0}
	fs.Children[name] = nd
	return nd
}

func (fs *FsNode) TotalSize() int {
	var tSize int = 0
	for _, child := range fs.Children {
		if child.Size == 0 {
			tSize += child.TotalSize()
		} else {
			tSize += child.Size
		}
	}
	return tSize
}

func (fs *FsNode) Lsr(depth int) {
	if depth == 0 {
		fmt.Printf("%*s %s\n", depth*2, "-", fs)
	} else {
		fmt.Printf(" %*s %s\n", depth*2, "-", fs)
	}
	for _, c := range fs.ChildNames() {
		child := fs.Children[c]
		if child.Size == 0 {
			child.Lsr(depth + 1)
		} else {
			fmt.Printf(" %*s %s\n", depth*2+2, "-", child)
		}
	}
}

func (fs *FsNode) FileNames() []string {
	fNames := make([]string, len(fs.Children))
	i := 0
	for name, node := range fs.Children {
		if node.Size > 0 {
			fNames[i] = name
			i++
		}
	}
	fNames = fNames[0:i]
	sort.Strings(fNames)
	return fNames
}

func (fs *FsNode) DirNames() []string {
	dNames := make([]string, len(fs.Children))
	i := 0
	for name, node := range fs.Children {
		if node.Size == 0 {
			dNames[i] = name
			i++
		}
	}
	dNames = dNames[0:i]
	sort.Strings(dNames)
	return dNames
}
