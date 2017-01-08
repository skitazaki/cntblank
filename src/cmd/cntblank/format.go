package main

import (
	"strings"
)

// Format represents file format.
type Format int

const (
	// Unknown is nil value of Format type
	Unknown Format = iota
	// CSV is delimiter separated values
	CSV
	// Excel is Microsoft office suite
	Excel
	// JSON is JSON object
	JSON
	// Text is plain text format
	Text
	// HTML is HTML file format
	HTML
)

func (f Format) String() string {
	switch f {
	case Unknown:
		return "Unknown"
	case CSV:
		return "CSV"
	case Excel:
		return "Excel"
	case JSON:
		return "JSON"
	case Text:
		return "Text"
	case HTML:
		return "HTML"
	default:
		return "undefined"
	}
}

func formatFrom(s string) Format {
	// Compare strings in lower case removing initial dot in case of
	// given string is file extention.
	switch strings.ToLower(strings.TrimPrefix(s, ".")) {
	case "csv", "tsv", "txt":
		return CSV
	case "excel", "xlsx":
		return Excel
	case "json":
		return JSON
	case "text":
		return Text
	case "html":
		return HTML
	default:
		return Unknown
	}
}
