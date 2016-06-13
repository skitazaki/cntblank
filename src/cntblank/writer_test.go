package main

import (
	"bytes"
	"encoding/csv"
	"testing"
)

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
