package main

import (
	"bytes"
	"testing"

	"csvhelper"
)

func TestReportWriterWithHeader(t *testing.T) {
	buffer := &bytes.Buffer{}
	dialect := &csvhelper.FileDialect{
		HasHeader: true,
	}
	w := NewReportWriter(buffer, "", dialect)
	s := make([]Report, 1)
	s[0] = Report{}
	err := w.Write(s)
	if err != nil {
		t.Errorf("Unexpected error: %s\n", err)
	}
	out := buffer.String()
	expected := "seq,Name,#Blank,%Blank,MinLength,MaxLength,#Int,#Float,#Bool,#Time,Minimum,Maximum,#True,#False\n"
	if out != expected {
		t.Errorf("out=%q want %q", out, expected)
	}
}

func TestReportWriterWithMetadata(t *testing.T) {
	buffer := &bytes.Buffer{}
	dialect := &csvhelper.FileDialect{
		HasMetadata: true,
		HasHeader:   true,
	}
	w := NewReportWriter(buffer, "", dialect)
	s := make([]Report, 1)
	s[0] = Report{}
	err := w.Write(s)
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
	w := NewReportWriter(buffer, "", nil)
	s := make([]Report, 1)
	s[0] = Report{}
	err := w.Write(s)
	if err != nil {
		t.Errorf("Unexpected error: %s\n", err)
	}
	out := buffer.String()
	expected := ""
	if out != expected {
		t.Errorf("out=%q want %q", out, expected)
	}
}
