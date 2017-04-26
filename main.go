package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
)

const maxFileSize = 1024 * 1024

func main() {
	var inFile string
	if len(os.Args) > 1 {
		inFile = os.Args[1]
	} else {
		inFile = "pkgfile"
	}
	gopath := os.Getenv("GOPATH")
	inFile = gopath + "/src/github.com/davidhenrygao/goInstallBinariesTool/" + inFile

	f, err := os.Open(inFile)
	if err != nil {
		fmt.Printf("Open file %s error: %v.\n", inFile, err)
		return
	}
	defer f.Close()
	fInfo, err := f.Stat()
	if err != nil {
		fmt.Printf("Read file info error: %v.\n", err)
		return
	}
	if fInfo.Size() > maxFileSize {
		fmt.Printf("File size exceeds %d bytes(%d bytes).\n",
			maxFileSize, fInfo.Size())
		return
	}
	fContent, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Printf("Read file error: %v.\n", err)
		return
	}

	regStr := "(?m)packages *= *\\[[\\s\\S]*?\\][\\.|,|;]?[\\r]?$"
	reg := regexp.MustCompile(regStr)
	strFind := reg.Find(fContent)
	if strFind == nil {
		fmt.Printf("input file format error!\n")
		return
	}
	regStr = "\\\".*\\\""
	reg = regexp.MustCompile(regStr)
	pkgs := reg.FindAll(strFind, -1)
	if pkgs == nil {
		fmt.Printf("input file pkg format error!\n")
		return
	}

	cmd := exec.Command("go", "get", "-u", "github.com/gpmgo/gopm")
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	fmt.Printf("Running go get gopm...\n")
	err = cmd.Run()
	if err != nil {
		fmt.Printf("Run cmd error: stdout(%s) \t stderr(%s).\n",
			out.String(), errBuf.String())
		return
	}

	var errFlag bool = false
	var pkgStr string
	for _, pkg := range pkgs {
		pkgStr = string(pkg[1 : len(pkg)-1])
		cmd = exec.Command("gopm", "get", "-g", "-d", pkgStr)
		fmt.Printf("gopm get %s.\n", pkgStr)
		err = cmd.Run()
		if err != nil {
			errFlag = true
			fmt.Printf("gopm get %s error: stdout(%s) \t stderr(%s).\n",
				pkgStr, out.String(), errBuf.String())
			continue
		}
		cmd = exec.Command("go", "install", pkgStr)
		fmt.Printf("go install %s.\n", pkgStr)
		err = cmd.Run()
		if err != nil {
			errFlag = true
			fmt.Printf("go install %s error: stdout(%s) \t stderr(%s).\n",
				pkgStr, out.String(), errBuf.String())
			continue
		}
	}

	if errFlag {
		fmt.Printf("Some error happen!\n")
	} else {
		fmt.Printf("Program Done.\n")
	}
}
