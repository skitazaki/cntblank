package main

import (
	"io"
	"os"
	"unicode/utf8"

	log "github.com/Sirupsen/logrus"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

// Command line options.
var (
	cli             = kingpin.New("cntblank", "Count blank cells on text-based tabular data.")
	cliVerbose      = cli.Flag("verbose", "Set verbose mode on.").Short('v').Bool()
	cliInEncoding   = cli.Flag("input-encoding", "Input encoding.").Short('e').String()
	cliOutEncoding  = cli.Flag("output-encoding", "Output encoding.").Short('E').String()
	cliInDelimiter  = cli.Flag("input-delimiter", "Input field delimiter.").String()
	cliOutDelimiter = cli.Flag("output-delimiter", "Output field delmiter.").String()
	cliNoHeader     = cli.Flag("without-header", "Tabular does not have header line.").Bool()
	cliOutNoHeader  = cli.Flag("output-without-header", "Output report does not have header line.").Bool()
	cliStrict       = cli.Flag("strict", "Check column size strictly.").Bool()
	cliSheet        = cli.Flag("sheet", "Excel sheet number which starts with 1.").Int()
	cliOutMeta      = cli.Flag("output-meta", "Put meta information.").Bool()
	cliOutput       = cli.Flag("output", "Output file.").Short('o').String()
	cliOutFormat    = cli.Flag("output-format", "Output format.").String()
	cliTabularFiles = cli.Arg("tabfile", "Tabular data files.").ExistingFiles()
)

func main() {
	log.SetOutput(os.Stderr)
	cli.Version(VERSION)
	cli.Author(AUTHOR)
	_, err := cli.Parse(os.Args[1:])
	if err != nil {
		log.Fatal(err)
		return
	}
	// Setup logging verbosity.
	if *cliVerbose {
		log.SetLevel(log.DebugLevel)
	}
	// Set output stream.
	var output io.Writer
	if len(*cliOutput) > 0 {
		fp, err := os.Create(*cliOutput)
		if err != nil {
			log.Fatal(err)
			return
		}
		defer fp.Close()
		output = fp
	} else {
		output = os.Stdout
	}
	inDialect, outDialect := populateIODialect()
	// Run main application logic.
	app, err := newApplication(output, *cliOutFormat, outDialect)
	if err != nil {
		log.Fatal(err)
		return
	}
	var files []string
	if len(*cliTabularFiles) > 0 {
		files = *cliTabularFiles
	} else {
		files = append(files, "")
	}
	for _, file := range files {
		err = app.run(file, inDialect)
		if err != nil {
			log.Error(err)
		}
	}
}

func populateIODialect() (inDialect *FileDialect, outDialect *FileDialect) {
	// Convert delimiter type from string to rune. default is TAB.
	inComma := '\t'
	outComma := '\t'
	if len(*cliInDelimiter) > 0 {
		comma, size := utf8.DecodeRuneInString(*cliInDelimiter)
		if size == utf8.RuneError {
			log.Warn("input delimiter option is invalid, but continue running.")
		} else {
			inComma = comma
		}
	}
	if len(*cliOutDelimiter) > 0 {
		comma, size := utf8.DecodeRuneInString(*cliOutDelimiter)
		if size == utf8.RuneError {
			log.Warn("output delimiter option is invalid, but continue running.")
		} else {
			outComma = comma
		}
	}
	// Check encoding options. default is "utf8".
	inEncoding := "utf8"
	outEncoding := "utf8"
	if len(*cliInEncoding) > 0 {
		if *cliInEncoding == "sjis" {
			inEncoding = *cliInEncoding
		} else {
			log.Warn("unknown input encoding: ", *cliInEncoding)
		}
	}
	if len(*cliOutEncoding) > 0 {
		if *cliOutEncoding == "sjis" {
			outEncoding = *cliOutEncoding
		} else {
			log.Warn("unknown output encoding: ", *cliOutEncoding)
		}
	}
	inDialect = &FileDialect{
		Encoding:        inEncoding,
		Comma:           inComma,
		Comment:         '#',
		FieldsPerRecord: -1,
		HasHeader:       !*cliNoHeader,
		SheetNumber:     *cliSheet,
	}
	if *cliStrict {
		inDialect.FieldsPerRecord = 0
	}
	outDialect = &FileDialect{
		Encoding:    outEncoding,
		Comma:       outComma,
		HasHeader:   !*cliOutNoHeader,
		HasMetadata: *cliOutMeta,
	}
	return
}
