package main

import (
	"fmt"
	"log"
	"net"
)

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	conn, err := net.Dial("tcp", "localhost:9090")
	checkErr(err)
	defer conn.Close()
	b := make([]byte, 1024)
	for {
		conn.Read(b)
		fmt.Println(string(b))
	}

}
