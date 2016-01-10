package main

import (
	"testing"
)

func TestExcelReader(t *testing.T) {
	dialect := &FileDialect{
		Comma:     ',',
		HasHeader: true,
	}
	reader, err := OpenFile("testdata/addrcode_jp.xlsx", dialect)
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
