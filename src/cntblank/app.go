package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	log "github.com/Sirupsen/logrus"
)

// Application object.
type Application struct {
	reports   []Report
	writer    *ReportWriter
	logfields log.Fields
}

// Run application main logic.
func (a *Application) Run(pathList []string, dialect *FileDialect) error {
	var targets []string
	if len(pathList) == 0 {
		targets = append(targets, "")
	} else {
		targets = expandDir(pathList)
	}
	a.reports = make([]Report, len(targets))
	for i, target := range targets {
		report, err := a.process(target, dialect)
		if err != nil {
			log.Errorf("[%d] error while processing %s: %v", i+1, target, err)
		}
		if report != nil {
			a.reports[i] = *report
		}
	}
	return a.putReport()
}

func (a *Application) process(target string, dialect *FileDialect) (report *Report, err error) {
	reader, err := OpenFile(target, dialect)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	return a.cntblank(reader, dialect.HasHeader)
}

// Run application core logic.
func (a *Application) cntblank(reader *Reader, hasHeader bool) (report *Report, err error) {
	logger := log.WithFields(a.logfields)
	report = new(Report)
	report.Path = reader.path
	if len(reader.path) > 0 {
		report.Filename = filepath.Base(reader.path)
	}
	report.MD5hex = reader.md5hex
	if hasHeader {
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

func (a *Application) putReport() error {
	log.Infof("write %d reports", len(a.reports))
	return a.writer.Write(a.reports)
}

// newApplication creates `Application` object to set some options.
func newApplication(writer io.Writer, format string, dialect *FileDialect) (a *Application, err error) {
	a = new(Application)
	a.writer = NewReportWriter(writer, format, dialect)
	return a, nil
}

// expandDir expands files in directory.
func expandDir(pathList []string) []string {
	// TODO: Optional recursive traversing.
	// TODO: Filter file name and/or extension to expand.
	var files []string
	for _, path := range pathList {
		if len(path) == 0 {
			continue
		}
		fileInfo, err := os.Stat(path)
		if err != nil {
			log.Infof("%v", err)
			continue
		}
		if fileInfo.IsDir() {
			log.Infof("expand directory: %s", path)
			filesInDir, err := ioutil.ReadDir(path)
			if err != nil {
				log.Fatal(err)
			}
			var children []string
			for _, file := range filesInDir {
				if strings.HasPrefix(file.Name(), ".") {
					// Skip hidden file.
					continue
				}
				if isTabularFormat(file.Name()) {
					// Only parse tablular format file.
					children = append(children, filepath.Join(path, file.Name()))
				} else if file.IsDir() {
					children = append(children, filepath.Join(path, file.Name()))
				}
			}
			c := expandDir(children)
			files = append(files, c...)
		} else {
			files = append(files, path)
		}
	}
	return files
}

func isTabularFormat(fname string) bool {
	// TODO: Require more flexible filter implementation.
	formats := []string{
		".csv",
		".tsv",
		".txt",
		".xlsx",
	}
	ext := strings.ToLower(path.Ext(fname))
	for _, fmt := range formats {
		if ext == fmt {
			return true
		}
	}
	return false
}
