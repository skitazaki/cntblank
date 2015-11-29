package main

// FileDialect is a configuration for "csv.Reader".
type FileDialect struct {
	Encoding         string // file encoding (utf8 or sjis only)
	Comma            rune   // field delimiter (set to ',' by NewReader)
	Comment          rune   // comment character for start of line
	FieldsPerRecord  int    // number of expected fields per record
	LazyQuotes       bool   // allow lazy quotes
	TrimLeadingSpace bool   // trim leading space
	HasHeader        bool   // CSV file has header line
}
