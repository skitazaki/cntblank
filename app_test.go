package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"testing"
)

func TestReport(t *testing.T) {
	report := new(Report)
	nullCount := report.parseRecord([]string{"a", "b", ""})
	if nullCount == 1 {
		t.Log("ok to return null count.")
	} else {
		t.Error("fail to return null count:", nullCount)
	}
	if report.records == 1 {
		t.Log("ok to parse one record.")
	} else {
		t.Error("fail to increment internal counter.")
	}
	if len(report.fields) == 3 {
		t.Log("automatically grow field length to three.")
	} else {
		t.Error("fail to grow field:", report.fields)
	}

	nullCount = report.parseRecord([]string{"$C", "", "", "-1", "123", "3.14"})
	if nullCount == 2 {
		t.Log("ok to return null count.")
	} else {
		t.Error("fail to return null count:", nullCount)
	}
	if report.records == 2 {
		t.Log("ok to parse more one record.")
	} else {
		t.Error("fail to increment internal counter.")
	}
	if len(report.fields) == 6 {
		t.Log("automatically grow field length to six.")
	} else {
		t.Error("fail to grow field:", report.fields)
	}
}

func TestReportTypeDetection(t *testing.T) {
	report := new(Report)
	report.parseRecord([]string{"1", "0xff", "T"})
	report.parseRecord([]string{"-1", "3.14", "true"})
	report.parseRecord([]string{"0", "42.0", "True"})
	report.parseRecord([]string{"123", "", "F"})
	report.parseRecord([]string{"-456", "-999.999", "false"})
	report.parseRecord([]string{"987654321", "-101.0", "False"})
	if report.records == 6 {
		t.Log("ok to parse six records.")
	} else {
		t.Error("fail to increment internal counter.")
	}
	if len(report.fields) == 3 {
		t.Log("automatically grow field length to three.")
	} else {
		t.Error("fail to grow field:", report.fields)
	}

	f := report.fields[0]
	if f.blank == 0 {
		t.Log("ok to be filled.")
	} else {
		t.Error(f)
	}
	if f.intType == 6 {
		t.Log("ok to count up integer type.")
	} else {
		t.Error(f)
	}
	if f.floatType == 0 {
		t.Log("ok to stay zero for float type.")
	} else {
		t.Error(f)
	}
	if f.boolType == 0 {
		t.Log("ok to stay zero for bool type.")
	} else {
		t.Error(f)
	}
	if f.minimum == -456 {
		t.Log("ok to calculate minimum value as -456.")
	} else {
		t.Error(f)
	}
	if f.maximum == 987654321 {
		t.Log("ok to calculate maximum value as 987654321.")
	} else {
		t.Error(f)
	}
	if f.minimumF == 0 {
		t.Log("ok not to calculate floating minimum value.")
	} else {
		t.Error(f)
	}
	if f.maximumF == 0 {
		t.Log("ok not to calculate floating maximum value.")
	} else {
		t.Error(f)
	}
	if f.trueCount == 0 {
		t.Log("ok not to calculate boolean value for true.")
	} else {
		t.Error(f)
	}
	if f.falseCount == 0 {
		t.Log("ok not to calculate boolean value for false.")
	} else {
		t.Error(f)
	}

	f = report.fields[1]
	if f.blank == 1 {
		t.Log("ok to detect one blank cell.")
	} else {
		t.Error(f)
	}
	if f.intType == 0 {
		t.Log("ok to stay zero integer type.")
	} else {
		t.Error(f)
	}
	if f.floatType == 4 {
		t.Log("ok to count up float type.")
	} else {
		t.Error(f)
	}
	if f.boolType == 0 {
		t.Log("ok to stay zero for bool type.")
	} else {
		t.Error(f)
	}
	if f.minimum == 0 {
		t.Log("ok not to calculate minimum value.")
	} else {
		t.Error(f)
	}
	if f.maximum == 0 {
		t.Log("ok not to calculate maximum value.")
	} else {
		t.Error(f)
	}
	if f.minimumF == -999.999 {
		t.Log("ok to calculate minimum value as -999.999.")
	} else {
		t.Error(f)
	}
	if f.maximumF == 42 {
		t.Log("ok to calculate maximum value as 42 although originally 42.0.")
	} else {
		t.Error(f)
	}
	if f.trueCount == 0 {
		t.Log("ok not to calculate boolean value for true.")
	} else {
		t.Error(f)
	}
	if f.falseCount == 0 {
		t.Log("ok not to calculate boolean value for false.")
	} else {
		t.Error(f)
	}

	f = report.fields[2]
	if f.intType == 0 {
		t.Log("ok to stay zero integer type.")
	} else {
		t.Error(f)
	}
	if f.floatType == 0 {
		t.Log("ok to stay zero for float type.")
	} else {
		t.Error(f)
	}
	if f.boolType == 6 {
		t.Log("ok to count up for bool type.")
	} else {
		t.Error(f)
	}
	if f.minimum == 0 {
		t.Log("ok not to calculate minimum value.")
	} else {
		t.Error(f)
	}
	if f.maximum == 0 {
		t.Log("ok not to calculate maximum value.")
	} else {
		t.Error(f)
	}
	if f.minimumF == 0 {
		t.Log("ok not to calculate floating minimum value.")
	} else {
		t.Error(f)
	}
	if f.maximumF == 0 {
		t.Log("ok not to calculate floating maximum value.")
	} else {
		t.Error(f)
	}
	if f.trueCount == 3 {
		t.Log("ok to calculate boolean value for true.")
	} else {
		t.Error(f)
	}
	if f.falseCount == 3 {
		t.Log("ok to calculate boolean value for false.")
	} else {
		t.Error(f)
	}
}

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
