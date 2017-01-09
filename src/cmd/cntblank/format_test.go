package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatString(t *testing.T) {
	a := assert.New(t)
	for _, tc := range []struct {
		format Format
		want   string
	}{
		{CSV, "CSV"},
		{Excel, "Excel"},
		{JSON, "JSON"},
		{Text, "Text"},
		{HTML, "HTML"},
	} {
		a.Equal(tc.want, tc.format.String())
	}
}

func TestFormatFrom(t *testing.T) {
	a := assert.New(t)
	for _, tc := range []struct {
		str  string
		want Format
	}{
		{"CSV", CSV},
		{"csv", CSV},
		{".tsv", CSV},
		{".txt", CSV},
		{"Excel", Excel},
		{"xlsx", Excel},
		{".xlsx", Excel},
		{"JSON", JSON},
		{"Text", Text},
		{"HTML", HTML},
		{".xls", Unknown},
	} {
		a.Equal(tc.want, formatFrom(tc.str))
	}
}
