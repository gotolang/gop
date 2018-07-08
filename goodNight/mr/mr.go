package main

import (
	"log"
	"time"
)

type patientMR struct {
	b basic
	d []diagnosis
	o []operation
	c []childbirth
	t []tumor
}

type basic struct {
	patientID     string
	inpNO         string
	name          string
	admitDatetime time.Time
	disDatetime   time.Time
}

type diagnosis struct {
	diagType     string
	index        int
	diagDesc     string
	diagDatetime time.Time
}

type operation struct {
}

type childbirth struct {
}

type tumor struct {
}

func main() {
	log.Println("mr")
}
