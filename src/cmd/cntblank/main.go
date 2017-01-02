package main

import (
	"io"
	"os"

	log "github.com/Sirupsen/logrus"
	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"csvhelper"
)

// Command line options.
var (
	cli             = kingpin.New("cntblank", "Count blank cells on text-based tabular data.")
	cliVerbose      = cli.Flag("verbose", "Set verbose mode on.").Short('v').Bool()
	cliInEncoding   = cli.Flag("input-encoding", "Input encoding.").Short('e').Default("utf8").String()
	cliOutEncoding  = cli.Flag("output-encoding", "Output encoding.").Short('E').Default("utf8").String()
	cliInDelimiter  = cli.Flag("input-delimiter", "Input field delimiter.").Default("\t").String()
	cliOutDelimiter = cli.Flag("output-delimiter", "Output field delmiter.").Default("\t").String()
	cliNoHeader     = cli.Flag("without-header", "Tabular does not have header line.").Bool()
	cliOutNoHeader  = cli.Flag("output-without-header", "Output report does not have header line.").Bool()
	cliStrict       = cli.Flag("strict", "Check column size strictly.").Bool()
	cliSheet        = cli.Flag("sheet", "Excel sheet number which starts with 1.").Int()
	cliRecursive    = cli.Flag("recursive", "Traverse directory recursively.").Short('r').Bool()
	cliOutMeta      = cli.Flag("output-meta", "Put meta information.").Bool()
	cliOutput       = cli.Flag("output", "Output file.").Short('o').String()
	cliOutFormat    = cli.Flag("output-format", "Output format.").String()
	cliTabularFiles = cli.Arg("tabfile", "Tabular data files.").Strings()
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
	app, err := newApplication(*cliRecursive, output, *cliOutFormat, outDialect)
	if err != nil {
		log.Fatal(err)
		return
	}
	files := *cliTabularFiles
	err = app.Run(files, inDialect)
	if err != nil {
		log.Error(err)
	}
}

func populateIODialect() (inDialect *csvhelper.FileDialect, outDialect *csvhelper.FileDialect) {
	inDialect, err := csvhelper.NewFileDialect(*cliInDelimiter, *cliInEncoding, !*cliNoHeader)
	if err != nil {
		// TODO: report error.
	}
	inDialect.SheetNumber = *cliSheet
	if *cliStrict {
		inDialect.FieldsPerRecord = 0
	}
	outDialect, err = csvhelper.NewFileDialect(*cliOutDelimiter, *cliOutEncoding, !*cliOutNoHeader)
	if err != nil {
		// TODO: report error.
	}
	outDialect.HasMetadata = *cliOutMeta
	return
}
