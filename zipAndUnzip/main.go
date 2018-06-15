package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func checkErr(err error, action string) {
	// fmt.Println(err)
	if err != nil {
		fmt.Println(action + " error")
		panic(err)
	}

	// os.Exit(-1)
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
	appsFromServer := strings.Split(string(p), ";")
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
		fmt.Println("index does not exist")
		return
	}
	urlAppDownload := "http://localhost:9090/download?app=" + appName
	// urlAppDownload := "http://172.42.1.221:9999/download"

	respAppDownload, err := http.Get(urlAppDownload)
	checkErr(err, "GET download app")
	defer respAppDownload.Body.Close()

	f, err := os.Create(appName)
	// f, err := os.OpenFile(appName, os.O_RDONLY|os.O_CREATE, 0666)
	checkErr(err, "Create file")
	defer f.Close()

	_, err = io.Copy(f, respAppDownload.Body)
	checkErr(err, "io copy file")
	fmt.Println("Download complete.")
}
