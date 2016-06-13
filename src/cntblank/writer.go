package main

import (
	"encoding/csv"
	"fmt"
	"path/filepath"
)

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
	showMetadata bool
	showHeader   bool
	w            *csv.Writer
}

// NewReportWriter returns a new ReportWriter that writes to w.
func NewReportWriter(w *csv.Writer, showMetadata bool) *ReportWriter {
	return &ReportWriter{
		showMetadata: showMetadata,
		showHeader:   true,
		w:            w,
	}
}

func (w *ReportWriter) Write(report Report) error {
	if w.showMetadata {
		preamble := make([]string, 4)
		if len(report.path) > 0 {
			preamble[0] = "# File"
			preamble[1] = report.path
			preamble[2] = filepath.Base(report.path)
			preamble[3] = report.md5hex
			w.w.Write(preamble)
		}
		preamble[0] = "# Field"
		preamble[1] = fmt.Sprint(len(report.fields))
		if report.hasHeader {
			preamble[2] = "(has header)"
		} else {
			preamble[2] = ""
		}
		preamble[3] = ""
		w.w.Write(preamble)
		preamble[0] = "# Record"
		preamble[1] = fmt.Sprint(report.records)
		preamble[2] = ""
		w.w.Write(preamble)
	}
	// Put header line.
	if w.showHeader {
		w.w.Write(ReportOutputFields)
	}
	// Put each field report.
	for i := 0; i < len(report.fields); i++ {
		r := report.fields[i]
		w.w.Write(r.format(report.records))
	}
	w.w.Flush()
	return w.w.Error()
}
