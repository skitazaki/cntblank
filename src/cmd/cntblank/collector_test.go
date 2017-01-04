package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExpandDir(t *testing.T) {
	pwd, err := os.Getwd()
	require.Nil(t, err)
	rootDir, err := filepath.Abs(filepath.Join(pwd, "..", "..", ".."))
	require.Nil(t, err)
	dir := filepath.Join(rootDir, "testdata")
	list := []string{filepath.Join(rootDir, "Makefile"), dir}
	c := newFileCollector(true, []string{})
	err = c.CollectAll(list)
	require.Nil(t, err)
	files := c.files
	assert.Equal(t, 4, len(files))
	assert.Equal(t, "Makefile", files[0].Name(), "first element have to be Makefile")
	for i, expected := range []struct {
		fname  string
		md5hex string
	}{
		{"addrcode_jp.xlsx", "04770cd319075f0141dd56ce7a12161a"},
		{"elementary_school_teacher_ja.csv", "4907382473dfd0d7828bc360c91c0a08"},
		{"prefecture_jp.tsv", "4884d04103df0fd8a9e792866ca0b870"},
	} {
		file := files[i+1]
		if !strings.HasSuffix(file.Dir(), "/testdata") {
			t.Errorf("invalid directory: %s", file.Dir())
		}
		assert.Equal(t, expected.fname, file.Name(), "invalid file name")
		md5hex, err := file.Checksum()
		require.Nil(t, err, "Checksum() should returns nil: %v", err)
		assert.Equal(t, expected.md5hex, md5hex, "invalid MD5")
	}
}

func TestExpandDirWithExt(t *testing.T) {
	pwd, err := os.Getwd()
	require.Nil(t, err)
	rootDir, err := filepath.Abs(filepath.Join(pwd, "..", "..", ".."))
	require.Nil(t, err)
	dir := filepath.Join(rootDir, "testdata")
	list := []string{filepath.Join(rootDir, "Makefile"), dir}
	c := newFileCollector(true, []string{".xlsx", ".csv", ".tsv"})
	err = c.CollectAll(list)
	require.Nil(t, err)
	files := c.files
	assert.Equal(t, 3, len(files))
	for i, expected := range []struct {
		fname  string
		md5hex string
	}{
		{"addrcode_jp.xlsx", "04770cd319075f0141dd56ce7a12161a"},
		{"elementary_school_teacher_ja.csv", "4907382473dfd0d7828bc360c91c0a08"},
		{"prefecture_jp.tsv", "4884d04103df0fd8a9e792866ca0b870"},
	} {
		file := files[i]
		if !strings.HasSuffix(file.Dir(), "/testdata") {
			t.Errorf("invalid directory: %s", file.Dir())
		}
		assert.Equal(t, expected.fname, file.Name(), "invalid file name")
		md5hex, err := file.Checksum()
		require.Nil(t, err, "Checksum() should returns nil: %v", err)
		assert.Equal(t, expected.md5hex, md5hex, "invalid MD5")
	}
}
