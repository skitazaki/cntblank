package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"fmt"
	"testing"
)

func TestReport(t *testing.T) {
	report := new(Report)
	for i, tc := range []struct {
		record     []string
		nullCount  int
		numRecords int
		numFields  int
	}{
		{[]string{"a", "b", ""}, 1, 1, 3},
		{[]string{"$C", "", "", "-1", "123", "3.14"}, 2, 2, 6},
	} {
		nullCount := report.parseRecord(tc.record)
		if nullCount != tc.nullCount {
			t.Errorf("[#%d] fail to return null count: actual=%d, expected=%d", i+1, nullCount, tc.nullCount)
		}
		if report.records != tc.numRecords {
			t.Errorf("[#%d] fail to increment internal counter: actual=%d, expected=%d", i+1, report.records, tc.numRecords)
		}
		if len(report.fields) != tc.numFields {
			t.Errorf("[#%d] fail to grow field: actual=%d, expected=%d", i+1, len(report.fields), tc.numFields)
		}
	}
	for i, tc := range []struct {
		columnName string
		blankCount string
		blankRatio string
	}{
		{"Column001", "0", "0.0000"},
		{"Column002", "1", "0.5000"},
		{"Column003", "2", "1.0000"},
	} {
		f := report.fields[i]
		r := f.format(report.records)
		if len(r) != 15 {
			t.Errorf("#%d field formatter returns invalid result, which has %d elements.", i+1, len(r))
		}
		if r[0] != fmt.Sprint(i+1) {
			t.Errorf("#%d field has invalid sequence number: actual=%s, expected=%s", i+1, r[0], "1")
		}
		if r[1] != tc.columnName {
			t.Errorf("#%d field has invalid column name: actual=%s, expected=%s", i+1, r[1], tc.columnName)
		}
		if r[2] != tc.blankCount {
			t.Errorf("#%d field has invalid blank count: actual=%s, expected=%s", i+1, r[2], tc.blankCount)
		}
		if r[3] != tc.blankRatio {
			t.Errorf("#%d field has invalid blank ratio: actual=%s, expected=%s", i+1, r[3], tc.blankRatio)
		}
	}
}

func TestReportTypeDetection(t *testing.T) {
	report := new(Report)
	report.parseRecord([]string{"1", "0xff", "T"})
	report.parseRecord([]string{"-1", "3.14", "true"})
	report.parseRecord([]string{"0", "42.0", "True"})
	report.parseRecord([]string{"123", "", "F"})
	report.parseRecord([]string{"-456", "-999.999", "false"})
	report.parseRecord([]string{"987654321", "-101.0あ", "False"})
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

	for i, tc := range []struct {
		blank      int
		minimum    int
		maximum    int
		minimumF   float64
		maximumF   float64
		trueCount  int
		falseCount int
		intType    int
		floatType  int
		boolType   int
		fullWidth  int
	}{
		{0, -456, 987654321, 0, 0, 0, 0, 6, 0, 0, 0},
		{1, 0, 0, -999.999, 42, 0, 0, 0, 3, 0, 1},
		{0, 0, 0, 0, 0, 3, 3, 0, 0, 6, 0},
	} {
		f := report.fields[i]
		if f.blank != tc.blank {
			t.Errorf("#%d fail to count blank: actual=%d, expected=%d", i+1, f.blank, tc.blank)
		}
		if f.minimum != tc.minimum {
			t.Errorf("#%d fail to calculate minimum value: actual=%d, expected=%d", i+1, f.minimum, tc.minimum)
		}
		if f.maximum != tc.maximum {
			t.Errorf("#%d fail to calculate maximum value: actual=%d, expected=%d", i+1, f.maximum, tc.maximum)
		}
		if f.minimumF != tc.minimumF {
			t.Errorf("#%d fail to calculate minimum float value: actual=%f, expected=%f", i+1, f.minimumF, tc.minimumF)
		}
		if f.maximumF != tc.maximumF {
			t.Errorf("#%d fail to calculate maximum float value: actual=%f, expected=%f", i+1, f.maximumF, tc.maximumF)
		}
		if f.trueCount != tc.trueCount {
			t.Errorf("#%d fail to count boolean true: actual=%d, expected=%d", i+1, f.trueCount, tc.trueCount)
		}
		if f.falseCount != tc.falseCount {
			t.Errorf("#%d fail to count boolean false: actual=%d, expected=%d", i+1, f.falseCount, tc.falseCount)
		}
		if f.intType != tc.intType {
			t.Errorf("#%d fail to count integer type: actual=%d, expected=%d", i+1, f.intType, tc.intType)
		}
		if f.floatType != tc.floatType {
			t.Errorf("#%d fail to count float type: actual=%d, expected=%d", i+1, f.floatType, tc.floatType)
		}
		if f.boolType != tc.boolType {
			t.Errorf("#%d fail to count bool type: actual=%d, expected=%d", i+1, f.boolType, tc.boolType)
		}
		if f.fullWidth != tc.fullWidth {
			t.Errorf("#%d fail to count full width: actual=%d, expected=%d", i+1, f.fullWidth, tc.fullWidth)
		}
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

func TestTrimCell(t *testing.T) {
	input := []byte(`key, value ," comment"
PI,3.1415926535897932384,normal
PI ,3.1415926535897932384 ,"blank after value "
 PI, 3.1415926535897932384, blank before value
 PI , 3.1415926535897932384 ," blank both of value "
`)
	buffer := new(bytes.Buffer)
	w := bufio.NewWriter(buffer)
	reader := csv.NewReader(bytes.NewBuffer(input))
	app, _ := newApplication(w, "", "", false)
	report, err := app.cntblank(reader, ",", false, false)
	if err != nil {
		t.Error(err)
	}
	if report.records == 4 {
		t.Log("ok to read three records with header line.")
	} else {
		t.Error("fail to count invalid records:", report.records)
	}
	if len(report.fields) == 3 {
		t.Log("ok to read three columns.")
	} else {
		t.Error("fail to count invalid columns:", len(report.fields))
	}

	for i, tc := range []struct {
		fieldName string
		minLength int
		maxLength int
	}{
		{"key", 2, 2},
		{"value", 21, 21},
		{"comment", 6, 19},
	} {
		f := report.fields[i]
		if f.name == tc.fieldName {
			t.Logf("#%d field name is ok.", i+1)
		} else {
			t.Errorf("#%d field name is invalid: actual=\"%s\", expected=\"%s\"", i+1, f.name, tc.fieldName)
		}
		if f.minLength == tc.minLength {
			t.Logf("#%d field minimum length is ok.", i+1)
		} else {
			t.Errorf("#%d field has invalid minLength: actual=%d, expected=%d", i+1, f.minLength, tc.minLength)
		}
		if f.maxLength == tc.maxLength {
			t.Logf("#%d field maximum length is ok.", i+1)
		} else {
			t.Errorf("#%d field has invalid maxLength: actual=%d, expected=%d", i+1, f.maxLength, tc.maxLength)
		}
	}
}
