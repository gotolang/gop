package main

import (
	"database/sql"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func struct2XML() error {
	type User struct {
		XMLName  xml.Name `xml:"Person"`
		ID       string   `xml:"id,attr"`
		Tel      string   `xml:"tel"`
		NickName string   `xml:"nickName"`
		Email    string   `xml:"email"`
	}

	db, err := sql.Open("postgres", "user=postgres password=postgres dbname=aiwriter sslmode=disable")
	if err != nil {
		return err
	}
	fmt.Println(db)
	rows, err := db.Query("select tel, nick_name, email from \"user\" where tel = $1", "18877777777")
	if err != nil {
		return err
	}

	defer rows.Close()
	var person User

	for rows.Next() {

		err = rows.Scan(&person.Tel, &person.NickName, &person.Email)
		if err != nil {
			return err
		}
	}
	person.ID = "1"
	b, err := xml.Marshal(person)
	if err != nil {
		return err
	}
	// file, err := os.Create("person.xml")
	// if err != nil {
	// 	log.Println(err)
	// 	return err
	// }
	err = ioutil.WriteFile("person.xml", b, 0644)
	if err != nil {
		return err
	}
	return nil
}

func xml2Struct() (interface{}, error) {
	type NickName struct {
		ID   string `xml:"id,attr"`
		Name string `xml:",chardata"`
	}
	type User struct {
		XMLName  xml.Name `xml:"user"`
		ID       string   `xml:"id,attr"`
		Tel      string   `xml:"tel"`
		Nickname NickName `xml:"nick_Name"`
		Email    string   `xml:"email"`
	}

	file, err := os.Open("user.xml")
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	var user User
	err = xml.Unmarshal(data, &user)
	if err != nil {
		return nil, err
	}
	// fmt.Println(user)
	return user, nil
}

func main() {
	err := struct2XML()
	if err != nil {
		log.Println(err)
	}
	v, err := xml2Struct()
	if err != nil {
		log.Println(err)
	}

	fmt.Println(v)

}
