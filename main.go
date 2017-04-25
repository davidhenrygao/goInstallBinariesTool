package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

func main() {
	var inFile string
	if len(os.Args) > 1 {
		inFile = os.Args[1]
	} else {
		inFile = "go.vim"
	}
	f, err := os.Open(inFile)
	if err != nil {
		fmt.Printf("Open file %s error: %v.\n", inFile, err)
		return
	}
	defer f.Close()

	cmd := exec.Command("go", "get", "-u", "github.com/gpmgo/gopm")
	var out bytes.Buffer
	cmd.Stdout = &out
	fmt.Printf("Running go get gopm...\n")
	err = cmd.Run()
	if err != nil {
		fmt.Printf("Run cmd error: %v.\n", err)
		return
	}
	fmt.Printf("Run cmd output: %s.\n", out.String())
	gopath := os.Getenv("GOPATH")
	fmt.Printf("GOPATH(%s).\n", gopath)
}
