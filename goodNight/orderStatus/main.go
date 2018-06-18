package main

import (
	"database/sql"

	_ "github.com/mattn/go-oci8"
)

func main() {
	db, err := sql.Open("oci8", "hw/a@dbtest")
	if err ! nil {
		log.Println(err)
		return
	}

	db.Query("")

}
