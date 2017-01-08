package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"

	"csvhelper"
)

func TestNewReportWriter(t *testing.T) {
	a := assert.New(t)
	w := NewReportWriter(nil, 0, nil)
	a.Nil(w)
	w = NewReportWriter(nil, CSV, nil)
	a.NotNil(w)
	w = NewReportWriter(nil, JSON, nil)
	a.NotNil(w)
	w = NewReportWriter(nil, HTML, nil)
	a.NotNil(w)
}

func TestReportWriterWithHeader(t *testing.T) {
	buffer := &bytes.Buffer{}
	dialect, err := csvhelper.NewFileDialect("", "", true)
	w := NewReportWriter(buffer, CSV, dialect)
	s := make([]Report, 1)
	s[0] = Report{}
	err = w.Write(s)
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
	dialect, err := csvhelper.NewFileDialect("", "", true)
	dialect.HasMetadata = true
	w := NewReportWriter(buffer, CSV, dialect)
	s := make([]Report, 1)
	s[0] = Report{}
	err = w.Write(s)
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
	w := NewReportWriter(buffer, CSV, nil)
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

func TestReportWriter_JSON(t *testing.T) {
	expected := `[{"header":false,"records":0,"fields":null}]`
	a := assert.New(t)
	buffer := &bytes.Buffer{}
	w := NewReportWriter(buffer, JSON, nil)
	s := make([]Report, 1)
	s[0] = Report{}
	err := w.Write(s)
	a.Nil(err)
	a.Equal(expected, buffer.String())
}
