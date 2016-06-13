package main

import (
	"fmt"
	"io"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
)

// Application object.
type Application struct {
	writer    *ReportWriter
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
	a.putReport(*report)
	return nil
}

// Run application core logic.
func (a *Application) cntblank(reader *Reader, dialect *FileDialect) (report *Report, err error) {
	logger := log.WithFields(a.logfields)
	report = new(Report)
	report.Path = reader.path
	report.Filename = filepath.Base(reader.path)
	report.Md5hex = reader.md5hex
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
		logger.Info("start parsing with ", len(report.Fields), " columns.")
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
		report.Records, len(report.Fields))
	return report, nil
}

func (a *Application) putReport(report Report) {
	err := a.writer.Write(report)
	if err != nil {
		log.Error(err)
	}
}

// newApplication creates `Application` object to set some options.
func newApplication(writer io.Writer, format string, dialect *FileDialect) (a *Application, err error) {
	a = new(Application)
	a.writer = NewReportWriter(writer, format, dialect)
	return a, nil
}
