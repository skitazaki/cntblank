package main

import (
	"bytes"
	"encoding/csv"
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
		{"Column004", "1", "0.5000"},
		{"Column005", "1", "0.5000"},
		{"Column006", "1", "0.5000"},
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

var formatTests = []struct {
	Sequence int
	Name     string
	Blank    int
	Length   [2]int     // Order by min, max
	Types    [4]int     // Order by int, float, bool, time
	Integer  [2]int     // Order by min, max
	Float    [2]float64 // Order by min, max
	Time     [2]string  // Order by min, max
	Expected []string
}{
	{
		Sequence: 1, Name: "Column001", Blank: 98, Length: [2]int{10, 10},
		Types: [4]int{0, 0, 0, 2},
		Time:  [2]string{"2015-10-29", "2015-11-05"},
		Expected: []string{
			"1", "Column001", // seq, Name
			"98", "0.9800", // #Blank, %Blank
			"10", "10", // MinLength, MaxLength
			"", "", "", "2", // #Int, #Float, #Bool, #Time
			"2015-10-29 00:00:00", "2015-11-05 00:00:00", // Minimum, Maximum
			"", "", // #True, #False
		},
	},
	{
		Sequence: 2, Name: "Column002", Blank: 50, Length: [2]int{3, 12},
		Types: [4]int{0, 50, 0, 0},
		Float: [2]float64{1.1, 2.2},
		Expected: []string{
			"2", "Column002", // seq, Name
			"50", "0.5000", // #Blank, %Blank
			"3", "12", // MinLength, MaxLength
			"", "50", "", "", // #Int, #Float, #Bool, #Time
			"1.1000", "2.2000", // Minimum, Maximum
			"", "", // #True, #False
		},
	},
}

func TestReportFieldFormat(t *testing.T) {
	for n, tt := range formatTests {
		field := new(ReportField)
		field.seq = tt.Sequence
		field.name = tt.Name
		field.blank = tt.Blank
		field.minLength = tt.Length[0]
		field.maxLength = tt.Length[1]
		field.intType = tt.Types[0]
		field.floatType = tt.Types[1]
		field.boolType = tt.Types[2]
		field.timeType = tt.Types[3]
		if len(tt.Integer) > 0 {
			field.minimum = tt.Integer[0]
			field.maximum = tt.Integer[1]
		}
		if len(tt.Float) > 0 {
			field.minimumF = tt.Float[0]
			field.maximumF = tt.Float[1]
		}
		if len(tt.Time) > 0 {
			field.minimumT, _ = time.Parse("2006-01-02", tt.Time[0])
			field.maximumT, _ = time.Parse("2006-01-02", tt.Time[1])
		}
		expected := tt.Expected
		f := field.format(100)
		if len(f) != len(expected) {
			t.Errorf("[#%d] field formatter returned invalid length: actual=%d, expected=%d", n, len(f), len(expected))
		}
		for i, v := range expected {
			if f[i] != v {
				t.Errorf("[#%d] #%d field is invalid string: actual=%s, expected=%s", n, i, f[i], v)
			}
		}
	}
}

func TestReportWriterWithMetadata(t *testing.T) {
	buffer := &bytes.Buffer{}
	w := NewReportWriter(csv.NewWriter(buffer), true)
	r := new(Report)
	err := w.Write(*r)
	if err != nil {
		t.Errorf("Unexpected error: %s\n", err)
	}
	out := buffer.String()
	expected := "# Field,0,,\n"
	expected += "# Record,0,,\n"
	expected += "seq,Name,#Blank,%Blank,MinLength,MaxLength,#Int,#Float,#Bool,#Time,Minimum,Maximum,#True,#False\n"
	if out != expected {
		t.Errorf("out=%q want %q", out, expected)
	}
}

func TestReportWriterWithoutMetadata(t *testing.T) {
	buffer := &bytes.Buffer{}
	w := NewReportWriter(csv.NewWriter(buffer), false)
	r := new(Report)
	err := w.Write(*r)
	if err != nil {
		t.Errorf("Unexpected error: %s\n", err)
	}
	out := buffer.String()
	expected := "seq,Name,#Blank,%Blank,MinLength,MaxLength,#Int,#Float,#Bool,#Time,Minimum,Maximum,#True,#False\n"
	if out != expected {
		t.Errorf("out=%q want %q", out, expected)
	}
}
