package main

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
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
	// App           string
	DownloadFrom  string
	LastUpdatedAt time.Time
}

var app string

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
}

func decodeTOML(tomlfile string) (string, error) {
	var downloadConf downloadConfig
	_, err := toml.DecodeFile(tomlfile, &downloadConf)
	checkErr(err)
	servers := downloadConf.Servers
	dst := downloadConf.Client.DownloadFrom

	ipAndPort, ok := servers[dst]

	if !ok {
		return "", errors.New(dst + " does not exist")
	}

	ip := ipAndPort.IP
	port := strconv.Itoa(ipAndPort.Port)
	app = downloadConf.Client.App

	url := "http://" + ip + ":" + port + "/download?app=" + app

	return url, nil
}

func getZipFromServer(url string) (io.Reader, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return resp.Body, nil
}

func unZip2Localfile(respBody io.Reader) (bool, error) {
	//createLocalfile
	localZip, err := os.Create(app + ".zip")
	if err != nil {
		return false, err
	}
	defer localZip.Close()

	//copy resp to localfile
	_, err = io.Copy(localZip, respBody)
	if err != nil {
		return false, err
	}

	//unzip file from server
	//is directory or file
	//if dir then mkdir and it's files
	rc, err := zip.OpenReader(localZip.Name())
	if err != nil {
		return false, err
	}
	defer rc.Close()

	for _, file := range rc.Reader.File {
		frc, err := file.Open()
		if err != nil {
			return false, err
		}
		defer frc.Close()

		if file.FileInfo().IsDir() {
			os.MkdirAll("./", file.Mode())
		} else {
			lf, err := os.Create(file.Name)
			if err != nil {
				return false, err
			}

			defer lf.Close()

			_, err = io.Copy(lf, frc)
			if err != nil {
				return false, err
			}
		}
	}

	return true, nil
}

func main() {

	tomlFile := "download.toml"
	url, err := decodeTOML(tomlFile)
	checkErr(err)

	respBody, err := getZipFromServer(url)
	checkErr(err)

	b, err := unZip2Localfile(respBody)
	checkErr(err)

	fmt.Println(b)

}
