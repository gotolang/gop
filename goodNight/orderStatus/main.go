package main

import (
	"database/sql"
	"log"
	"strconv"
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	_ "github.com/mattn/go-oci8"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func gbk2Utf8(str []byte) (utf8Str []byte, err error) {
	reader := transform.NewReader(bytes.NewReader(str), simplifiedchinese.GBK.NewDecoder())
	utf8Str, err = ioutil.ReadAll(reader)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return utf8Str, nil
}

func userInteractive(db *sql.DB) (deptCode string, wardCode string, deptName string, err error) {
	// input name of dept
	fmt.Print("Input code to search:")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan() 

	deptQ := `
	select a.dept_code, b.ward_code, a.dept_name from dept_dict a, dept_vs_ward b 
	where a.dept_code = b.dept_code and 
	upper(a.input_code) like '%'||:1||'%'
	`

	// search dept_code according to name of dept
	rows, err := db.Query(deptQ, scanner.Text())
	if err != nil {
		log.Println(err)
		return "", "", "", err
	}

	defer rows.Close()

	var depts map[string][2]string

	for rows.Next() {
		var deptCode, wardCode, deptName string
		err := rows.Scan(&deptCode, &wardCode, &deptName)
		if err != nil {
			log.Println(err)
			return "", "", "", err
		}
		dn, err := gbk2Utf8([]byte(deptName))
		if err != nil {
			log.Println(err)
			return "", "", "", err
		}
		deptName = string(dn)
		// map the result (map[int]string)
		depts[deptCode] = [wardCode, deptName]
		fmt.Println(wardCode, deptCode, deptName)
	}
	err = Rows.Err()
	if err != nil {
		log.Println(err)
		return "", "", "", err
	}
	// user select the dept
	fmt.Print("input dept code:")
	scanner.Scan()
	dept4Update := scanner.Text()
	wardAndName, ok := depts[dept4Update]
	if !ok {
		log.Println("dept code does not exist")
		continue
	}
	return dept4Update, wardAndName[0], wardAndName[1], nil
}

func updateLSOrderStatus(db *sql.DB, deptCode string, wardCode string) (int64, error) {
	sqlOrder := `
	update orders d set d.order_status = 2
	where (d.order_status = 3 or (d.order_status = 3 and d.stop_date_time is null))
	  and d.repeat_indicator = 0
	  and d.patient_id in (select patient_id from pats_in_hospital p where p.ward_code = :1 and p.dept_code = :2)
	  `
	r, err := db.Exec(sqlOrder, wardCode, deptCode)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	affectedRows, err := r.RowsAffected()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return affectedRows, nil
} 

func updateCQOrderStatus(db *sql.DB, deptCode string, wardCode string) (int64, error) {
	sqlOrder := `
	update orders d set d.order_status = 2
	where (d.order_status = 3 and d.stop_date_time is null)
	  and d.repeat_indicator = 1
	  and d.patient_id in (select patient_id from pats_in_hospital p where p.ward_code = :1 and p.dept_code = :2)`

	r, err = db.Exec(sqlOrder, wardCode, deptCode)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	affectedRows, err = r.RowsAffected()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return affectedRows, nil
}

func main() {
	os.Setenv("NLS_LANG", "AMERICAN_AMERICA.US7ASCII")
	db, err := sql.Open("oci8", "golang/golang@dbtest")
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()

	wardCode, deptCode, deptName, err := userInteractive(db)
	if err != nil {
		log.Println(err)
		return
	}

	affectedRows, err := updateLSOrderStatus(db, wardCode, deptCode)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println("临时医嘱影响行数：" + strconv.FormatInt(affectedRows, 10))

	affectedRows, err = updateCQOrderStatus(db, wardCode, deptCode)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println("长期医嘱影响行数：" + strconv.FormatInt(affectedRows, 10))
}
