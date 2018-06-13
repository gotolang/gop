package main

import (
	"fmt"
	"log"
	"net/http"
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

	ip := downloadConf.Servers["beta"].IP
	port := downloadConf.Servers["beta"].Port
	app := downloadConf.Client.App

	// 获取文件个数

	//
	url := "http://" + ip + string(port) + "/download?app=" + app
	resp, err := http.Get(url)
	checkErr(err)

}
