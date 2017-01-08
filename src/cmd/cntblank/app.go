package main

import (
	"fmt"
	"io"

	log "github.com/Sirupsen/logrus"

	"csvhelper"
)

// Application object.
type Application struct {
	collector *FileCollector
	reports   []Report
	writer    ReportWriter
	logfields log.Fields
}

// Run application main logic.
func (a *Application) Run(pathList []string, dialect *csvhelper.FileDialect) error {
	var files []File
	if len(pathList) == 0 {
		files = append(files, File{})
	} else {
		err := a.collector.CollectAll(pathList)
		if err != nil {
			return err
		}
		files = a.collector.files
	}
	a.reports = make([]Report, len(files))
	for i, file := range files {
		report := newReport(file)
		err := a.process(report, dialect)
		if err != nil {
			log.Errorf("[%d] error while processing %s: %v", i+1, file.path, err)
			continue
		}
		a.reports[i] = *report
	}
	return a.putReport()
}

func (a *Application) process(report *Report, dialect *csvhelper.FileDialect) error {
	reader, err := OpenFile(report.Path, dialect)
	if err != nil {
		return err
	}
	defer reader.Close()

	return a.cntblank(report, reader, dialect.HasHeader)
}

// Run application core logic.
func (a *Application) cntblank(report *Report, reader *Reader, hasHeader bool) error {
	logger := log.WithFields(a.logfields)
	if hasHeader {
		// Use first line as header name if flag is not specified.
		record, err := reader.Read()
		if err == io.EOF {
			return fmt.Errorf("reader is empty")
		} else if err != nil {
			return err
		}
		err = report.header(record)
		if err != nil {
			logger.Error(err)
			return err
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
	return nil
}

func (a *Application) putReport() error {
	log.Infof("write %d reports", len(a.reports))
	return a.writer.Write(a.reports)
}

// newApplication creates `Application` object to set some options.
func newApplication(recursive bool, writer io.Writer, format string, dialect *csvhelper.FileDialect) (a *Application, err error) {
	f := CSV // default output format is CSV
	if format != "" {
		f = formatFrom(format)
		if f == Unknown {
			return nil, fmt.Errorf("unknown format %q", format)
		}
	}
	a = new(Application)
	a.collector = newFileCollector(recursive, []string{
		".csv",
		".tsv",
		".txt",
		".xlsx",
	})
	a.writer = NewReportWriter(writer, f, dialect)
	return a, nil
}
