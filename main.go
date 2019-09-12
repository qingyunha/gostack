package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/go-delve/delve/pkg/proc"
	"github.com/go-delve/delve/pkg/proc/native"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %v <pid>\n", os.Args[0])
		os.Exit(1)
	}
	pid, _ := strconv.Atoi(os.Args[1])
	p, err := native.Attach(pid, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	gs, _, err := proc.GoroutinesInfo(p, 0, 0)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, g := range gs {
		if g.Status == proc.Gdead || g.Status == proc.Gidle {
			continue
		}
		stack, err := g.Stacktrace(50, false)
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
		fmt.Printf("goroutine %v [%v]\n", g.ID, g.WaitReason)
		printStack(stack)
		fmt.Println()
	}
	p.Detach(false)
}

func printStack(stack []proc.Stackframe) {
	for _, s := range stack {
		fmt.Printf("%v\n", s.Call.Fn.Name)
		fmt.Printf("\t%v:%v\n", s.Call.File, s.Call.Line)
	}
}
