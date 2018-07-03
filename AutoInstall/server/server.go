package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/BurntSushi/toml"
)

const appLoc string = "C:\\HIS\\HIS装机程序\\"

type apps struct {
	Apps map[string]app
}

type app struct {
	Chinese string
	Zip     string
	Dir     string
	Ini     string
	Exe     string
	Desktop bool
}

func decodeTOML(tomlfile string) (map[string]app, error) {
	var applications apps

	_, err := toml.DecodeFile(tomlfile, &applications)
	if err != nil {
		return nil, err
	}

	// fmt.Println(applications.Apps)
	return applications.Apps, nil
}

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

		apps, err := decodeTOML("apps.toml")
		if err != nil {
			log.Println(err)
		}

		for i := range apps {
			b, err := json.Marshal(apps[i])
			if err != nil {
				log.Println(err)
			}

			_, err = w.Write(b)

			if err != nil {
				log.Println(err)
			}
			// fmt.Println(n)
		}
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
