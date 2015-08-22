package main

import (
	"io"
	"os"

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
	cliTabularFile  = cli.Arg("tabfile", "Tabular data file.").ExistingFile()
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
	// Run main application logic.
	app, err := newApplication(output, *cliOutEncoding, *cliOutDelimiter, *cliOutMeta)
	if err != nil {
		log.Fatal(err)
		return
	}
	app.run(*cliTabularFile, *cliInEncoding, *cliInDelimiter, *cliNoHeader, *cliStrict)
}
