package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Required minimum 2 args")
		os.Exit(1)
	}

	dir := os.Args[1]
	args := os.Args[2:]

	env, err := ReadDir(dir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	os.Exit(RunCmd(args, env))
}
