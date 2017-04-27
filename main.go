package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

const maxFileSize = 1024 * 1024
const githubpath string = "/github.com/davidhenrygao/goInstallBinariesTool/"

func closeFile(f *os.File) {
	err := f.Close()
	if err != nil {
		fmt.Printf("Close file error: %v.\n", err)
	}
}

func openAndCheckFile(file string) *os.File {
	f, err := os.Open(file)
	if err != nil {
		fmt.Printf("Open file %s error: %v.\n", file, err)
		return nil
	}
	fInfo, err := f.Stat()
	if err != nil {
		fmt.Printf("Read file info error: %v.\n", err)
		closeFile(f)
		return nil
	}
	if fInfo.Size() > maxFileSize {
		fmt.Printf("File size exceeds %d bytes(It is %d bytes).\n",
			maxFileSize, fInfo.Size())
		closeFile(f)
		return nil
	}
	return f
}

func findPkgsInFile(f *os.File) [][]byte {
	fContent, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Printf("Read file error: %v.\n", err)
		return nil
	}
	regStr := "(?m)packages *= *\\[[\\s\\S]*?\\][\\.|,|;]?[\\r]?$"
	reg := regexp.MustCompile(regStr)
	strFind := reg.Find(fContent)
	if strFind == nil {
		fmt.Printf("input file format error!\n")
		return nil
	}
	regStr = "\\\".*\\\""
	reg = regexp.MustCompile(regStr)
	pkgs := reg.FindAll(strFind, -1)
	if pkgs == nil {
		fmt.Printf("input file pkg format error!\n")
	}
	return pkgs
}

func execCmd(cmdline string) error {
	cmdStrs := strings.Fields(cmdline)
	if len(cmdStrs) == 0 {
		return nil
	}
	prog := cmdStrs[0]
	args := cmdStrs[1:]
	cmd := exec.Command(prog, args...)
	var errBuf bytes.Buffer
	cmd.Stderr = &errBuf
	fmt.Printf("Run Cmd: %s\n", cmdline)
	err := cmd.Run()
	if err != nil {
		fmt.Printf("%s error: \n%s\n", cmdline, errBuf.String())
		return err
	}
	return nil
}

func main() {
	var inFile string
	if len(os.Args) > 1 {
		inFile = os.Args[1]
	} else {
		inFile = "pkgfile"
	}
	gopath := os.Getenv("GOPATH") + "/src"
	inFile = gopath + githubpath + inFile

	f := openAndCheckFile(inFile)
	if f == nil {
		return
	}
	defer closeFile(f)

	pkgs := findPkgsInFile(f)
	if pkgs == nil {
		return
	}

	err := execCmd("go get -u github.com/gpmgo/gopm")
	if err != nil {
		return
	}

	var errFlag = false
	var pkgStr string
	for _, pkg := range pkgs {
		pkgStr = string(pkg[1 : len(pkg)-1])
		err = execCmd("gopm get -g -d " + pkgStr)
		if err != nil {
			errFlag = true
			continue
		}
		err = execCmd("go install " + pkgStr)
		if err != nil {
			errFlag = true
			continue
		}
	}

	if errFlag {
		fmt.Printf("Some error happen!\n")
	} else {
		fmt.Printf("Program Done.\n")
	}
}
