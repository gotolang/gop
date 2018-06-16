package main

import (
	"archive/zip"
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func checkErr(err error, action string) {
	// fmt.Println(err)
	if err != nil {
		log.Println(action + " error")
		panic(err)
	}
}

func unzip(filename string) {
	zipReader, err := zip.OpenReader(filename)
	checkErr(err, "Zip open reader")
	defer zipReader.Close()

	for _, file := range zipReader.Reader.File {
		fmt.Println(file.Name)

		zipfile, err := file.Open()
		checkErr(err, "zip file open")
		defer zipfile.Close()

		osfile, err := os.Create(file.Name)
		checkErr(err, "zip create file")
		defer osfile.Close()

		_, err = io.Copy(osfile, zipfile)
		checkErr(err, "zip io copy")

		fmt.Println("Extract complete.")
	}
}

func main() {
	// urlAppList := "http://172.42.1.221:9999/applist"
	urlAppList := "http://localhost:9090/applist"
	respAppList, err := http.Get(urlAppList)
	checkErr(err, "GET applist")
	defer respAppList.Body.Close()

	p := make([]byte, 1024)
	respAppList.Body.Read(p)
	var apps map[int64]string
	apps = make(map[int64]string)
	appsFromServer := strings.Split(string(bytes.Trim(p, "\x00")), ";")
	// fmt.Println(len(appsFromServer))
	for i, v := range appsFromServer {
		fmt.Println(i, v)
		apps[int64(i)] = v
		// fmt.Println("assignment:", i, v)
	}

	userInput := bufio.NewScanner(os.Stdin)
	fmt.Printf("Which APP do you want to download? Input the number:")
	userInput.Scan()
	text := userInput.Text()
	// fmt.Println(text)
	appIndex, err := strconv.ParseInt(text, 10, 64)
	checkErr(err, "ParseInt")
	appName, ok := apps[appIndex]
	if !ok {
		log.Println("index does not exist")
		return
	}

	urlAppDownload := "http://localhost:9090/download?app=" + appName
	// urlAppDownload := "http://172.42.1.221:9999/download"
	// fmt.Println(urlAppDownload)
	respAppDownload, err := http.Get(urlAppDownload)
	checkErr(err, "GET download app")
	// fmt.Println(respAppDownload.StatusCode)
	defer respAppDownload.Body.Close()

	file, err := os.Create(appName)
	// f, err := os.OpenFile(appName, os.O_RDONLY|os.O_CREATE, 0666)
	checkErr(err, "Create file")
	defer file.Close()

	_, err = io.Copy(file, respAppDownload.Body)
	checkErr(err, "io copy file")
	fmt.Println("Download complete.")

	fmt.Println("Start extract files...")
	unzip(appName)

}
