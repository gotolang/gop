package main

import (
	"log"
	"os/exec"
)

func main() {
	log.Println("call windows exe installer")
	filepath := "c:\\emrs.exe"
	cmd := exec.Command("cmd", "/C", "start", filepath)
	err := cmd.Start()
	if err != nil {
		log.Println(err)
		return
	}
}
