//go:build ignore

//go:generate go run generate.go ../proto-src

package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

var packageRe = regexp.MustCompile(`(?m)package PB_(.*?);`)

const baseImportPath = "github.com/flipperdevices/go-flipper/internal/proto"

func main() {
	for _, path := range os.Args[1:] {
		matches, err := filepath.Glob(filepath.Join(path, "*proto"))
		if err != nil {
			log.Fatalln(err)
		}
		args := []string{
			"--proto_path", path,
			"--go_out=module=" + filepath.Dir(baseImportPath) + ":..",
		}
		for _, m := range matches {
			p, err := composeImportPath(m)
			if err != nil {
				log.Fatalln("Can't compose import path", m, err)
			}
			args = append(args, fmt.Sprintf("--go_opt=M%s=%s%s", filepath.Base(m), baseImportPath, p))
		}
		args = append(args, matches...)
		fmt.Println(args)
		cmd := exec.Command("protoc", args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			log.Fatal("Failed generating", path)
		}
	}
}

func composeImportPath(file string) (string, error) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}
	matches := packageRe.FindStringSubmatch(string(b))
	switch len(matches) {
	case 2:
		return "/" + strings.ToLower(matches[1]), nil
	case 0:
		return "", nil
	}
	return "", errors.New("regex failed")
}
