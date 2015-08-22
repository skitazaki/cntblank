package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
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
	path    string
	records int
	fields  []ReportField
}

// Output field.
type ReportField struct {
	seq       int
	name      string
	blank     int
	minLength int
	maxLength int
}

func (r *ReportField) header() []string {
	s := make([]string, 6)
	s[0] = "seq"
	s[1] = "Name"
	s[2] = "#Blank"
	s[3] = "%Blank"
	s[4] = "MinLength"
	s[5] = "MaxLength"
	return s
}
func (r *ReportField) format(total int) []string {
	s := make([]string, 6)
	s[0] = fmt.Sprint(r.seq)
	s[1] = r.name
	s[2] = fmt.Sprint(r.blank)
	ratio := float64(r.blank) / float64(total) * 100
	s[3] = fmt.Sprintf("%02.4f", ratio)
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
	return s
}

// Run application main logic.
func (a *Application) run(path *string, encoding *string, delimiter *string,
	noHeader bool, strict bool) error {
	var buffer *bufio.Reader
	if path != nil && len(*path) > 0 {
		fp, err := os.Open(*path)
		// TODO: Check `fp` is file or directory.
		// http://www.reddit.com/r/golang/comments/2fjwyk/isdir_in_go/
		if err != nil {
			log.Error(err)
			return err
		}
		defer fp.Close()
		buffer = bufio.NewReader(fp)
		a.logfields = log.Fields{"path": *path}
		if a.putMeta {
			preamble := make([]string, 2)
			preamble[0] = "# File"
			preamble[1] = *path
			a.writer.Write(preamble)
		}
	} else {
		buffer = bufio.NewReader(os.Stdin)
		a.logfields = log.Fields{}
	}

	logger := log.WithFields(a.logfields)

	var reader *csv.Reader
	if encoding != nil && len(*encoding) > 0 {
		if *encoding == "sjis" {
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
	a.putReport(*report)
	return nil
}

// Run application core logic.
func (a *Application) cntblank(reader *csv.Reader, delimiter *string, noHeader bool, strict bool) (report *Report, err error) {
	logger := log.WithFields(a.logfields)
	if delimiter != nil && len(*delimiter) > 0 {
		comma, err := utf8.DecodeRuneInString(*delimiter)
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
	fields := make(map[int]*ReportField)
	if noHeader {
		logger.Info("start parsing without header row")
	} else {
		// Use first line as header name if flag is not specified.
		record, err := reader.Read()
		if err == io.EOF {
			return nil, fmt.Errorf("reader is empty")
		} else if err != nil {
			logger.Error(err)
			return nil, err
		}
		for i := 0; i < len(record); i++ {
			f := new(ReportField)
			f.seq = i + 1
			f.name = record[i]
			fields[i] = f
		}
		logger.Info("start parsing with ", len(fields), " columns.")
	}
	recordCount := 0
	errCount := 0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			logger.Error(err, ", #record", recordCount-1)
			recordCount++
			errCount++
			if errCount > 100 {
				return nil, fmt.Errorf("too many error lines")
			}
			continue
		}
		nullCount := 0
		for i := 0; i < len(record); i++ {
			f, ok := fields[i]
			if !ok {
				f = new(ReportField)
				f.seq = i + 1
				f.name = fmt.Sprintf("Column%03d", i+1)
				fields[i] = f
			}
			if len(record[i]) == 0 {
				nullCount++
				f.blank++
			} else {
				stringLength := utf8.RuneCountInString(record[i])
				if f.minLength == 0 || f.minLength > stringLength {
					f.minLength = stringLength
				}
				if f.maxLength < stringLength {
					f.maxLength = stringLength
				}
			}
		}
		if nullCount > 0 {
			logger.Debugf("record #%d has %d fields with %d NULL(s).",
				recordCount, len(record), nullCount)
		}
		recordCount++
		if recordCount%1000000 == 0 {
			logger.Info("==> Processed ", recordCount, " lines <==")
		}
	}
	columnSize := len(fields)
	logger.Infof("finish parsing %d records with %d columns. %d errors detected.",
		recordCount, columnSize, errCount)
	report = new(Report)
	report.records = recordCount
	report.fields = make([]ReportField, columnSize)
	for i := 0; i < columnSize; i++ {
		report.fields[i] = *fields[i]
	}
	return report, nil
}

func (a *Application) putReport(report Report) {
	if a.putMeta {
		preamble := make([]string, 2)
		preamble[0] = "# Field"
		preamble[1] = fmt.Sprint(len(report.fields))
		a.writer.Write(preamble)
		preamble[0] = "# Record"
		preamble[1] = fmt.Sprint(report.records)
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
func newApplication(writer io.Writer, encoding *string, delimiter *string, meta bool) (a *Application, err error) {
	a = new(Application)
	if encoding != nil && len(*encoding) > 0 {
		if *encoding == "sjis" {
			log.Info("use ShiftJIS encoder for output.")
			encoder := japanese.ShiftJIS.NewEncoder()
			writer = transform.NewWriter(writer, encoder)
		} else {
			log.Warn("unknown encoding: ", *encoding)
		}
	}
	a.writer = csv.NewWriter(writer)
	if delimiter != nil && len(*delimiter) > 0 {
		comma, err := utf8.DecodeRuneInString(*delimiter)
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
