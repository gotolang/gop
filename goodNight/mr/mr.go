package main

import (
	"log"
	"time"
)

// Basic 病人基本信息 pat_master_index
type Basic struct {
	PatientID   string    // 病人ID
	InpNo       string    // 住院号
	Name        string    // 姓名
	Sex         string    // 性别
	DateOfBirth time.Time //年龄
}

// Hospitalization 病人住院信息 pat_visit
type Hospitalization struct {
	Basic
	VisitID         int
	DateOfAdmit     time.Time
	DateOfDischarge time.Time
	DeptOfAdmit     string // 入院科室
	DeptOfDischarge string // 出院科室
}

// DiagType 诊断类型字典 diagnosis_type_dict
type DiagType struct {
	No         int    // 诊断类型代码
	TypeOfDiag string // 诊断类型描述
}

// ICD 诊断编码 diagnosis_dict
type ICD struct {
	Code string // ICD10编码
	Desc string // ICD10诊断
}

// Diagnosis 诊断信息 diagnosis
type Diagnosis struct {
	No int
	DiagType
	ICD
	DateOfDiag    time.Time
	DaysOfTreat   int
	ResultOfTreat string
}

// Operation 手术信息 operation
type Operation struct {
}

// Child 新生儿信息
type Child struct {
	Basic
}

// Childbirth 分娩信息
type Childbirth struct {
	Child
	DateOfchildbirth time.Time // 分娩日期
}

// Tumor 肿瘤信息
type Tumor struct {
	Basic
}

// MR 住院病人病案信息
type MR struct {
	Hospitalization
	Diags []Diagnosis // 诊断：门急诊诊断、入院诊断、出院诊断、病理诊断、其他诊断...
	Opers []Operation // 手术：手术1、手术2...
	Cbs   []Childbirth
	Tms   []Tumor
}

func mr2json() {

}

func mr2xml() {

}

func main() {
	log.Println("mr")
}
