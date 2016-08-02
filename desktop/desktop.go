// Package desktop contains OS specific functions.
package desktop

import (
	"os"
	"os/exec"
	"runtime"
	//	"fmt"
)

var commands = map[string]string{
	"windows": "start",
	"darwin":  "open",
	"linux":   "xdg-open",
}

func run(cmd *exec.Cmd) {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err := cmd.Start()
	if err != nil {
		panic(err)
	}
	err = cmd.Wait()
	if err != nil {
		panic(err)
	}
}

// Open calls the OS default program for a URI.
func Open(uri string) (err error) {
	if runtime.GOOS == "windows" {
		run(exec.Command("cmd", "/c", "start", uri))
	}
	if runtime.GOOS == "darwin" {
		err = exec.Command("open", uri).Start()
	}
	if runtime.GOOS == "linux" {
		err = exec.Command("xdg-open", uri).Start()
	}
	//	fmt.Println("Client started.")
	return
}
