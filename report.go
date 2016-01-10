package main

import (
	"encoding/csv"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	valid "github.com/asaskevich/govalidator"
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

// Calculate report.
type Report struct {
	path      string
	hasHeader bool
	records   int
	fields    []*ReportField
}

// Output field.
type ReportField struct {
	seq        int
	name       string
	blank      int
	minLength  int
	maxLength  int
	minimum    int
	maximum    int
	minimumF   float64
	maximumF   float64
	minimumT   time.Time
	maximumT   time.Time
	trueCount  int
	falseCount int
	intType    int
	floatType  int
	boolType   int
	timeType   int
	fullWidth  int
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

func (r *ReportField) format(total int) []string {
	s := make([]string, 14)
	s[0] = fmt.Sprint(r.seq)
	s[1] = r.name
	s[2] = fmt.Sprint(r.blank)
	ratio := float64(r.blank) / float64(total)
	s[3] = fmt.Sprintf("%.4f", ratio)
	if r.minLength > 0 {
		s[4] = fmt.Sprint(r.minLength)
	} else {
		s[4] = ""
	}
	if r.maxLength > 0 {
		s[5] = fmt.Sprint(r.maxLength)
	} else {
		s[5] = ""
	}
	if r.intType > 0 {
		s[6] = fmt.Sprint(r.intType)
	} else {
		s[6] = ""
	}
	if r.floatType > 0 {
		s[7] = fmt.Sprint(r.floatType)
	} else {
		s[7] = ""
	}
	if r.boolType > 0 {
		s[8] = fmt.Sprint(r.boolType)
	} else {
		s[8] = ""
	}
	if r.timeType > 0 {
		s[9] = fmt.Sprint(r.timeType)
	} else {
		s[9] = ""
	}
	// Min/Max comparison.
	if r.timeType > r.floatType {
		s[10] = r.minimumT.Format("2006-01-02 15:04:05")
		s[11] = r.maximumT.Format("2006-01-02 15:04:05")
	} else if r.floatType > 0 && r.intType > 0 {
		if float64(r.minimum) <= r.minimumF {
			s[10] = fmt.Sprint(r.minimum)
		} else {
			s[10] = fmt.Sprintf("%.4f", r.minimumF)
		}
		if float64(r.maximum) <= r.maximumF {
			s[11] = fmt.Sprint(r.maximum)
		} else {
			s[11] = fmt.Sprintf("%.4f", r.maximumF)
		}
	} else if r.floatType > 0 {
		s[10] = fmt.Sprintf("%.4f", r.minimumF)
		s[11] = fmt.Sprintf("%.4f", r.maximumF)
	} else if r.intType > 0 {
		s[10] = fmt.Sprint(r.minimum)
		s[11] = fmt.Sprint(r.maximum)
	} else {
		s[10] = ""
		s[11] = ""
	}
	if r.boolType > 0 {
		s[12] = fmt.Sprint(r.trueCount)
		s[13] = fmt.Sprint(r.falseCount)
	} else {
		s[12] = ""
		s[13] = ""
	}
	return s
}

func (r *Report) header(record []string) error {
	if len(record) == 0 {
		return fmt.Errorf("header record has no elements")
	}
	for i := 0; i < len(record); i++ {
		f := new(ReportField)
		f.seq = i + 1
		name := strings.TrimSpace(record[i])
		if name != "" {
			f.name = strings.Replace(name, "\n", "", -1)
		} else {
			f.name = fmt.Sprintf("Column%03d", i+1)
		}
		r.fields = append(r.fields, f)
	}
	r.hasHeader = true
	return nil
}

func (r *Report) parseRecord(record []string) (nullCount int) {
	r.records++
	size := len(record)
	if size > len(r.fields) {
		for i := len(r.fields); i < size; i++ {
			f := new(ReportField)
			f.seq = i + 1
			f.name = fmt.Sprintf("Column%03d", i+1)
			r.fields = append(r.fields, f)
		}
	}
	for i := 0; i < size; i++ {
		f := r.fields[i]
		val := strings.TrimSpace(record[i])
		if len(val) == 0 {
			nullCount++
			f.blank++
			continue
		}
		stringLength := utf8.RuneCountInString(val)
		if f.minLength == 0 || f.minLength > stringLength {
			f.minLength = stringLength
		}
		if f.maxLength < stringLength {
			f.maxLength = stringLength
		}
		if valid.IsFullWidth(val) {
			f.fullWidth++
		}
		if valInt, err := strconv.Atoi(val); err == nil {
			if f.intType == 0 {
				f.minimum = valInt
				f.maximum = valInt
			} else {
				if valInt < f.minimum {
					f.minimum = valInt
				}
				if valInt > f.maximum {
					f.maximum = valInt
				}
			}
			f.intType++
		}
		if valFloat, err := strconv.ParseFloat(val, 64); err == nil {
			if f.floatType == 0 {
				f.minimumF = valFloat
				f.maximumF = valFloat
			} else {
				if valFloat < f.minimumF {
					f.minimumF = valFloat
				}
				if valFloat > f.maximumF {
					f.maximumF = valFloat
				}
			}
			f.floatType++
		}
		if valBool, err := strconv.ParseBool(val); err == nil {
			if valBool {
				f.trueCount++
			} else {
				f.falseCount++
			}
			f.boolType++
		}
		if valTime, err := parseDateTime(val); err == nil {
			if f.timeType == 0 {
				f.minimumT = valTime
				f.maximumT = valTime
			} else {
				if valTime.Before(f.minimumT) {
					f.minimumT = valTime
				}
				if valTime.After(f.maximumT) {
					f.maximumT = valTime
				}
			}
			f.timeType++
		}
	}
	return nullCount
}

func (w *ReportWriter) Write(report Report) error {
	if w.showMetadata {
		preamble := make([]string, 3)
		if len(report.path) > 0 {
			preamble[0] = "# File"
			preamble[1] = report.path
			preamble[2] = filepath.Base(report.path)
			w.w.Write(preamble)
		}
		preamble[0] = "# Field"
		preamble[1] = fmt.Sprint(len(report.fields))
		if report.hasHeader {
			preamble[2] = "(has header)"
		} else {
			preamble[2] = ""
		}
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
