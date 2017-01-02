package csvhelper

import (
	"encoding/csv"
	"fmt"
	"io"
	"unicode/utf8"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

// FileDialect is a configuration for reader and writer.
type FileDialect struct {
	Comma            rune   // field delimiter (set to ',' by NewReader)
	Comment          rune   // comment character for start of line
	Encoding         string // file encoding (utf8 or sjis only)
	FieldsPerRecord  int    // number of expected fields per record
	HasHeader        bool   // CSV file has header line
	HasMetadata      bool   // meta data before header line
	LazyQuotes       bool   // allow lazy quotes
	SheetNumber      int    // sheet number in Excel file which starts with 1
	TrimLeadingSpace bool   // trim leading space
}

var defaults = FileDialect{
	Comma:            ',',
	Comment:          '#',
	Encoding:         "utf8",
	FieldsPerRecord:  -1,
	HasHeader:        true,
	LazyQuotes:       true,
	SheetNumber:      1,
	TrimLeadingSpace: true,
}

// NewFileDialect creates new FileDialect instance.
func NewFileDialect(delimiter, encoding string, hasHeader bool) (*FileDialect, error) {
	var comma rune
	if len(delimiter) > 0 {
		c, size := utf8.DecodeRuneInString(delimiter)
		if size == utf8.RuneError {
			err := fmt.Errorf("delimiter is invalid rune %q", delimiter)
			return nil, err
		}
		comma = c
	}
	if len(encoding) > 0 {
		// TODO: define the list of encodings and check the given is included the list.
	}
	return &FileDialect{
		Comma:            comma,
		Encoding:         encoding,
		HasHeader:        hasHeader,
		Comment:          defaults.Comment,
		FieldsPerRecord:  defaults.FieldsPerRecord,
		LazyQuotes:       defaults.LazyQuotes,
		TrimLeadingSpace: defaults.TrimLeadingSpace,
	}, nil
}

// NewCsvReader creates new csv reader instance.
func NewCsvReader(r io.Reader, d *FileDialect) (reader *csv.Reader) {
	// TODO: separate the logic to switch decoder based on encoding
	if d.Encoding == "sjis" {
		decoder := japanese.ShiftJIS.NewDecoder()
		reader = csv.NewReader(transform.NewReader(r, decoder))
	} else {
		reader = csv.NewReader(r)
	}
	reader.Comma = d.Comma
	reader.Comment = d.Comment
	reader.FieldsPerRecord = d.FieldsPerRecord
	return reader
}

// NewCsvWriter creates new csv writer instance.
func NewCsvWriter(w io.Writer, d *FileDialect) (writer *csv.Writer) {
	// TODO: separate the logic to switch encoder based on encoding
	if d.Encoding == "sjis" {
		encoder := japanese.ShiftJIS.NewEncoder()
		writer = csv.NewWriter(transform.NewWriter(w, encoder))
	} else {
		writer = csv.NewWriter(w)
	}
	writer.Comma = d.Comma
	return writer
}
