package main

import (
	"fmt"
	"log"
	"time"

	"github.com/BurntSushi/toml"
)

type downloadConfig struct {
	Title   string
	Servers map[string]server
	Client  client
}

type server struct {
	IP   string
	Port int
}

type client struct {
	App           string
	DownloadFrom  string
	LastUpdatedAt time.Time
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
}

func main() {
	var downloadConf downloadConfig
	_, err := toml.DecodeFile("download.toml", &downloadConf)
	checkErr(err)
	fmt.Println(downloadConf.Client.App)
	fmt.Println(downloadConf.Client.DownloadFrom)
	fmt.Println(downloadConf.Servers["beta"].IP)
	fmt.Println(downloadConf.Servers["beta"].Port)

}
