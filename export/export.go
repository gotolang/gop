package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	fmt.Println("export")
	db, err := sql.Open("postgres", "user=postgres password=postgres dbname=aiwriter sslmode=disable")
	if err != nil {
		log.Println(err)
	}
	fmt.Println(db)
	rows, err := db.Query("select nick_name from \"user\" where tel = $1", "15012535569")
	if err != nil {
		log.Println(err)
	}

	defer rows.Close()

	for rows.Next() {
		var nickName string
		err = rows.Scan(&nickName)
		if err != nil {
			log.Println(err)
		}
		fmt.Println(nickName)
	}
}
