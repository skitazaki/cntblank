package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/tealeg/xlsx"

	"csvhelper"
)

// ReportWriter writes multiple reports as specified file format.
type ReportWriter interface {
	Write([]Report) error
}

// ReportCSVWriter is a writer object to wrap csv writer.
type ReportCSVWriter struct {
	dialect *csvhelper.FileDialect
	w       io.Writer
}

// ReportJSONWriter is a writer object to write report as JSON format.
type ReportJSONWriter struct {
	w io.Writer
}

// ReportHTMLWriter is a writer object to write report as HTML format.
type ReportHTMLWriter struct {
	w io.Writer
}

// ReportExcelWriter is a writer object to write report as Excel format.
type ReportExcelWriter struct {
	w       io.Writer
	dialect *csvhelper.FileDialect
}

// NewReportWriter returns a new ReportWriter that writes to w by given format.
func NewReportWriter(w io.Writer, format Format, dialect *csvhelper.FileDialect) ReportWriter {
	switch format {
	case CSV:
		if dialect == nil {
			dialect = &csvhelper.FileDialect{}
		}
		return &ReportCSVWriter{
			dialect: dialect,
			w:       w,
		}
	case Excel:
		if dialect == nil {
			dialect = &csvhelper.FileDialect{}
		}
		return &ReportExcelWriter{
			dialect: dialect,
			w:       w,
		}
	case JSON:
		return &ReportJSONWriter{
			w: w,
		}
	case HTML:
		return &ReportHTMLWriter{
			w: w,
		}
	}
	log.Errorf("NewReportWriter: not implemented format %q", format)
	return nil
}

func (w *ReportCSVWriter) Write(reports []Report) error {
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

func (w *ReportCSVWriter) writeCsvOne(writer *csv.Writer, report Report) error {
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

func (w *ReportJSONWriter) Write(reports []Report) error {
	b, err := json.Marshal(reports)
	if err != nil {
		return err
	}
	w.w.Write(b)
	return nil
}

func (w *ReportHTMLWriter) Write(reports []Report) error {
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

func (w *ReportExcelWriter) Write(reports []Report) error {
	file := xlsx.NewFile()
	sheetFiles, err := file.AddSheet("Files")
	if err != nil {
		return err
	}
	sheetFields, err := file.AddSheet("Fields")
	if err != nil {
		return err
	}
	var row *xlsx.Row
	// Put header line on Files sheet.
	row = sheetFiles.AddRow()
	for _, k := range []string{
		"No.",
		"Path",
		"File name",
		"MD5 Checksum",
		"Has header",
		"#Fields",
		"#Records",
	} {
		w.addString(row, k)
	}
	// Put header line on Fields sheet.
	row = sheetFields.AddRow()
	for _, k := range []string{
		"No.",
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
		"MinTime",
		"MaxTime",
		"#True",
		"#False",
	} {
		w.addString(row, k)
	}
	for i, report := range reports {
		log.Debugf("[%d] write report of %q (%s)", i+1, report.Path, report.MD5hex)
		row = sheetFiles.AddRow()
		w.addInt(row, i+1)
		w.addString(row, report.Path)
		w.addString(row, report.Filename)
		w.addString(row, report.MD5hex)
		w.addBool(row, report.HasHeader)
		w.addInt(row, len(report.Fields))
		w.addInt(row, report.Records)
		if i > 0 {
			// Append blank row to separate files
			row = sheetFields.AddRow()
			w.addString(row, "")
		}
		// Put preamble about meta-data
		row = sheetFields.AddRow()
		w.addString(row, fmt.Sprintf("# File No.%d", i+1))
		w.addString(row, report.Filename)
		w.addString(row, report.MD5hex)
		row = sheetFields.AddRow()
		w.addString(row, "# Contents")
		if report.HasHeader {
			w.addString(row, "has header")
		} else {
			w.addString(row, "does not have header")
		}
		w.addString(row, "")
		w.addInt(row, len(report.Fields))
		w.addString(row, "fields")
		w.addString(row, "")
		w.addInt(row, report.Records)
		w.addString(row, "records")
		w.writeFields(sheetFields, report.Fields, report.Records)
	}
	return file.Write(w.w)
}

func (w *ReportExcelWriter) writeFields(sheet *xlsx.Sheet, fields []*ReportField, total int) {
	var row *xlsx.Row
	for i, field := range fields {
		row = sheet.AddRow()
		w.addInt(row, i+1)
		w.addString(row, field.Name)
		w.addInt(row, field.Blank)
		if total > 0 {
			w.addFloat(row, float64(field.Blank)/float64(total))
		} else {
			w.addString(row, "N/A divided by 0")
		}
		w.addInt(row, field.MinLength)
		w.addInt(row, field.MaxLength)
		w.addInt(row, field.TypeInt)
		w.addInt(row, field.TypeFloat)
		w.addInt(row, field.TypeBool)
		w.addInt(row, field.TypeTime)
		if field.Minimum != nil {
			w.addFloat(row, *field.Minimum)
		} else {
			w.addString(row, "")
		}
		if field.Maximum != nil {
			w.addFloat(row, *field.Maximum)
		} else {
			w.addString(row, "")
		}
		if field.MinTime != nil {
			w.addTime(row, *field.MinTime)
		} else {
			w.addString(row, "")
		}
		if field.MaxTime != nil {
			w.addTime(row, *field.MaxTime)
		} else {
			w.addString(row, "")
		}
		if field.BoolTrue != nil {
			w.addInt(row, *field.BoolTrue)
		} else {
			w.addString(row, "")
		}
		if field.BoolFalse != nil {
			w.addInt(row, *field.BoolFalse)
		} else {
			w.addString(row, "")
		}
	}
}

func (w *ReportExcelWriter) addString(row *xlsx.Row, value string) {
	cell := row.AddCell()
	cell.SetString(value)
}

func (w *ReportExcelWriter) addBool(row *xlsx.Row, value bool) {
	cell := row.AddCell()
	cell.SetBool(value)
}

func (w *ReportExcelWriter) addInt(row *xlsx.Row, value int) {
	cell := row.AddCell()
	cell.SetInt(value)
}

func (w *ReportExcelWriter) addFloat(row *xlsx.Row, value float64) {
	cell := row.AddCell()
	cell.SetFloat(value)
}

func (w *ReportExcelWriter) addTime(row *xlsx.Row, value time.Time) {
	cell := row.AddCell()
	cell.SetDateTime(value)
}
