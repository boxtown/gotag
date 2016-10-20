package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/boxtown/gotag"
)

var skips string

func init() {
	flag.StringVar(&skips, "skip", "", "Comma-separated list of test flags to skip")
}

func main() {
	for _, v := range strings.Split(skips, ",") {
		gotag.Skip(strings.TrimSpace(v))
	}
	cmd := exec.Command("go", "test")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println(err.Error())
	}
}
