package main

import (
	"log"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatal("command line must be in 'go-envdir /path/to/env/dir command arg1 arg2' format")
		return
	}

	dir := os.Args[1]

	env, err := ReadDir(dir)
	if err != nil {
		log.Fatalf("failed read configuration dir, %s", err.Error())
		return
	}

	code := RunCmd(os.Args[2:], env)
	os.Exit(code)

}
