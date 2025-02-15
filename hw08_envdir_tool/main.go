package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		panic("mmm hueta")
	}
	
	dir := os.Args[1]
	args := os.Args[2:]

	env, err := ReadDir(dir)

	if err != nil {
		fmt.Println(err)
		return
	}

	os.Exit(RunCmd(args, env))
}
