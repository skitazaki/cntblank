package main

import (
	"io"
	"os"
	"unicode/utf8"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
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
	cliStrict       = cli.Flag("strict", "Check column size strictly.").Bool()
	cliOutMeta      = cli.Flag("output-meta", "Put meta information.").Bool()
	cliOutput       = cli.Flag("output", "Output file.").Short('o').String()
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
	// Convert delimiter type from string to rune.
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
	// Check encoding options.
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
	// Run main application logic.
	app, err := newApplication(output, outEncoding, outComma, *cliOutMeta)
	if err != nil {
		log.Fatal(err)
		return
	}
	if len(*cliTabularFiles) > 0 {
		for _, file := range *cliTabularFiles {
			app.run(file, inEncoding, inComma, *cliNoHeader, *cliStrict)
		}
	} else {
		app.run("", inEncoding, inComma, *cliNoHeader, *cliStrict)
	}
}
