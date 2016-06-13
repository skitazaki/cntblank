package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"

	log "github.com/Sirupsen/logrus"
)

// Format is output format
type Format int

const (
	// Csv is delimiter separated values
	Csv = iota + 1
	// Excel is Microsoft office suite
	Excel
	// JSON is JSON object
	JSON
	// Text uses template
	Text
	// HTML uses template
	HTML
)

func (s Format) String() string {
    switch s {
    case Csv:
        return "CSV"
    case Excel:
        return "Excel"
    case JSON:
        return "JSON"
    case Text:
        return "Text"
    case HTML:
        return "HTML"
    default:
        return "Unknown"
    }
}
// ReportOutputFields is a header line to write.
var ReportOutputFields = []string{
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
}

// ReportWriter is a writer object to wrap csv writer.
type ReportWriter struct {
	dialect *FileDialect
	w       io.Writer
	format  Format
}

// NewReportWriter returns a new ReportWriter that writes to w.
func NewReportWriter(w io.Writer, format string, dialect *FileDialect) *ReportWriter {
	if dialect == nil {
		dialect = &FileDialect{}
	}
	var f Format
	switch strings.ToLower(format) {
	case "csv":
		f = Csv
	case "excel":
		f = Excel
	case "json":
		f = JSON
	case "text":
		f = Text
	case "html":
		f = HTML
	default:
		f = Csv
	}
	return &ReportWriter{
		dialect: dialect,
		w:       w,
		format:  f,
	}
}

func (w *ReportWriter) Write(report Report) error {
	switch w.format {
	case Csv:
		return w.writeCsv(report)
	}
	log.Errorf("not implemented yet: format=%v", w.format)
	return nil
}

func (w *ReportWriter) writeCsv(report Report) error {
	wr := w.w
	if w.dialect.Encoding == "sjis" {
		log.Info("use ShiftJIS encoder for output.")
		encoder := japanese.ShiftJIS.NewEncoder()
		wr = transform.NewWriter(wr, encoder)
	}
	writer := csv.NewWriter(wr)
	if w.dialect.Comma != 0 {
		writer.Comma = w.dialect.Comma
	}

	if w.dialect.HasMetadata {
		preamble := make([]string, 4)
		if len(report.path) > 0 {
			preamble[0] = "# File"
			preamble[1] = report.path
			preamble[2] = filepath.Base(report.path)
			preamble[3] = report.md5hex
			writer.Write(preamble)
		}
		preamble[0] = "# Field"
		preamble[1] = fmt.Sprint(len(report.fields))
		if report.hasHeader {
			preamble[2] = "(has header)"
		} else {
			preamble[2] = ""
		}
		preamble[3] = ""
		writer.Write(preamble)
		preamble[0] = "# Record"
		preamble[1] = fmt.Sprint(report.records)
		preamble[2] = ""
		writer.Write(preamble)
	}
	// Put header line.
	if w.dialect.HasHeader {
		writer.Write(ReportOutputFields)
	}
	// Put each field report.
	for i := 0; i < len(report.fields); i++ {
		r := report.fields[i]
		writer.Write(r.format(report.records))
	}
	writer.Flush()
	return writer.Error()
}
