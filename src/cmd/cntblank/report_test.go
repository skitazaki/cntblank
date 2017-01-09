package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewReport_Empty(t *testing.T) {
	a := assert.New(t)
	r := newReport(File{})
	a.Empty(r.Path, "Path sould be empty")
	a.Empty(r.Filename, "Filename should be empty")
	a.Empty(r.MD5hex, "MD5hex should be empty")
	a.Empty(r.HasHeader, "HasHeader should be false")
	a.Empty(r.Records, "Records should be zero")
	a.Empty(r.Fields, "Fields should be empty")
}

func TestNewReport_Filename(t *testing.T) {
	a := assert.New(t)
	r := newReport(File{path: "/path/to/file"})
	a.Equal("/path/to/file", r.Path, "Path is different")
	a.Equal("file", r.Filename, "Filename is different")
	a.Empty(r.MD5hex, "MD5hex should be empty")
	a.Empty(r.HasHeader, "HasHeader should be false")
	a.Empty(r.Records, "Records should be zero")
	a.Empty(r.Fields, "Fields should be empty")
}

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
		if report.Records != tc.numRecords {
			t.Errorf("[#%d] fail to increment internal counter: actual=%d, expected=%d", i+1, report.Records, tc.numRecords)
		}
		if len(report.Fields) != tc.numFields {
			t.Errorf("[#%d] fail to grow field: actual=%d, expected=%d", i+1, len(report.Fields), tc.numFields)
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
		f := report.Fields[i]
		r := f.format(report.Records)
		if len(r) != 14 {
			t.Errorf("#%d field formatter returns invalid result, which has %d elements.", i+1, len(r))
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
		t.Logf("#%d record parsed, detected %d time types.", i+1, report.Fields[3].TypeTime)
	}
	if report.Records != 6 {
		t.Error("fail to increment internal counter.")
	}
	if len(report.Fields) != 4 {
		t.Error("fail to grow field:", report.Fields)
	}

	for i, tc := range []struct {
		blank      int
		minimum    float64
		maximum    float64
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
		{0, -456.0, 987654321.0, "", "", 1, 1, 6, 6, 2, 0, 0},
		{1, -999.999, 42, "", "", 0, 0, 0, 3, 0, 0, 1},
		{0, 0, 0, "", "", 3, 3, 0, 0, 6, 0, 0},
		{0, 20150123.0, 20150123.0, "2015-01-02 03:04", "2015-01-23 00:00", 0, 0, 1, 1, 0, 6, 0},
	} {
		f := report.Fields[i]
		if f.Blank != tc.blank {
			t.Errorf("#%d fail to count blank: actual=%d, expected=%d", i+1, f.Blank, tc.blank)
		}
		if f.TypeFloat > 0 && *f.Minimum != tc.minimum {
			t.Errorf("#%d fail to calculate minimum value: actual=%f, expected=%f", i+1, *f.Minimum, tc.minimum)
		}
		if f.TypeFloat > 0 && *f.Maximum != tc.maximum {
			t.Errorf("#%d fail to calculate maximum value: actual=%f, expected=%f", i+1, *f.Maximum, tc.maximum)
		}
		if tc.minimumT != "" {
			tt, err := time.Parse("2006-01-02 15:4", tc.minimumT)
			if err != nil {
				t.Fatal(err)
			}
			if *f.MinTime != tt {
				t.Errorf("#%d fail to calculate minimum time value: actual=%v, expected=%v", i+1, f.MinTime, tt)
			}
		}
		if tc.maximumT != "" {
			tt, err := time.Parse("2006-01-02 15:4", tc.maximumT)
			if err != nil {
				t.Fatal(err)
			}
			if *f.MaxTime != tt {
				t.Errorf("#%d fail to calculate maximum time value: actual=%v, expected=%v", i+1, f.MaxTime, tt)
			}
		}
		if f.TypeBool > 0 && *f.BoolTrue != tc.trueCount {
			t.Errorf("#%d fail to count boolean true: actual=%d, expected=%d", i+1, *f.BoolTrue, tc.trueCount)
		}
		if f.TypeBool > 0 && *f.BoolFalse != tc.falseCount {
			t.Errorf("#%d fail to count boolean false: actual=%d, expected=%d", i+1, *f.BoolFalse, tc.falseCount)
		}
		if f.TypeInt != tc.intType {
			t.Errorf("#%d fail to count integer type: actual=%d, expected=%d", i+1, f.TypeInt, tc.intType)
		}
		if f.TypeFloat != tc.floatType {
			t.Errorf("#%d fail to count float type: actual=%d, expected=%d", i+1, f.TypeFloat, tc.floatType)
		}
		if f.TypeBool != tc.boolType {
			t.Errorf("#%d fail to count bool type: actual=%d, expected=%d", i+1, f.TypeBool, tc.boolType)
		}
		if f.TypeTime != tc.timeType {
			t.Errorf("#%d fail to count time type: actual=%d, expected=%d", i+1, f.TypeTime, tc.timeType)
		}
		if f.fullWidth != tc.fullWidth {
			t.Errorf("#%d fail to count full width: actual=%d, expected=%d", i+1, f.fullWidth, tc.fullWidth)
		}
	}
}

var formatTests = []struct {
	Name     string
	Blank    int
	Length   [2]int     // Order by min, max
	Types    [4]int     // Order by int, float, bool, time
	Float    [2]float64 // Order by min, max
	Time     [2]string  // Order by min, max
	Expected []string
}{
	{
		Name: "Column001", Blank: 98, Length: [2]int{10, 10},
		Types: [4]int{0, 0, 0, 2},
		Time:  [2]string{"2015-10-29", "2015-11-05"},
		Expected: []string{
			"", "Column001", // seq, Name
			"98", "0.9800", // #Blank, %Blank
			"10", "10", // MinLength, MaxLength
			"", "", "", "2", // #Int, #Float, #Bool, #Time
			"2015-10-29 00:00:00", "2015-11-05 00:00:00", // Minimum, Maximum
			"", "", // #True, #False
		},
	},
	{
		Name: "Column002", Blank: 50, Length: [2]int{3, 12},
		Types: [4]int{0, 50, 0, 0},
		Float: [2]float64{1.1, 2.2},
		Expected: []string{
			"", "Column002", // seq, Name
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
		field.Name = tt.Name
		field.Blank = tt.Blank
		field.MinLength = tt.Length[0]
		field.MaxLength = tt.Length[1]
		field.TypeInt = tt.Types[0]
		field.TypeFloat = tt.Types[1]
		field.TypeBool = tt.Types[2]
		field.TypeTime = tt.Types[3]
		if len(tt.Float) > 0 {
			field.Minimum = &tt.Float[0]
			field.Maximum = &tt.Float[1]
		}
		if len(tt.Time) > 0 {
			t1, _ := time.Parse("2006-01-02", tt.Time[0])
			t2, _ := time.Parse("2006-01-02", tt.Time[1])
			field.MinTime = &t1
			field.MaxTime = &t2
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

func TestReportFieldHeader(t *testing.T) {
	a := assert.New(t)
	header := ReportField{}.header()
	a.Equal(14, len(header))
	for i, s := range []string{
		"seq",
		"Name",
		"#Blank",
		"%Blank",
		"MinLength",
		"MaxLength",
		"#Int",
		"#Float",
		"#Bool",
		"#Time",
		"Minimum",
		"Maximum",
		"#True",
		"#False",
	} {
		a.Equal(s, header[i], "differ header[%d] index element", i)
	}
}
