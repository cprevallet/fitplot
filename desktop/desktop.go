/* The desktop package exposes Open function which acts like the user clicked on a uri. */
package desktop

import (
	"os"
	"os/exec"
	"runtime"
	"fmt"
)

var commands = map[string]string{
	"windows": "start",
	"darwin":  "open",
	"linux":   "xdg-open",
}

var Version = "0.1.0"

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

// Open calls the OS default program for uri
func Open(uri string) (err error) {
	if runtime.GOOS == "windows" {
		run(exec.Command("cmd", "/c", "start", uri))
	}
	if runtime.GOOS == "darwin" {
		run(exec.Command("open", uri))
	}
	if runtime.GOOS == "linux" {
		run(exec.Command("xdg-open", uri))
	}
	fmt.Println("Client started.")
	return
}
