package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode/utf8"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"

	log "github.com/Sirupsen/logrus"
)

// Application object.
type Application struct {
	writer    *csv.Writer
	putMeta   bool
	logfields log.Fields
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
	trueCount  int
	falseCount int
	intType    int
	floatType  int
	boolType   int
}

func (r *ReportField) header() []string {
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
		"Minimum",
		"Maximum",
		"MinimumF",
		"MaximumF",
		"#True",
		"#False",
	}
}
func (r *ReportField) format(total int) []string {
	s := make([]string, 15)
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
	if r.intType > 0 {
		s[9] = fmt.Sprint(r.minimum)
		s[10] = fmt.Sprint(r.maximum)
	} else {
		s[9] = ""
		s[10] = ""
	}
	if r.floatType > 0 {
		s[11] = fmt.Sprintf("%.4f", r.minimumF)
		s[12] = fmt.Sprintf("%.4f", r.maximumF)
	} else {
		s[11] = ""
		s[12] = ""
	}
	if r.boolType > 0 {
		s[13] = fmt.Sprint(r.trueCount)
		s[14] = fmt.Sprint(r.falseCount)
	} else {
		s[13] = ""
		s[14] = ""
	}
	return s
}

func (r *Report) header(record []string) error {
	if len(record) == 0 {
		return fmt.Errorf("header record has no elements.")
	}
	for i := 0; i < len(record); i++ {
		f := new(ReportField)
		f.seq = i + 1
		f.name = strings.TrimSpace(record[i])
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
		} else if valFloat, err := strconv.ParseFloat(val, 64); err == nil {
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
		} else if valBool, err := strconv.ParseBool(val); err == nil {
			if valBool {
				f.trueCount++
			} else {
				f.falseCount++
			}
			f.boolType++
		}
	}
	return nullCount
}

// Run application main logic.
func (a *Application) run(path string, encoding string, delimiter string,
	noHeader bool, strict bool) error {
	var buffer *bufio.Reader
	if len(path) > 0 {
		fp, err := os.Open(path)
		// TODO: Check `fp` is file or directory.
		// http://www.reddit.com/r/golang/comments/2fjwyk/isdir_in_go/
		if err != nil {
			log.Error(err)
			return err
		}
		defer fp.Close()
		buffer = bufio.NewReader(fp)
		a.logfields = log.Fields{"path": path}
	} else {
		buffer = bufio.NewReader(os.Stdin)
		a.logfields = log.Fields{}
	}

	logger := log.WithFields(a.logfields)

	var reader *csv.Reader
	if len(encoding) > 0 {
		if encoding == "sjis" {
			logger.Info("use ShiftJIS decoder for input.")
			decoder := japanese.ShiftJIS.NewDecoder()
			r := transform.NewReader(buffer, decoder)
			reader = csv.NewReader(r)
		} else {
			logger.Warn("unknown encoding: ", encoding)
			reader = csv.NewReader(buffer)
		}
	} else {
		reader = csv.NewReader(buffer)
	}

	report, err := a.cntblank(reader, delimiter, noHeader, strict)
	if err != nil {
		logger.Error(err)
		return err
	}
	report.path = path
	a.putReport(*report)
	return nil
}

// Run application core logic.
func (a *Application) cntblank(reader *csv.Reader, delimiter string, noHeader bool, strict bool) (report *Report, err error) {
	logger := log.WithFields(a.logfields)
	if len(delimiter) > 0 {
		comma, err := utf8.DecodeRuneInString(delimiter)
		if err == utf8.RuneError {
			logger.Warn(err)
			logger.Info("input delimiter option is invalid, but continue running.")
			reader.Comma = '\t'
		} else {
			reader.Comma = comma
		}
	} else {
		reader.Comma = '\t'
	}
	reader.Comment = '#'
	if strict {
		reader.FieldsPerRecord = 0
	} else {
		reader.FieldsPerRecord = -1
	}
	lines := 0
	report = new(Report)
	if noHeader {
		logger.Info("start parsing without header row")
	} else {
		// Use first line as header name if flag is not specified.
		record, err := reader.Read()
		lines++
		if err == io.EOF {
			return nil, fmt.Errorf("reader is empty")
		} else if err != nil {
			logger.Error(err)
			return nil, err
		}
		err = report.header(record)
		if err != nil {
			logger.Error(err)
			return nil, err
		}
		logger.Info("start parsing with ", len(report.fields), " columns.")
	}
	errCount := 0
	for {
		record, err := reader.Read()
		lines++
		if err == io.EOF {
			lines--
			break
		} else if err != nil {
			logger.Error(err, ", #line", lines)
			errCount++
			if errCount > 100 {
				return nil, fmt.Errorf("too many error lines")
			}
			continue
		}
		nullCount := report.parseRecord(record)
		if nullCount > 0 {
			logger.Debugf("line #%d has %d fields with %d NULL(s).",
				lines, len(record), nullCount)
		}
		if lines%1000000 == 0 {
			logger.Info("==> Processed ", lines, " lines <==")
		}
	}
	logger.Infof("finish parsing %d lines to get %d records with %d columns. %d errors detected.",
		lines, report.records, len(report.fields), errCount)
	return report, nil
}

func (a *Application) putReport(report Report) {
	if a.putMeta {
		preamble := make([]string, 3)
		if len(report.path) > 0 {
			preamble[0] = "# File"
			preamble[1] = report.path
			preamble[2] = filepath.Base(report.path)
			a.writer.Write(preamble)
		}
		preamble[0] = "# Field"
		preamble[1] = fmt.Sprint(len(report.fields))
		if report.hasHeader {
			preamble[2] = "(has header)"
		} else {
			preamble[2] = ""
		}
		a.writer.Write(preamble)
		preamble[0] = "# Record"
		preamble[1] = fmt.Sprint(report.records)
		preamble[2] = ""
		a.writer.Write(preamble)
	}
	// Put header line.
	a.writer.Write(new(ReportField).header())
	// Put each field report.
	for i := 0; i < len(report.fields); i++ {
		r := report.fields[i]
		a.writer.Write(r.format(report.records))
	}
	a.writer.Flush()
}

// Create `Application` object to set some options.
func newApplication(writer io.Writer, encoding string, delimiter string, meta bool) (a *Application, err error) {
	a = new(Application)
	if len(encoding) > 0 {
		if encoding == "sjis" {
			log.Info("use ShiftJIS encoder for output.")
			encoder := japanese.ShiftJIS.NewEncoder()
			writer = transform.NewWriter(writer, encoder)
		} else {
			log.Warn("unknown encoding: ", encoding)
		}
	}
	a.writer = csv.NewWriter(writer)
	if len(delimiter) > 0 {
		comma, err := utf8.DecodeRuneInString(delimiter)
		if err == utf8.RuneError {
			log.Warn(err)
			log.Info("output delimiter option is invalid, but continue running.")
			a.writer.Comma = '\t'
		} else {
			a.writer.Comma = comma
		}
	} else {
		a.writer.Comma = '\t'
	}
	a.putMeta = meta
	return a, nil
}
