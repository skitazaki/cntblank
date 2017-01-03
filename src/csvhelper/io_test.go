package csvhelper

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFileDialect_Default(t *testing.T) {
	d, err := NewFileDialect("", "", false)
	require.Nil(t, err, "NewFileDialect returns error: %v", err)
	assert.Equal(t, ',', d.Comma)
	assert.Equal(t, "", d.Encoding)
	assert.Equal(t, false, d.HasHeader)
	assert.Equal(t, '#', d.Comment)
	assert.Equal(t, -1, d.FieldsPerRecord)
	assert.Equal(t, true, d.LazyQuotes)
	assert.Equal(t, true, d.TrimLeadingSpace)
	assert.Equal(t, 0, d.SheetNumber, "SheetNumber should not be initialized.")
	assert.Equal(t, false, d.HasMetadata)
}

func TestNewFileDialect_Comma(t *testing.T) {
	for i, tc := range []struct {
		expected  rune
		delimiter string
	}{
		{'\t', "\t"},
		{':', ":"},
		{':', "::"}, // Use first 1-letter as rune
	} {
		d, err := NewFileDialect(tc.delimiter, "", false)
		require.Nil(t, err, "NewFileDialect returns error: %v", err)
		assert.Equal(t, tc.expected, d.Comma, "for loop index %d", i)
	}
}

func TestNewFileDialect_Encoding(t *testing.T) {
	for i, tc := range []struct {
		expected string
		encoding string
	}{
		{"", ""},
		{"utf8", "utf8"},
		{"sjis", "sjis"},
		{"cp932", "cp932"},
	} {
		d, err := NewFileDialect("", tc.encoding, false)
		require.Nil(t, err, "NewFileDialect returns error: %v", err)
		assert.Equal(t, tc.expected, d.Encoding, "for loop index %d", i)
	}
}

func TestNewFileDialect_Header(t *testing.T) {
	for i, tc := range []struct {
		expected bool
		header   bool
	}{
		{true, true},
		{false, false},
	} {
		d, err := NewFileDialect("", "", tc.header)
		require.Nil(t, err, "NewFileDialect returns error: %v", err)
		assert.Equal(t, tc.expected, d.HasHeader, "for loop index %d", i)
	}
}
