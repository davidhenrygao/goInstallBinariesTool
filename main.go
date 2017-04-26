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
const githubpath string = "/github.com/davidhenrygao/goInstallBinariesTool/"

func OpenAndCheckFile(file string) *os.File {
	f, err := os.Open(file)
	if err != nil {
		fmt.Printf("Open file %s error: %v.\n", file, err)
		return nil
	}
	fInfo, err := f.Stat()
	if err != nil {
		fmt.Printf("Read file info error: %v.\n", err)
		f.Close()
		return nil
	}
	if fInfo.Size() > maxFileSize {
		fmt.Printf("File size exceeds %d bytes(It is %d bytes).\n",
			maxFileSize, fInfo.Size())
		f.Close()
		return nil
	}
	return f
}

func FindPkgsInFile(f *os.File) [][]byte {
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

/*
func ExecCmd(cmdline string) error {

}
*/

func main() {
	var inFile string
	if len(os.Args) > 1 {
		inFile = os.Args[1]
	} else {
		inFile = "pkgfile"
	}
	gopath := os.Getenv("GOPATH") + "/src"
	inFile = gopath + githubpath + inFile

	f := OpenAndCheckFile(inFile)
	if f == nil {
		return
	}
	defer f.Close()

	pkgs := FindPkgsInFile(f)
	if pkgs == nil {
		return
	}

	cmd := exec.Command("go", "get", "-u", "github.com/gpmgo/gopm")
	var errBuf bytes.Buffer
	cmd.Stderr = &errBuf
	fmt.Printf("Running go get gopm...\n")
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Run cmd error: \n%s\n",
			errBuf.String())
		return
	}

	var errFlag bool = false
	var pkgStr string
	for _, pkg := range pkgs {
		pkgStr = string(pkg[1 : len(pkg)-1])
		cmd = exec.Command("gopm", "get", "-g", "-d", pkgStr)
		cmd.Stderr = &errBuf
		fmt.Printf("gopm get %s.\n", pkgStr)
		err = cmd.Run()
		if err != nil {
			errFlag = true
			fmt.Printf("gopm get %s error: \n%s\n",
				pkgStr, errBuf.String())
			continue
		}
		cmd = exec.Command("go", "install", pkgStr)
		cmd.Stderr = &errBuf
		fmt.Printf("go install %s.\n", pkgStr)
		err = cmd.Run()
		if err != nil {
			errFlag = true
			fmt.Printf("go install %s error: \n%s\n",
				pkgStr, errBuf.String())
			continue
		}
	}

	if errFlag {
		fmt.Printf("Some error happen!\n")
	} else {
		fmt.Printf("Program Done.\n")
	}
}
