package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/gotolang/nt"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-oci8"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// Basic 病人基本信息 pat_master_index
type Basic struct {
	PatientID   string         `db:"PATIENT_ID"`    // 病人ID
	InpNo       sql.NullString `db:"INP_NO"`        // 住院号
	Name        string         `db:"NAME"`          // 姓名
	Sex         string         `db:"SEX"`           // 性别
	DateOfBirth nt.NullTime    `db:"DATE_OF_BIRTH"` //年龄
}

// Hospitalization 病人住院信息 pat_visit
type Hospitalization struct {
	// PatientID   string `db:"PATIENT_ID"`   // 病人ID
	// InpNo       sql.NullString `db:"INP_NO"`    // 住院号
	// Name        string `db:"NAME"`    // 姓名
	// Sex         string `db:"SEX"`    // 性别
	// DateOfBirth nt.NullTime `db:"DATE_OF_BIRTH"` //年龄
	Basic
	VisitID         int            `db:"VISIT_ID"`
	DateOfAdmit     time.Time      `db:"ADMISSION_DATE_TIME"`
	DateOfDischarge nt.NullTime    `db:"DISCHARGE_DATE_TIME"`
	DeptOfAdmit     string         `db:"DEPT_ADMISSION_TO"`   // 入院科室
	DeptOfDischarge sql.NullString `db:"DEPT_DISCHARGE_FROM"` // 出院科室
}

// DiagType 诊断类型字典 diagnosis_type_dict
type DiagType struct {
	No         int    // 诊断类型代码
	TypeOfDiag string // 诊断类型描述
}

// ICD 诊断编码 diagnosis_dict
type ICD struct {
	Code sql.NullString `db:"DIAGNOSIS_CODE"` // ICD10编码
	Desc string         `db:"DIAGNOSIS_DESC"` // ICD10诊断
}

// Diagnosis 诊断信息 diagnosis
type Diagnosis struct {
	TypeOfDiag string `db:"DIAGNOSIS_TYPE"`
	No         int    `db:"DIAGNOSIS_NO"`
	ICD
	DateOfDiag    time.Time      `db:"DIAGNOSIS_DATE"`
	DaysOfTreat   sql.NullInt64  `db:"TREAT_DAYS"`
	ResultOfTreat sql.NullString `db:"TREAT_RESULT"`
}

// Operation 手术信息 operation
type Operation struct {
	OperationNo        string         `db:"OPERATION_NO"`
	OperationCode      string         `db:"OPERATION_CODE"`
	OperationDesc      string         `db:"OPERATION_DESC"`
	DateOfOperation    nt.NullTime    `db:"OPERATING_DATE"`
	Heal               sql.NullString `db:"HEAL"`
	WoundGrade         sql.NullString `db:"WOUND_GRADE"`
	AnaesthesiaMethod  sql.NullString `db:"ANAESTHESIA_METHOD"`
	Operator           sql.NullString `db:"OPERATOR"`
	FirstAssistant     sql.NullString `db:"FIRST_ASSISTANT"`
	SecondAssistant    sql.NullString `db:"SECOND_ASSISTANT"`
	Anesthetist        sql.NullString `db:"ANESTHETIST"`
	DoctorOfAnesthesia sql.NullString `db:"ANESTHESIA_DOCTOR"`
}

// // Child 新生儿信息
// type Child struct {
// 	Basic
// }

// // Childbirth 分娩信息
// type Childbirth struct {
// 	Child
// 	DateOfchildbirth time.Time // 分娩日期
// }

// // Tumor 肿瘤信息
// type Tumor struct {
// 	Basic
// }

// MR 住院病人病案信息
type MR struct {
	Hospitalization
	Diags []Diagnosis // 诊断：门急诊诊断、入院诊断、出院诊断、病理诊断、其他诊断...
	Opers []Operation // 手术：手术1、手术2...
	// Cbs   []Childbirth
	// Tms   []Tumor
}

func mr2json(mr interface{}) (string, error) {
	b, err := json.Marshal(mr)
	if err != nil {
		log.Println(err)
		return "", err
	}
	// fmt.Println(string(b))
	return string(b), nil
}

func mr2xml(mr MR) {

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

func generateMrRecords() {

	var mr MR
	var mrs []MR
	db, err := sqlx.Open("oci8", "hw/a@dbtest")
	if err != nil {
		log.Println(err)
		return
	}

	sqlText := `select x.patient_id,
					x.inp_no,
					x.name,
					x.sex,
					x.date_of_birth,
					v.visit_id,
					v.admission_date_time,
					v.discharge_date_time,
					v.dept_admission_to,
					v.dept_discharge_from
				from pat_master_index x, pat_visit v
				where x.patient_id = v.patient_id
				and x.patient_id = :pid`

	rows, err := db.Queryx(sqlText, "D00193374")
	if err != nil {
		log.Println(err)
		return
	}

	defer rows.Close()

	var patVisit Hospitalization
	// var patVisits []Hospitalization
	for rows.Next() {

		err = rows.StructScan(&patVisit)
		if err != nil {
			log.Println(err)
			return
		}

		// fmt.Println(patVisit)
		mr.Hospitalization = patVisit

		var opers []Operation
		var oper Operation

		sqlTextOperation := `
						select 
							o.operation_no,
							o.operation_code,
							o.operation_desc,
							o.operating_date,
							o.heal,
							o.wound_grade,
							o.anaesthesia_method,
							o.operator,
							o.first_assistant,
							o.second_assistant,
							o.anesthetist,
							o.anesthesia_doctor
						from operation o
						where o.patient_id = :pid and o.visit_id = :vid
						`
		rows2, err := db.Queryx(sqlTextOperation, patVisit.PatientID, patVisit.VisitID)
		if err != nil {
			log.Println(err)
		}
		defer rows2.Close()

		for rows2.Next() {

			err = rows2.StructScan(&oper)
			if err != nil {
				log.Println(err)
				return
			}
			opers = append(opers, oper)
		}
		mr.Opers = opers

		var diags []Diagnosis
		var diag Diagnosis

		sqlTextDiagnosis := `select s.diagnosis_type,
							s.diagnosis_no,
							s.diagnosis_code,
							s.diagnosis_desc,
							s.diagnosis_date,
							s.treat_days,
							s.treat_result
					from diagnosis s
					where s.patient_id = :pid
						and s.visit_id = :vid`
		rows3, err := db.Queryx(sqlTextDiagnosis, patVisit.PatientID, patVisit.VisitID)
		if err != nil {
			log.Println(err)
		}
		defer rows3.Close()

		for rows3.Next() {

			err = rows3.StructScan(&diag)
			if err != nil {
				log.Println(err)
				return
			}
			diags = append(diags, diag)
		}
		mr.Diags = diags
		mrs = append(mrs, mr)
	}

	mrjson, err := mr2json(mrs)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(mrjson)
}

func main() {
	generateMrRecords()

}
