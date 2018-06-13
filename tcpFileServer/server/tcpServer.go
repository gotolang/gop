package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
)

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()
	files, err := ioutil.ReadDir("/Users/damao/Downloads/phstock")
	checkErr(err)
	for _, file := range files {
		fmt.Println(file.Name())
		_, err = conn.Write([]byte(file.Name()))
		// _, err = conn.Write([]byte(file.Size()))
		// for {
		// 	_, err := f.Read(buffer)
		// 	if err == io.EOF {
		// 		break
		// 	}
		// 	conn.Write(buffer)
		// }
	}
}

func downloadFiles() {

}

func main() {
	l, err := net.Listen("tcp", "localhost:9090")
	checkErr(err)
	defer l.Close()

	for {
		conn, err := l.Accept()
		checkErr(err)
		go handleClient(conn)
	}

}
