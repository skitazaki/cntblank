# Count Blank Cells

Count blank cells on text-based tabular data.

Short examples:

```bash
$ cntblank --output-meta examples/prefecture_jp.tsv
INFO[0000] start parsing with 2 columns header.          path=examples/prefecture_jp.tsv
INFO[0000] finish parsing 47 records with 2 columns. 0 errors detected.  path=examples/prefecture_jp.tsv
# File  examples/prefecture_jp.tsv
# Field 2
# Record        47
seq     Name    #Blank  %Blank
1       都道府県コード  0       0.0000
2       都道府県        0       0.0000
```

Since file argument is optional, it accepts standard input.

```bash
$ curl -sL https://raw.githubusercontent.com/datasets/language-codes/master/data/ietf-language-tags.csv | ./cntblank --input-delimiter=,
INFO[0005] start parsing with 6 columns header.
INFO[0005] finish parsing 688 records with 6 columns. 0 errors detected.
seq     Name    #Blank  %Blank
1       lang    0       0.0000
2       langType        0       0.0000
3       territory       213     30.9593
4       revGenDate      0       0.0000
5       defs    0       0.0000
6       file    0       0.0000
```

## Full Usage

`--help` shows the details.

- Default input/output encoding is "UTF-8" and it also accepts only "sjis" value on the option.
- Default input/output delimiter is TAB.
- Meta-information is file path, field length, and number of records.

```text
usage: cntblank [<flags>] [<tabfile>]

Count blank cells on text-based tabular data.

Flags:
  --help               Show help (also see --help-long and --help-man).
  -v, --verbose        Set verbose mode on.
  -e, --input-encoding=INPUT-ENCODING
                       Input encoding.
  -E, --output-encoding=OUTPUT-ENCODING
                       Output encoding.
  --input-delimiter=INPUT-DELIMITER
                       Input field delimiter.
  --output-delimiter=OUTPUT-DELIMITER
                       Output field delmiter.
  --without-header     Tabular does not have header line.
  --strict             Check column size strictly.
  --output-meta        Put meta information.
  -o, --output=OUTPUT  Output file.

Args:
  [<tabfile>]  Tabular data file.
```


## Development setup

- Golang 1.4

### Library Dependency

- github.com/Sirupsen/logrus
- gopkg.in/alecthomas/kingpin.v2

### Build

```bash
$ go build -o cntblank
```

`build.sh` is a build script to generate binary files for multiple architecture
using docker container.

