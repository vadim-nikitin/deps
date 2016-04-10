/**
* deps - dependency management utility for Linux executables
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

func getDeps(lddOutput string) []dependency {
    lines := strings.Split(lddOutput, "\n")
    depsCount := len(lines) - 1
    deps := make([]dependency, depsCount)
        
    for i, l := range lines {
        dep := strings.Split(l, " ")
        
        switch {
        case len(dep) > 3:
            deps[i].Name = strings.TrimSpace(dep[0])
            deps[i].Path = strings.TrimSpace(dep[2])
            deps[i].Addr = strings.TrimSpace(dep[3])
        case len(dep) > 1:
            idx := strings.LastIndex(dep[0], "/")
            deps[i].Name = strings.TrimSpace(string(dep[0][idx + 1:]))
            deps[i].Path = strings.TrimSpace(dep[0])
            deps[i].Addr = strings.TrimSpace(dep[1])
        }
    }

    return deps
}

func walkDeps(what, where string) {
    ldd := exec.Command("ldd", what)
    var lddOut bytes.Buffer
    ldd.Stdout = &lddOut
    ldd.Run()
    
    deps := getDeps(lddOut.String())  
    for _, d := range deps {
        fmt.Printf("-Dependency:\n  name: %s\n  path: %s\n  addr: %s\n", d.Name, d.Path, d.Addr)
        cp := exec.Command("cp", "--parent", d.Path, where)
        cp.Run()
        walkDeps(d.Path, where)
    }
}

func main() {
    what := os.Args[1]
    where := os.Args[2]
    os.Mkdir(where, 0777)
    walkDeps(what, where)
}
