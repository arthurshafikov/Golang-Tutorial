package main

import (
	"os"
)

func main() {
	if len(os.Args) < 2 {
		panic("Not enough args")
	}
	env, err := ReadDir(os.Args[1])
	if err != nil {
		panic(err)
	}
	cmd := os.Args[2:]
	RunCmd(cmd, env)
}
