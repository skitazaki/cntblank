package main

import (
	"fmt"
	"testing"
	"time"
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
		if len(r) != 14 {
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
	for i, s := range [][]string{
		{"1", "0xff", "T", "2015/01/23"},
		{"-1", "3.14", "true", "2015/1/23"},
		{"0", "42.0", "True", "2015/1/2 3:45"},
		{"123", "", "F", "2015-01-02 03:04"},
		{"-456", "-999.999", "false", "2015-01-23"},
		{"987654321", "-101.0あ", "False", "20150123"},
	} {
		report.parseRecord(s)
		t.Logf("#%d record parsed, detected %d time types.", i+1, report.fields[3].timeType)
	}
	if report.records != 6 {
		t.Error("fail to increment internal counter.")
	}
	if len(report.fields) != 4 {
		t.Error("fail to grow field:", report.fields)
	}

	for i, tc := range []struct {
		blank      int
		minimum    int
		maximum    int
		minimumF   float64
		maximumF   float64
		minimumT   string
		maximumT   string
		trueCount  int
		falseCount int
		intType    int
		floatType  int
		boolType   int
		timeType   int
		fullWidth  int
	}{
		{0, -456, 987654321, -456.0, 987654321.0, "", "", 1, 1, 6, 6, 2, 0, 0},
		{1, 0, 0, -999.999, 42, "", "", 0, 0, 0, 3, 0, 0, 1},
		{0, 0, 0, 0, 0, "", "", 3, 3, 0, 0, 6, 0, 0},
		{0, 20150123, 20150123, 20150123.0, 20150123.0, "2015-01-02 03:04", "2015-01-23 00:00", 0, 0, 1, 1, 0, 6, 0},
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
		if tc.minimumT != "" {
			tt, err := time.Parse("2006-01-02 15:4", tc.minimumT)
			if err != nil {
				t.Fatal(err)
			}
			if f.minimumT != tt {
				t.Errorf("#%d fail to calculate minimum time value: actual=%v, expected=%v", i+1, f.minimumT, tt)
			}
		}
		if tc.maximumT != "" {
			tt, err := time.Parse("2006-01-02 15:4", tc.maximumT)
			if err != nil {
				t.Fatal(err)
			}
			if f.maximumT != tt {
				t.Errorf("#%d fail to calculate maximum time value: actual=%v, expected=%v", i+1, f.maximumT, tt)
			}
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
		if f.timeType != tc.timeType {
			t.Errorf("#%d fail to count time type: actual=%d, expected=%d", i+1, f.timeType, tc.timeType)
		}
		if f.fullWidth != tc.fullWidth {
			t.Errorf("#%d fail to count full width: actual=%d, expected=%d", i+1, f.fullWidth, tc.fullWidth)
		}
	}
}

func TestReportFieldFormat(t *testing.T) {
	field := new(ReportField)
	field.seq = 1
	field.name = "Column001"
	field.blank = 100
	field.timeType = 2
	field.minimumT, _ = time.Parse("2006-01-02", "2015-10-29")
	field.maximumT, _ = time.Parse("2006-01-02", "2015-11-05")

	expected := []string{
		"1", "Column001", // seq, Name
		"100", "1.0000", // #Blank, %Blank
		"", "", // MinLength, MaxLength
		"", "", "", "2", // #Int, #Float, #Bool, #Time
		"2015-10-29 00:00:00", "2015-11-05 00:00:00", // Minimum, Maximum
		"", "", // #True, #False
	}
	f := field.format(100)
	if len(f) != len(expected) {
		t.Errorf("field formatter returned invalid length: actual=%d, expected=%d", len(f), len(expected))
	}
	for i, v := range expected {
		if f[i] != v {
			t.Errorf("#%d field is invalid string: actual=%s, expected=%s", i, f[i], v)
		}
	}
}
