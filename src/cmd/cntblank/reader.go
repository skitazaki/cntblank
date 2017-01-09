package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/tealeg/xlsx"

	"csvhelper"
)

// Reader is a generic file reader.
type Reader struct {
	path      string
	line      int
	columns   map[int]int
	err       int
	fp        *os.File
	csvReader *csv.Reader
	slices    [][]string
	logger    *log.Entry
}

// NewReader returns a new Reader that reads from r using dialect.
func NewReader(r io.Reader, dialect *csvhelper.FileDialect) (reader *Reader, err error) {
	reader = &Reader{
		columns: make(map[int]int),
		logger:  log.WithFields(nil),
	}
	reader.csvReader = csvhelper.NewCsvReader(r, dialect)
	return
}

// OpenFile returns a new Reader that reads from path using dialect.
func OpenFile(path string, dialect *csvhelper.FileDialect) (reader *Reader, err error) {
	if path == "" {
		return NewReader(os.Stdin, dialect)
	}
	extension := filepath.Ext(path)
	if extension == ".xlsx" {
		slices, err := xlsx.FileToSlice(path)
		if err != nil {
			return nil, err
		}
		var feeds [][]string
		if dialect.SheetNumber == 0 {
			feeds = slices[0]
		} else {
			if dialect.SheetNumber > len(slices) {
				return nil, fmt.Errorf("%s has only %d sheets, given sheet number is %d",
					path, len(slices), dialect.SheetNumber)
			}
			feeds = slices[dialect.SheetNumber-1]
		}
		reader = &Reader{
			columns: make(map[int]int),
			slices:  feeds,
		}
	} else {
		fp, err := os.Open(path)
		// TODO: Check `fp` is file or directory.
		// http://www.reddit.com/r/golang/comments/2fjwyk/isdir_in_go/
		if err != nil {
			return nil, err
		}
		reader, err = NewReader(fp, dialect)
		reader.fp = fp
	}
	reader.path = path
	reader.logger = log.WithFields(log.Fields{"path": path})
	return
}

func (r *Reader) Read() (record []string, err error) {
	// TODO: do not depend on `csv.Reader` interface.
	// This is a tentative implementation.
	if r.csvReader != nil {
		record, err = r.csvReader.Read()
		if err == io.EOF {
			// Report the summary.
			r.logger.Infof("finish parsing %d lines with %d errors", r.line, r.err)
			for col, count := range r.columns {
				r.logger.Infof("  column size %d has %d lines", col, count)
			}
			return nil, err
		} else if err != nil {
			r.logger.Error(err, ", #line", r.line)
			r.err++
			if r.err > 100 {
				r.logger.Error("too many error lines")
				return nil, fmt.Errorf("too many error lines")
			}
			return nil, err
		}
	} else if r.slices != nil {
		if len(r.slices) <= r.line {
			r.logger.Infof("finish parsing %d lines with %d errors", r.line, r.err)
			return nil, io.EOF
		}
		record = r.slices[r.line]
	}
	r.line++
	length := len(record)
	_, ok := r.columns[length]
	if ok {
		r.columns[length]++
	} else {
		r.columns[length] = 1
	}
	// Show simple progress report.
	if r.line%1000000 == 0 {
		r.logger.Infof("==> Processed %d lines <==", r.line)
	}
	return record, nil
}

// Close closes a internal file pointer.
func (r *Reader) Close() {
	if r.fp != nil {
		r.fp.Close()
	}
}
