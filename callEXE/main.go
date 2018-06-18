package main

import (
	"fmt"
	"os/exec"
	"runtime"
)

func main() {
	if runtime.GOOS == "windows" {
		fmt.Println(runtime.GOOS)
	} else if runtime.GOOS == "darwin" {
		fmt.Println("macos")
		cmd := exec.Command("/Applications/Postman.app", "")
		err := cmd.Start()
		if err != nil {
			fmt.Println(err)
		}

	} else {
		fmt.Println("linux")
	}

}
