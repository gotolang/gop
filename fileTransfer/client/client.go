package main

import (
	"fmt"
	"log"
	"time"
	"github.com/burntsushi/toml"
)

type server struct {
	ip string
	port int
}

type client struct {
	app string
	lastUpdateAt time.Time
}

type downloadConfig struct {
	serverConf server
	clientConf client
}

func checkErr(err error) {
	log.Fatal(err)
	return
}

func main() {
	var downloadConf downloadConfig
	r, err := toml.DecodeFile("download.ini", &downloadConf)
	checkErr(err)
	fmt.Println(r)
	fmt.Println(downloadConf)
}