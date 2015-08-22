package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"testing"
)

func TestApplication(t *testing.T) {
	input := []byte(`key,value
A,B
C,D
`)
	buffer := new(bytes.Buffer)
	w := bufio.NewWriter(buffer)
	reader := csv.NewReader(bytes.NewBuffer(input))
	app, _ := newApplication(w, "", "", false)
	report, err := app.cntblank(reader, ",", false, false)
	if err != nil {
		t.Error(err)
	}
	if report.records == 2 {
		t.Log("ok to read two records with header line.")
	} else {
		t.Error("fail to count invalid records:", report.records)
	}
	if len(report.fields) == 2 {
		t.Log("ok to read two columns.")
	} else {
		t.Error("fail to count invalid columns:", report.fields)
	}
}

func TestMinMaxLength(t *testing.T) {
	input := []byte(`key,value
ABC,0123456789
PI,3.1415926535897932384
ネイピア数,2.718281828459045235360287471352
`)
	buffer := new(bytes.Buffer)
	w := bufio.NewWriter(buffer)
	reader := csv.NewReader(bytes.NewBuffer(input))
	app, _ := newApplication(w, "", "", false)
	report, err := app.cntblank(reader, ",", false, false)
	if err != nil {
		t.Error(err)
	}
	if report.records == 3 {
		t.Log("ok to read three records with header line.")
	} else {
		t.Error("fail to count invalid records:", report.records)
	}
	if len(report.fields) == 2 {
		t.Log("ok to read two columns.")
	} else {
		t.Error("fail to count invalid columns:", report.fields)
	}

	f := report.fields[0]
	if f.name == "key" {
		t.Log("first field name is ok.")
	} else {
		t.Error("first field name is invalid: ", f.name)
	}
	if f.minLength == 2 {
		t.Log("minimum length of first field is 2.")
	} else {
		t.Error("first field has invalid minLength:", f.minLength)
	}
	if f.maxLength == 5 {
		t.Log("maximum length of first field is 5.")
	} else {
		t.Error("first field has invalid maxLength:", f.maxLength)
	}

	f = report.fields[1]
	if f.name == "value" {
		t.Log("second field name is ok.")
	} else {
		t.Error("second field name is invalid:", f.name)
	}
	if f.minLength == 10 {
		t.Log("minimum length of second field is 10.")
	} else {
		t.Error("second field has invalid minLength:", f.minLength)
	}
	if f.maxLength == 32 {
		t.Log("maximum length of second field is 32.")
	} else {
		t.Error("second field has invalid maxLength:", f.maxLength)
	}
}
