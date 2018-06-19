package main

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"
)

func main() {
	if runtime.GOOS == "windows" {
		fmt.Println(runtime.GOOS)
		cmd := exec.Command("C:\\Windows\\system32\\mstsc.exe") //

		err := cmd.Start()
		if err != nil {
			log.Println(err)
			return
		}
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
