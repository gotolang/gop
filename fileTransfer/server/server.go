package main

import (
	"crypto/md5"
	"fmt"
	"html/template"
	"io"
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

}

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {

	} else {

	}
}

func main() {
	http.HandleFunc("/login", login)
	http.HandleFunc("/download", download)
	http.HandleFunc("/upload", upload)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal(err)
	}
}
