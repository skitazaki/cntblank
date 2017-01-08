package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"time"

	log "github.com/Sirupsen/logrus"

	"csvhelper"
)

// ReportWriter is a writer object to wrap csv writer.
type ReportWriter struct {
	dialect *csvhelper.FileDialect
	w       io.Writer
	format  Format
}

// NewReportWriter returns a new ReportWriter that writes to w.
func NewReportWriter(w io.Writer, format string, dialect *csvhelper.FileDialect) *ReportWriter {
	if dialect == nil {
		dialect = &csvhelper.FileDialect{}
	}
	f := CSV // default format is CSV
	if format != "" {
		f = formatFrom(format)
	}
	return &ReportWriter{
		dialect: dialect,
		w:       w,
		format:  f,
	}
}

func (w *ReportWriter) Write(reports []Report) error {
	switch w.format {
	case CSV:
		return w.writeCsv(reports)
	case JSON:
		return w.writeJSON(reports)
	case HTML:
		return w.writeHTML(reports)
	}
	log.Errorf("not implemented yet: format=%v", w.format)
	return nil
}

func (w *ReportWriter) writeCsv(reports []Report) error {
	writer := csvhelper.NewCsvWriter(w.w, w.dialect)
	for i, report := range reports {
		if i > 0 {
			writer.Write(nil)
		}
		log.Debugf("[%d] write csv file", i+1)
		w.writeCsvOne(writer, report)
	}
	writer.Flush()
	return writer.Error()
}

func (w *ReportWriter) writeCsvOne(writer *csv.Writer, report Report) error {
	if w.dialect.HasMetadata {
		preamble := make([]string, 4)
		if len(report.Path) > 0 {
			preamble[0] = "# File"
			preamble[1] = report.Path
			preamble[2] = report.Filename
			preamble[3] = report.MD5hex
			writer.Write(preamble)
		}
		preamble[0] = "# Field"
		preamble[1] = fmt.Sprint(len(report.Fields))
		if report.HasHeader {
			preamble[2] = "(has header)"
		} else {
			preamble[2] = ""
		}
		preamble[3] = ""
		writer.Write(preamble)
		preamble[0] = "# Record"
		preamble[1] = fmt.Sprint(report.Records)
		preamble[2] = ""
		writer.Write(preamble)
	}
	// Put header line.
	if w.dialect.HasHeader {
		writer.Write(ReportField{}.header())
	}
	// Put each field report.
	for i, f := range report.Fields {
		r := f.format(report.Records)
		r[0] = fmt.Sprint(i + 1)
		writer.Write(r)
	}
	return writer.Error()
}

func (w *ReportWriter) writeJSON(reports []Report) error {
	b, err := json.Marshal(reports)
	if err != nil {
		return err
	}
	w.w.Write(b)
	return nil
}

func (w *ReportWriter) writeHTML(reports []Report) error {
	path := "templates/index.html"
	b, err := Asset(path)
	if err != nil {
		// Asset was not found.
		return err
	}
	fmap := template.FuncMap{
		"deref": func(data interface{}) string {
			switch vv := data.(type) {
			case *string:
				return fmt.Sprint(*vv)
			case *int:
				if vv == nil {
					return ""
				}
				return RenderInteger("", *vv)
			case *float64:
				if vv == nil {
					return ""
				}
				return RenderFloat("", *vv)
			case *time.Time:
				if vv == nil {
					return ""
				}
				return (*vv).Format("2006-01-02 15:04:05")
			default:
				return fmt.Sprint(vv)
			}
		},
		"plus1": func(i int) int {
			return i + 1
		},
		"renderInt": func(i int) string {
			return RenderInteger("#,###.", i)
		},
	}
	tmpl, err := template.New("name").Funcs(fmap).Parse(fmt.Sprintf("%s", b))
	if err != nil {
		return err
	}
	return tmpl.Execute(w.w, reports)
}
