/**
  deps - dependency management utility for Linux executables
*/
package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type dependency struct {
	Name string
	Path string
	Addr string
}

type GetDepsFunc func(lddOutput string) []dependency

func getDependencies(lddOutput string) []dependency {
	lines := strings.Split(lddOutput, "\n")
	depsCount := len(lines) - 1
	deps := make([]dependency, depsCount)

	for i, l := range lines {
		dep := strings.Split(l, " ")

		// ldd has a little weird output:
		// 1) name => path (addr)
		//    libc.so.6 => /lib/i386-linux-gnu/libc.so.6 (0xb753d000)
		// 2) path (addr)
		//    /lib/ld-linux.so.2 (0xb7746000)
		// 3) and a strange (this is VDSO provided by kernel, no need to copy it):
		//     linux-gate.so.1 =>  (0xb7745000)
		// so handle different types of it

		switch {
		case len(dep) > 3: // case 1
			deps[i].Name = strings.TrimSpace(dep[0])
			deps[i].Path = strings.TrimSpace(dep[2])
			deps[i].Addr = strings.TrimSpace(dep[3])
		case len(dep) > 1: // case 2
			idx := strings.LastIndex(dep[0], "/")
			deps[i].Name = strings.TrimSpace(string(dep[0][idx+1:]))
			deps[i].Path = strings.TrimSpace(dep[0])
			deps[i].Addr = strings.TrimSpace(dep[1])
		}
	}

	return deps
}

var visitedDeps map[string]bool

func walkDeps(path, dir string, getDeps GetDepsFunc) error {
	// Do not walk among visited paths,
	// or else the recursion will be infinite
	if _, ok := visitedDeps[path]; ok {
		fmt.Printf("! %s already copied\n", path)
		return nil
	}

	// Check this path as visited
	visitedDeps[path] = true

	// Get dependency list by executing ldd
	ldd := exec.Command("ldd", path)
	var lddOut bytes.Buffer
	ldd.Stdout = &lddOut
	ldd.Run()

	// Parse ldd's output
	deps := getDeps(lddOut.String())

	// Walk through new dependencies recursively to gather them all
	for _, d := range deps {
		fmt.Printf("-Dependency:\n  name: %s\n  path: %s\n  addr: %s\n", d.Name, d.Path, d.Addr)

		cp := exec.Command("cp", "--parent", d.Path, dir)
		cp.Run()

		if err := walkDeps(d.Path, dir, getDeps); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: deps EXECUTABLE PATH\nGather dependencies of an executable")
		return
	}

	what := os.Args[1]
	whereToCopy := os.Args[2]

	if err := os.Mkdir(whereToCopy, 0777); err != nil {
		panic("Could not create output directory - " + err.Error())
	}

	visitedDeps = make(map[string]bool)
	if err := walkDeps(what, whereToCopy, getDependencies); err != nil {
		panic("An error has occured while gathering dependencies - " + err.Error())
	}

	fmt.Println("Finished!")
}
