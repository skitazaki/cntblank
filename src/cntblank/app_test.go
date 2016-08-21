package main

import (
	"bytes"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"
)

func TestApplication(t *testing.T) {
	input := []byte(`key,value
A,B
C,D
`)
	buffer := &bytes.Buffer{}
	app, _ := newApplication(buffer, "", &FileDialect{})
	dialect := &FileDialect{
		Comma:     ',',
		HasHeader: true,
	}
	reader, err := NewReader(bytes.NewBuffer(input), dialect)
	report, err := app.cntblank(reader, dialect.HasHeader)
	if err != nil {
		t.Error(err)
	}
	if report.Records == 2 {
		t.Log("ok to read two records with header line.")
	} else {
		t.Error("fail to count invalid records:", report.Records)
	}
	if len(report.Fields) == 2 {
		t.Log("ok to read two columns.")
	} else {
		t.Error("fail to count invalid columns:", report.Fields)
	}
}

func TestMinMaxLength(t *testing.T) {
	input := []byte(`key,value
ABC,0123456789
PI,3.1415926535897932384
ネイピア数,2.718281828459045235360287471352
`)
	buffer := &bytes.Buffer{}
	app, _ := newApplication(buffer, "", &FileDialect{})
	dialect := &FileDialect{
		Comma:     ',',
		HasHeader: true,
	}
	reader, err := NewReader(bytes.NewBuffer(input), dialect)
	report, err := app.cntblank(reader, dialect.HasHeader)
	if err != nil {
		t.Error(err)
	}
	if report.Records == 3 {
		t.Log("ok to read three records with header line.")
	} else {
		t.Error("fail to count invalid records:", report.Records)
	}
	if len(report.Fields) == 2 {
		t.Log("ok to read two columns.")
	} else {
		t.Error("fail to count invalid columns:", report.Fields)
	}

	f := report.Fields[0]
	if f.Name == "key" {
		t.Log("first field name is ok.")
	} else {
		t.Error("first field name is invalid: ", f.Name)
	}
	if f.MinLength == 2 {
		t.Log("minimum length of first field is 2.")
	} else {
		t.Error("first field has invalid minLength:", f.MinLength)
	}
	if f.MaxLength == 5 {
		t.Log("maximum length of first field is 5.")
	} else {
		t.Error("first field has invalid maxLength:", f.MaxLength)
	}

	f = report.Fields[1]
	if f.Name == "value" {
		t.Log("second field name is ok.")
	} else {
		t.Error("second field name is invalid:", f.Name)
	}
	if f.MinLength == 10 {
		t.Log("minimum length of second field is 10.")
	} else {
		t.Error("second field has invalid minLength:", f.MinLength)
	}
	if f.MaxLength == 32 {
		t.Log("maximum length of second field is 32.")
	} else {
		t.Error("second field has invalid maxLength:", f.MaxLength)
	}
}

func TestTrimCell(t *testing.T) {
	input := []byte(`key, value ," comment"
PI,3.1415926535897932384,normal
PI ,3.1415926535897932384 ,"blank after value "
 PI, 3.1415926535897932384, blank before value
 PI , 3.1415926535897932384 ," blank both of value "
`)
	buffer := &bytes.Buffer{}
	app, _ := newApplication(buffer, "", &FileDialect{})
	dialect := &FileDialect{
		Comma:     ',',
		HasHeader: true,
	}
	reader, err := NewReader(bytes.NewBuffer(input), dialect)
	report, err := app.cntblank(reader, dialect.HasHeader)
	if err != nil {
		t.Error(err)
	}
	if report.Records == 4 {
		t.Log("ok to read three records with header line.")
	} else {
		t.Error("fail to count invalid records:", report.Records)
	}
	if len(report.Fields) == 3 {
		t.Log("ok to read three columns.")
	} else {
		t.Error("fail to count invalid columns:", len(report.Fields))
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
		f := report.Fields[i]
		if f.Name == tc.fieldName {
			t.Logf("#%d field name is ok.", i+1)
		} else {
			t.Errorf("#%d field name is invalid: actual=\"%s\", expected=\"%s\"", i+1, f.Name, tc.fieldName)
		}
		if f.MinLength == tc.minLength {
			t.Logf("#%d field minimum length is ok.", i+1)
		} else {
			t.Errorf("#%d field has invalid minLength: actual=%d, expected=%d", i+1, f.MinLength, tc.minLength)
		}
		if f.MaxLength == tc.maxLength {
			t.Logf("#%d field maximum length is ok.", i+1)
		} else {
			t.Errorf("#%d field has invalid maxLength: actual=%d, expected=%d", i+1, f.MaxLength, tc.maxLength)
		}
	}
}

func TestExpandDir(t *testing.T) {
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("%v", err)
	}
	rootDir, err := filepath.Abs(filepath.Join(pwd, "..", ".."))
	if err != nil {
		t.Fatalf("%v", err)
	}
	dir := filepath.Join(rootDir, "testdata")
	list := []string{filepath.Join(rootDir, "Makefile"), dir}
	files := expandDir(list)
	if len(files) == 4 {
		t.Log(files)
	} else {
		t.Errorf("expected files is %d, but actual files is %d", 4, len(files))
	}
	if path.Base(files[0]) != "Makefile" {
		t.Errorf("first element is %s", files[0])
	}
	for i, expected := range []string{
		"addrcode_jp.xlsx",
		"elementary_school_teacher_ja.csv",
		"prefecture_jp.tsv",
	} {
		d, f := path.Split(files[i+1])
		if !strings.HasSuffix(d, "testdata/") {
			t.Errorf("invalid directory: %s", d)
		}
		if f != expected {
			t.Errorf("invalid file name: %s, expected=%s", f, expected)
		}
	}
}
