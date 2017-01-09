package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	valid "github.com/asaskevich/govalidator"
)

// Report presents tabular contents description.
type Report struct {
	Path      string         `json:"path,omitempty"`
	Filename  string         `json:"filename,omitempty"`
	MD5hex    string         `json:"md5,omitempty"`
	HasHeader bool           `json:"header"`
	Records   int            `json:"records"`
	Fields    []*ReportField `json:"fields"`
}

// ReportField represents output field.
type ReportField struct {
	Name      string     `json:"name"`
	Blank     int        `json:"blank"`
	MinLength int        `json:"minLength"`
	MaxLength int        `json:"maxLength"`
	Minimum   *float64   `json:"minimum,omitempty"`
	Maximum   *float64   `json:"maximum,omitempty"`
	MinTime   *time.Time `json:"minTime,omitempty"`
	MaxTime   *time.Time `json:"maxTime,omitempty"`
	BoolTrue  *int       `json:"boolTrue,omitempty"`
	BoolFalse *int       `json:"boolFalse,omitempty"`
	TypeInt   int        `json:"typeInt,omitempty"`
	TypeFloat int        `json:"typeFloat,omitempty"`
	TypeBool  int        `json:"typeBool,omitempty"`
	TypeTime  int        `json:"typeTime,omitempty"`
	fullWidth int
}

func (r ReportField) header() []string {
	return []string{
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
}

func (r *ReportField) format(total int) []string {
	s := make([]string, 14)
	s[1] = r.Name
	s[2] = fmt.Sprint(r.Blank)
	ratio := float64(r.Blank) / float64(total)
	s[3] = fmt.Sprintf("%.4f", ratio)
	if r.MinLength > 0 {
		s[4] = fmt.Sprint(r.MinLength)
	} else {
		s[4] = ""
	}
	if r.MaxLength > 0 {
		s[5] = fmt.Sprint(r.MaxLength)
	} else {
		s[5] = ""
	}
	if r.TypeInt > 0 {
		s[6] = fmt.Sprint(r.TypeInt)
	} else {
		s[6] = ""
	}
	if r.TypeFloat > 0 {
		s[7] = fmt.Sprint(r.TypeFloat)
	} else {
		s[7] = ""
	}
	if r.TypeBool > 0 {
		s[8] = fmt.Sprint(r.TypeBool)
	} else {
		s[8] = ""
	}
	if r.TypeTime > 0 {
		s[9] = fmt.Sprint(r.TypeTime)
	} else {
		s[9] = ""
	}
	// Min/Max comparison.
	if r.TypeTime > r.TypeFloat {
		s[10] = r.MinTime.Format("2006-01-02 15:04:05")
		s[11] = r.MaxTime.Format("2006-01-02 15:04:05")
	} else if r.TypeFloat > 0 {
		s[10] = fmt.Sprintf("%.4f", *r.Minimum)
		s[11] = fmt.Sprintf("%.4f", *r.Maximum)
	} else if r.TypeInt > 0 {
		s[10] = fmt.Sprint(*r.Minimum)
		s[11] = fmt.Sprint(*r.Maximum)
	} else {
		s[10] = ""
		s[11] = ""
	}
	if r.TypeBool > 0 {
		s[12] = fmt.Sprint(*r.BoolTrue)
		s[13] = fmt.Sprint(*r.BoolFalse)
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
		name := strings.TrimSpace(record[i])
		if name != "" {
			f.Name = strings.Replace(name, "\n", "", -1)
		} else {
			f.Name = fmt.Sprintf("Column%03d", i+1)
		}
		r.Fields = append(r.Fields, f)
	}
	r.HasHeader = true
	return nil
}

func (r *Report) parseRecord(record []string) (nullCount int) {
	r.Records++
	size := len(record)
	if size > len(r.Fields) {
		for i := len(r.Fields); i < size; i++ {
			f := new(ReportField)
			f.Name = fmt.Sprintf("Column%03d", i+1)
			f.Blank = r.Records - 1 // suppose all cells are blank until up to here.
			r.Fields = append(r.Fields, f)
		}
	}
	for i := 0; i < size; i++ {
		f := r.Fields[i]
		val := strings.TrimSpace(record[i])
		if len(val) == 0 {
			nullCount++
			f.Blank++
			continue
		}
		stringLength := utf8.RuneCountInString(val)
		if f.MinLength == 0 || f.MinLength > stringLength {
			f.MinLength = stringLength
		}
		if f.MaxLength < stringLength {
			f.MaxLength = stringLength
		}
		if valid.IsFullWidth(val) {
			f.fullWidth++
		}
		if valInt, err := strconv.Atoi(val); err == nil {
			v := float64(valInt)
			if f.Minimum == nil {
				f.Minimum = new(float64)
				*f.Minimum = v
			}
			if f.Maximum == nil {
				f.Maximum = new(float64)
				*f.Maximum = v
			}
			if v < *f.Minimum {
				*f.Minimum = v
			}
			if v > *f.Maximum {
				*f.Maximum = v
			}
			f.TypeInt++
		}
		if valFloat, err := strconv.ParseFloat(val, 64); err == nil {
			if f.Minimum == nil {
				f.Minimum = new(float64)
				*f.Minimum = valFloat
			}
			if f.Maximum == nil {
				f.Maximum = new(float64)
				*f.Maximum = valFloat
			}
			if valFloat < *f.Minimum {
				*f.Minimum = valFloat
			}
			if valFloat > *f.Maximum {
				*f.Maximum = valFloat
			}
			f.TypeFloat++
		}
		if valBool, err := strconv.ParseBool(val); err == nil {
			if f.TypeBool == 0 {
				f.BoolTrue = new(int)
				f.BoolFalse = new(int)
			}
			if valBool {
				*f.BoolTrue++
			} else {
				*f.BoolFalse++
			}
			f.TypeBool++
		}
		if valTime, err := parseDateTime(val); err == nil {
			if f.TypeTime == 0 {
				f.MinTime = new(time.Time)
				f.MaxTime = new(time.Time)
				*f.MinTime = valTime
				*f.MaxTime = valTime
			}
			if valTime.Before(*f.MinTime) {
				*f.MinTime = valTime
			}
			if valTime.After(*f.MaxTime) {
				*f.MaxTime = valTime
			}
			f.TypeTime++
		}
	}
	return nullCount
}

func newReport(f File) *Report {
	r := new(Report)
	if f.path != "" {
		r.Path = f.path
		r.Filename = f.Name()
		if md5hex, err := f.Checksum(); err == nil {
			r.MD5hex = md5hex
		}
	}
	return r
}
