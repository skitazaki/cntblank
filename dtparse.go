package main

import (
	"time"
)

var TimeLayouts = []string{
	time.ANSIC,
	time.UnixDate,
	time.RubyDate,
	time.RFC822,
	time.RFC822Z,
	time.RFC850,
	time.RFC1123,
	time.RFC1123Z,
	time.RFC3339,
	time.RFC3339Nano,
	time.Kitchen,
	time.Stamp,
	time.StampMilli,
	time.StampMicro,
	time.StampNano,
	"20060102",
	"2006/1/2",
	"2006/1/2 15:4",
	"2006-01-02",
	"2006-01-02 15:4",
	"2006-01-02 15:4:5",
}

func parseDateTime(s string) (t time.Time, err error) {
	for _, layout := range TimeLayouts {
		t, err = time.Parse(layout, s)
		if err == nil {
			return
		}
	}
	return
}
