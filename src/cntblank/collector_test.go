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
	rootDir, err := filepath.Abs(filepath.Join(pwd, "..", ".."))
	if err != nil {
		t.Fatalf("%v", err)
	}
	dir := filepath.Join(rootDir, "testdata")
	list := []string{filepath.Join(rootDir, "Makefile"), dir}
	c, err := newFileCollector()
	if err != nil {
		t.Fatalf("%v", err)
	}
	err = c.CollectAll(list)
	if err != nil {
		t.Fatalf("%v", err)
	}
	files := c.files
	if len(files) == 3 {
		t.Log(files)
	} else {
		t.Fatalf("expected files is %d, but actual files is %d", 3, len(files))
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
