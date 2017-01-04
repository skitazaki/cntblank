package main

import (
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"
)

func TestExpandDir(t *testing.T) {
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("%v", err)
	}
	rootDir, err := filepath.Abs(filepath.Join(pwd, "..", "..", ".."))
	if err != nil {
		t.Fatalf("%v", err)
	}
	dir := filepath.Join(rootDir, "testdata")
	list := []string{filepath.Join(rootDir, "Makefile"), dir}
	c := newFileCollector(true, []string{})
	err = c.CollectAll(list)
	if err != nil {
		t.Fatalf("%v", err)
	}
	files := c.files
	if len(files) != 4 {
		t.Fatalf("expected files is %d, but actual files is %d: %v", 4, len(files), files)
	}
	if path.Base(files[0].path) != "Makefile" {
		t.Errorf("first element have to be Makefile: %s", files[0].path)
	}
	for i, expected := range []struct {
		fname  string
		md5hex string
	}{
		{"addrcode_jp.xlsx", "04770cd319075f0141dd56ce7a12161a"},
		{"elementary_school_teacher_ja.csv", "4907382473dfd0d7828bc360c91c0a08"},
		{"prefecture_jp.tsv", "4884d04103df0fd8a9e792866ca0b870"},
	} {
		file := files[i+1]
		d, f := path.Split(file.path)
		if !strings.HasSuffix(d, "testdata/") {
			t.Errorf("invalid directory: %s", d)
		}
		if f != expected.fname {
			t.Errorf("invalid file name: %s, expected=%s", f, expected)
		}
		md5hex, err := file.Checksum()
		if err != nil {
			t.Errorf("%v", err)
		}
		if md5hex != expected.md5hex {
			t.Errorf("invalid MD5: %s, expected=%s", md5hex, expected)
		}
	}
}

func TestExpandDirWithExt(t *testing.T) {
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("%v", err)
	}
	rootDir, err := filepath.Abs(filepath.Join(pwd, "..", "..", ".."))
	if err != nil {
		t.Fatalf("%v", err)
	}
	dir := filepath.Join(rootDir, "testdata")
	list := []string{filepath.Join(rootDir, "Makefile"), dir}
	c := newFileCollector(true, []string{".xlsx", ".csv", ".tsv"})
	err = c.CollectAll(list)
	if err != nil {
		t.Fatalf("%v", err)
	}
	files := c.files
	if len(files) != 3 {
		t.Fatalf("expected files is %d, but actual files is %d: %v", 3, len(files), files)
	}
	for i, expected := range []string{
		"addrcode_jp.xlsx",
		"elementary_school_teacher_ja.csv",
		"prefecture_jp.tsv",
	} {
		d, f := path.Split(files[i].path)
		if !strings.HasSuffix(d, "testdata/") {
			t.Errorf("invalid directory: %s", d)
		}
		if f != expected {
			t.Errorf("invalid file name: %s, expected=%s", f, expected)
		}
	}
}
