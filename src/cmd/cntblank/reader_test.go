package main

import (
	"os"
	"path/filepath"
	"testing"

	"csvhelper"
)

func getTestfilePath(fname string) (path string, err error) {
	pwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	projectDir, err := filepath.Abs(filepath.Join(pwd, "..", "..", "..", "testdata"))
	if err != nil {
		return "", err
	}
	return filepath.Join(projectDir, fname), nil
}

func TestReader(t *testing.T) {
	dialect := &csvhelper.FileDialect{
		HasHeader: true,
		Comma:     '\t',
	}
	path, err := getTestfilePath("prefecture_jp.tsv")
	if err != nil {
		t.Fatalf("%v", err)
	}
	reader, err := OpenFile(path, dialect)
	if err != nil {
		t.Fatalf("%v", err)
	}
	for i, expected := range []struct {
		record []string
	}{
		{[]string{"都道府県コード", "都道府県"}},
		{[]string{"01", "北海道"}},
		{[]string{"02", "青森県"}},
		{[]string{"03", "岩手県"}},
		{[]string{"04", "宮城県"}},
	} {
		r, err := reader.Read()
		if err != nil {
			t.Fatalf("Line#%d: %v", i+1, err)
		}
		if len(r) != len(expected.record) {
			t.Errorf("Line#%d should have %d elements, but %d exists", i+1, len(expected.record), len(r))
		}
		for j, v := range expected.record {
			if r[j] != v {
				t.Errorf("Line#%d Column#%d should be \"%s\", but \"%s\"", i+1, j+1, v, r[j])
			}
		}
	}
}

func TestExcelReader(t *testing.T) {
	dialect := &csvhelper.FileDialect{
		HasHeader: true,
	}
	path, err := getTestfilePath("addrcode_jp.xlsx")
	if err != nil {
		t.Fatalf("%v", err)
	}
	reader, err := OpenFile(path, dialect)
	if err != nil {
		t.Fatalf("%v", err)
	}
	for i, expected := range []struct {
		record []string
	}{
		{[]string{"団体コード", "都道府県名\n（漢字）", "市区町村名\n（漢字）", "都道府県名\n（カナ）", "市区町村名\n（カナ）", "", ""}},
		{[]string{"010006", "北海道", "", "ﾎｯｶｲﾄﾞｳ", "", "", ""}},
		{[]string{"011002", "北海道", "札幌市", "ﾎｯｶｲﾄﾞｳ", "ｻｯﾎﾟﾛｼ", "", ""}},
	} {
		r, err := reader.Read()
		if err != nil {
			t.Fatalf("Line#%d: %v", i+1, err)
		}
		if len(r) != len(expected.record) {
			t.Errorf("Line#%d should have %d elements, but %d exists", i+1, len(expected.record), len(r))
		}
		for j, v := range expected.record {
			if r[j] != v {
				t.Errorf("Line#%d Column#%d should be \"%s\", but \"%s\"", i+1, j+1, v, r[j])
			}
		}
	}
}

func TestExcelReaderSheetOption(t *testing.T) {
	dialect := &csvhelper.FileDialect{
		HasHeader:   false,
		SheetNumber: 2,
	}
	path, err := getTestfilePath("addrcode_jp.xlsx")
	if err != nil {
		t.Fatalf("%v", err)
	}
	reader, err := OpenFile(path, dialect)
	if err != nil {
		t.Fatalf("%v", err)
	}
	for i, expected := range []struct {
		record []string
	}{
		{[]string{"011002", "札幌市", "さっぽろし", "", ""}},
		{[]string{"011011", "札幌市中央区", "さっぽろしちゅうおうく", "", ""}},
	} {
		r, err := reader.Read()
		if err != nil {
			t.Fatalf("Line#%d: %v", i+1, err)
		}
		if len(r) != len(expected.record) {
			t.Errorf("Line#%d should have %d elements, but %d exists", i+1, len(expected.record), len(r))
		}
		for j, v := range expected.record {
			if r[j] != v {
				t.Errorf("Line#%d Column#%d should be \"%s\", but \"%s\"", i+1, j+1, v, r[j])
			}
		}
	}
}

func TestExcelReaderExceedSheetNumber(t *testing.T) {
	dialect := &csvhelper.FileDialect{
		SheetNumber: 10,
	}
	path, err := getTestfilePath("addrcode_jp.xlsx")
	if err != nil {
		t.Fatalf("%v", err)
	}
	_, err = OpenFile(path, dialect)
	if err == nil {
		t.Fatalf("%d should exceed acutual sheet number", dialect.SheetNumber)
	}
}
