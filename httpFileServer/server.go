package main

import (
	"log"
	"net/http"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("/Users/damao/Downloads/phstock")))
	err := http.ListenAndServe(":9090", nil)

	if err != nil {
		log.Fatal(err)
	}
}
