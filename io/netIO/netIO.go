//https://github.com/toml-lang/toml/archive/master.zip

package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	url := "https://github.com/toml-lang/toml/archive/master.zip"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	file, err := os.Create("toml.zip")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer file.Close()

	io.Copy(file, resp.Body)
	if err := resp.Body.Close(); err != nil {
		fmt.Println(err)
	}
}
