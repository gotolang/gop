package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const appLoc string = "/Users/damao/Downloads/applist/"

func download(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		r.ParseForm()
		appName := r.Form.Get("app")
		// fmt.Println("ioutil read file " + appLoc + appName)
		// f, err := ioutil.ReadFile(appLoc + appName)
		// if err != nil {
		// 	fmt.Println(err)
		// 	return
		// }
		// fmt.Println("Start write to client...")
		// _, err = w.Write(f)
		// if err != nil {
		// 	fmt.Println(err)
		// 	return
		// }
		file, err := os.Open(appLoc + appName)
		if err != nil {
			log.Println(err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		defer file.Close()
		_, err = io.Copy(w, file)
		if err != nil {
			log.Println(err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	} else {

	}
}

func listApps(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		fmt.Println("=====================start=======================")
		log.Println(r.Host + " request app list...")

		zips, err := ioutil.ReadDir(appLoc)
		if err != nil {
			fmt.Println(err)
			return
		}

		var apps []byte
		// apps = make([]byte, 1024)
		for i, v := range zips {
			// fmt.Println(i, v.Name())
			apps = append(apps, v.Name()...)
			if i != len(zips)-1 {
				apps = append(apps, ";"...)
			}
		}
		for i, v := range apps {
			fmt.Println(i, string(v))
		}
		n, err := w.Write(apps)
		// fmt.Fprint(w, apps)
		if err != nil {
			log.Println(err)
		}
		fmt.Println(n)
		fmt.Println("====================end=========================")
	}

}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello,world!"))
}

func main() {

	http.HandleFunc("/applist", listApps)
	http.HandleFunc("/download", download)
	http.HandleFunc("/hello", hello)

	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal(err)
	}
}
