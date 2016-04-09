package main

import (
	"encoding/csv"
	"fmt"
	"io"

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
func (a *Application) run(path string, dialect *FileDialect) error {
	reader, err := OpenFile(path, dialect)
	if err != nil {
		return err
	}
	defer reader.Close()

	report, err := a.cntblank(reader, dialect)
	if err != nil {
		return err
	}
	report.path = path
	report.md5hex = reader.md5hex
	a.putReport(*report)
	return nil
}

// Run application core logic.
func (a *Application) cntblank(reader *Reader, dialect *FileDialect) (report *Report, err error) {
	logger := log.WithFields(a.logfields)
	report = new(Report)
	if dialect.HasHeader {
		// Use first line as header name if flag is not specified.
		record, err := reader.Read()
		if err == io.EOF {
			return nil, fmt.Errorf("reader is empty")
		} else if err != nil {
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
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			if reader.err > 100 {
				break
			}
			continue
		}
		nullCount := report.parseRecord(record)
		if nullCount > 0 {
			logger.Debugf("line #%d has %d fields with %d NULL(s).",
				reader.line, len(record), nullCount)
		}
	}
	logger.Infof("get %d records with %d columns",
		report.records, len(report.fields))
	return report, nil
}

func (a *Application) putReport(report Report) {
	writer := NewReportWriter(a.writer, a.putMeta)
	err := writer.Write(report)
	if err != nil {
		log.Error(err)
	}
}

// newApplication creates `Application` object to set some options.
func newApplication(writer io.Writer, dialect *FileDialect) (a *Application, err error) {
	a = new(Application)
	if dialect.Encoding == "sjis" {
		log.Info("use ShiftJIS encoder for output.")
		encoder := japanese.ShiftJIS.NewEncoder()
		writer = transform.NewWriter(writer, encoder)
	}
	a.writer = csv.NewWriter(writer)
	a.writer.Comma = dialect.Comma
	a.putMeta = dialect.HasMetadata
	return a, nil
}
