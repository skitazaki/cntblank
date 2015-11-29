package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"

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

// Run application main logic.
func (a *Application) run(path string, encoding string, delimiter rune,
	noHeader bool, strict bool) error {
	dialect := &FileDialect{
		Encoding:        encoding,
		Comma:           delimiter,
		Comment:         '#',
		FieldsPerRecord: -1,
		HasHeader:       !noHeader,
	}
	if strict {
		dialect.FieldsPerRecord = 0
	}
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
	if dialect.Encoding == "sjis" {
		logger.Info("use ShiftJIS decoder for input.")
		decoder := japanese.ShiftJIS.NewDecoder()
		r := transform.NewReader(buffer, decoder)
		reader = csv.NewReader(r)
	} else {
		reader = csv.NewReader(buffer)
	}

	report, err := a.cntblank(reader, dialect)
	if err != nil {
		logger.Error(err)
		return err
	}
	report.path = path
	a.putReport(*report)
	return nil
}

// Run application core logic.
func (a *Application) cntblank(reader *csv.Reader, dialect *FileDialect) (report *Report, err error) {
	logger := log.WithFields(a.logfields)
	reader.Comma = dialect.Comma
	reader.Comment = dialect.Comment
	reader.FieldsPerRecord = dialect.FieldsPerRecord
	lines := 0
	report = new(Report)
	if dialect.HasHeader {
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
	} else {
		logger.Info("start parsing without header row")
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
	writer := NewReportWriter(a.writer, a.putMeta)
	err := writer.Write(report)
	if err != nil {
		log.Error(err)
	}
}

// Create `Application` object to set some options.
func newApplication(writer io.Writer, encoding string, delimiter rune, meta bool) (a *Application, err error) {
	a = new(Application)
	if encoding == "sjis" {
		log.Info("use ShiftJIS encoder for output.")
		encoder := japanese.ShiftJIS.NewEncoder()
		writer = transform.NewWriter(writer, encoder)
	}
	a.writer = csv.NewWriter(writer)
	a.writer.Comma = delimiter
	a.putMeta = meta
	return a, nil
}
