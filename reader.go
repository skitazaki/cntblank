package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"

	log "github.com/Sirupsen/logrus"
	"github.com/tealeg/xlsx"
)

// FileDialect is a configuration for "csv.Reader".
type FileDialect struct {
	Encoding         string // file encoding (utf8 or sjis only)
	Comma            rune   // field delimiter (set to ',' by NewReader)
	Comment          rune   // comment character for start of line
	FieldsPerRecord  int    // number of expected fields per record
	LazyQuotes       bool   // allow lazy quotes
	TrimLeadingSpace bool   // trim leading space
	HasHeader        bool   // CSV file has header line
	HasMetadata      bool   // meta data before header line
}

// Reader is a generic file reader.
type Reader struct {
	Path      string
	line      int
	columns   map[int]int
	err       int
	fp        *os.File
	csvReader *csv.Reader
	slices    [][][]string
	logger    *log.Entry
}

// NewReader returns a new Reader that reads from r using dialect.
func NewReader(r io.Reader, dialect *FileDialect) (reader *Reader, err error) {
	reader = &Reader{
		columns: make(map[int]int),
		logger:  log.WithFields(nil),
	}
	err = reader.setupCsvReader(bufio.NewReader(r), dialect)
	return
}

// OpenFile returns a new Reader that reads from path using dialect.
func OpenFile(path string, dialect *FileDialect) (reader *Reader, err error) {
	if path == "" {
		return NewReader(os.Stdin, dialect)
	}
	extension := filepath.Ext(path)
	if extension == ".xlsx" {
		slices, err := xlsx.FileToSlice(path)
		if err != nil {
			return nil, err
		}
		reader = &Reader{
			columns: make(map[int]int),
			logger:  log.WithFields(nil),
			slices:  slices,
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
	reader.Path = path
	reader.logger = log.WithFields(log.Fields{"path": path})
	return
}

func (r *Reader) setupCsvReader(reader io.Reader, dialect *FileDialect) error {
	if dialect.Encoding == "sjis" {
		r.logger.Info("use ShiftJIS decoder for input.")
		decoder := japanese.ShiftJIS.NewDecoder()
		r.csvReader = csv.NewReader(transform.NewReader(reader, decoder))
	} else {
		r.csvReader = csv.NewReader(reader)
	}
	r.csvReader.Comma = dialect.Comma
	r.csvReader.Comment = dialect.Comment
	r.csvReader.FieldsPerRecord = dialect.FieldsPerRecord
	return nil
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
		// TODO: Set sheet number or name at runtime.
		if len(r.slices[0]) <= r.line {
			r.logger.Infof("finish parsing %d lines with %d errors", r.line, r.err)
			return nil, io.EOF
		}
		record = r.slices[0][r.line]
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
