package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-oci8"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

type mrBasic struct {
	AKB020      string         `db:"AKB020"`
	AKB026      string         `db:"AKB026"`
	YKC700      string         `db:"YKC700"`
	YZY201      string         `db:"YZY201"`
	PKA435      string         `db:"PKA435"`
	PKA439      sql.NullString `db:"PKA439"`
	BZ1         sql.NullString `db:"BZ1"`
	BZ2         sql.NullString `db:"BZ2"`
	BZ3         sql.NullString `db:"BZ3"`
	DRBZ        sql.NullInt64  `db:"DRBZ"`
	EXPLANATION sql.NullString `db:"EXPLANATION"`
}

func gbk2Utf8(str []byte) (utf8Str []byte, err error) {
	reader := transform.NewReader(bytes.NewReader(str), simplifiedchinese.GBK.NewDecoder())
	utf8Str, err = ioutil.ReadAll(reader)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return utf8Str, nil
}

func main() {
	orclDB, err := sqlx.Open("oci8", "hw/damao421@dbserver")
	if err != nil {
		log.Println("oracle database open error, because ", err)
		return
	}
	defer orclDB.Close()

	msDB, err := sqlx.Open("mssql", "sqlserver://yy005106:yy005106@172.42.1.199")
	if err != nil {
		log.Println("sqlserver database open error, because ", err)
		return
	}
	defer msDB.Close()

	logFile, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Println(err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	sqlTextms := `SELECT A.AKB020
					,A.AKB026
					,A.YKC700
					,A.YZY201
					,A.PKA435
					,A.PKA439
					,A.BZ1
					,A.BZ2
					,A.BZ3
					,A.DRBZ
					,B.EXPLANATION
				FROM [yy005106].[dbo].[HIS_SNYD_BASYBRXX_TEST] A LEFT JOIN [yy005106].[dbo].[HIS_SNYD_BASYCOLUMN] B ON A.PKA435=B.PKA435 AND A.YZY201 = B.YZY201
				WHERE (A.PKA439 IS NULL) AND B.ISNULL = 'F' AND A.YKC700='0051061609769418'`

	rows, err := msDB.Queryx(sqlTextms)
	if err != nil {
		log.Println("sqlserver database query error, because ", err)
		return
	}
	defer rows.Close()

	var basics []mrBasic

	for rows.Next() {
		var basic mrBasic
		err = rows.StructScan(&basic)
		if err != nil {
			log.Println("sqlserver database StructScan error, because ", err)
			return
		}
		basics = append(basics, basic)
	}

	fmt.Printf("Need to update %d rows. Are you sure?\n", len(basics))

	sqlTextorcl := `select b.pka439
						from INPBILL.his_snyd_basybrxx b
					where b.ykc700 = :jydjh
						and b.yzy201 = :no
						and b.pka435 = :item`

	sqlText4Update := `UPDATE [yy005106].[dbo].[HIS_SNYD_BASYBRXX_TEST]
							SET PKA439 = $1
						WHERE AKB020 = '005106' AND YKC700 = $2 AND YZY201 = $3`

	for _, v := range basics {
		// fmt.Println(i, v.YKC700, v.YZY201, v.PKA435, v.EXPLANATION.String)
		var valueFromHIS sql.NullString
		err = orclDB.QueryRow(sqlTextorcl, v.YKC700, v.YZY201, v.PKA435).Scan(&valueFromHIS)
		if err != nil {
			log.Printf("oracle database scan error, because %s", err)
		}
		hisValue := string(gbk2Utf8([]byte(valueFromHIS.String)))
		log.Println(v.YKC700, v.YZY201, v.PKA435, v.EXPLANATION.String, hisValue)
		result, err := msDB.Exec(sqlText4Update, hisValue, v.YKC700, v.YZY201)
		if err != nil {
			log.Printf("sqlserver database update error, because %s", err)
		}
		fmt.Println(result.RowsAffected())
	}

}
