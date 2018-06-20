package main

import (
	"crypto/md5"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func upload(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		curTime := time.Now().Unix()
		md5Hash := md5.New()
		io.WriteString(md5Hash, strconv.FormatInt(curTime, 10))
		token := fmt.Sprintf("%x", md5Hash.Sum(nil))

		t, _ := template.ParseFiles("upload.gtpl")
		t.Execute(w, token)
	} else {
		r.ParseMultipartForm(32 << 20)
		file, handle, err := r.FormFile("uploadfile")
		if err != nil {
			log.Fatal(err)
			return
		}

		defer file.Close()

		f, err := os.OpenFile("./"+handle.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			log.Fatal(err)
			return
		}
		defer f.Close()

		io.Copy(f, file)
	}
}

func howManyFiles(w http.ResponseWriter, r *http.Request) {

}

func download(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		r.ParseForm()
		appName := r.Form.Get("app")
		appLoc := "/Users/damao/Downloads/applist/"
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
		fmt.Println(time.Now())
		fmt.Println(r.Host + " request app list...")

		zips, err := ioutil.ReadDir("/Users/damao/Downloads/applist")
		if err != nil {
			fmt.Println(err)
			return
		}
		for i, v := range zips {
			fmt.Println(i, v.Name())
			w.Write([]byte(v.Name()))
			if i != len(zips)-1 {
				w.Write([]byte(";"))
			}
			// w.Write([]byte(v.Name() + " "))
		}
		fmt.Println("====================end=========================")
	}

}

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {

	} else {

	}
}

func main() {
	http.HandleFunc("/applist", listApps)
	http.HandleFunc("/login", login)
	http.HandleFunc("/download", download)
	http.HandleFunc("/upload", upload)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal(err)
	}
}
