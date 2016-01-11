[![Build Status](https://secure.travis-ci.org/skitazaki/cntblank.png?branch=master)](http://travis-ci.org/skitazaki/cntblank)

# Count Blank Cells

Count blank cells on text-based tabular data.

Short examples:

```bash
$ ./cntblank --output-meta testdata/prefecture_jp.tsv
INFO[0000] start parsing with 2 columns.                 path=testdata/prefecture_jp.tsv
INFO[0000] finish parsing 48 lines to get 47 records with 2 columns. 0 errors detected.  path=testdata/prefecture_jp.tsv
# File          testdata/prefecture_jp.tsv      prefecture_jp.tsv
# Field         2                               (has header)
# Record        47
seq     Name    #Blank  %Blank  MinLength       MaxLength       #Int    #Float  #Bool   #Time    Minimum Maximum #True   #False
1       都道府県コード  0       0.0000  2       2       47       47               1       47
2       都道府県        0       0.0000  3       4
```

Since file argument is optional, it accepts standard input.

```bash
$ curl -sL https://raw.githubusercontent.com/datasets/language-codes/master/data/ietf-language-tags.csv |
      ./cntblank --input-delimiter=, 2>/dev/null | cut -f1-6 | expand -t 15
seq            Name           #Blank         %Blank         MinLength      MaxLength
1              lang           0              0.0000         2              14
2              langType       0              0.0000         2              4
3              territory      213            0.3096         2              3
4              revGenDate     0              0.0000         10             10
5              defs           0              0.0000         1              1
6              file           0              0.0000         6              18
```

Also accepts Excel file whose extension is ".xlsx".

```bash
$ ./cntblank --output-delimiter=, testdata/addrcode_jp.xlsx
INFO[0000] start parsing with 7 columns.
INFO[0000] finish parsing 1789 lines with 0 errors       path=testdata/addrcode_jp.xlsx
INFO[0000] get 1788 records with 7 columns
seq,Name,#Blank,%Blank,MinLength,MaxLength,#Int,#Float,#Bool,#Time,Minimum,Maximum,#True,#False
1,団体コード,0,0.0000,6,6,1788,1788,,,10006,473821,,
2,都道府県名（漢字）,0,0.0000,3,4,,,,,,,,
3,市区町村名（漢字）,47,0.0263,2,7,,,,,,,,
4,都道府県名（カナ）,0,0.0000,4,7,,,,,,,,
5,市区町村名（カナ）,47,0.0263,2,13,,,,,,,,
6,Column006,1788,1.0000,,,,,,,,,,
7,Column007,1788,1.0000,,,,,,,,,,
```

## Full Usage

`--help` shows the details.

- Default input/output encoding is "UTF-8" and it also accepts only "sjis" value on the option.
- Default input/output delimiter is TAB.
- Meta-information is file path, field length, and number of records.
- If no file path arguments are given, process standard input.

```text
usage: cntblank [<flags>] [<tabfile>...]

Count blank cells on text-based tabular data.

Flags:
      --help            Show context-sensitive help (also try --help-long and --help-man).
  -v, --verbose         Set verbose mode on.
  -e, --input-encoding=INPUT-ENCODING
                        Input encoding.
  -E, --output-encoding=OUTPUT-ENCODING
                        Output encoding.
      --input-delimiter=INPUT-DELIMITER
                        Input field delimiter.
      --output-delimiter=OUTPUT-DELIMITER
                        Output field delmiter.
      --without-header  Tabular does not have header line.
      --strict          Check column size strictly.
      --sheet=SHEET     Excel sheet number which starts with 1.
      --output-meta     Put meta information.
  -o, --output=OUTPUT   Output file.
      --version         Show application version.

Args:
  [<tabfile>]  Tabular data files.
```

### Output fields

| Name | Description |
|------|-------------|
| seq | Sequential number which starts with one. |
| Name | Field name from first header line, otherwise "ColumnNNN" where NNN is sequential number. |
| #Blank | Count of blank cells. |
| %Blank | Percentage of blank cells. |
| MinLength | Minimum length of valid cells. |
| MaxLength | Maximum length of valid cells. |
| #Int | Count of integer type cells. This may be blank. |
| #Float | Count of float type cells. This may be blank. |
| #Bool | Count of bool type cells. This may be blank. |
| #Time | Count of time type cells. This may be blank. |
| Minimum | Minimum value after guessing data type. |
| Maximum | Maximum value after guessing data type. |
| #True | Count of cells which should be treated as boolean true. |
| #False | Count of cells which should be treated as boolean false. |

Note that "1" is interpreted as boolean true and "0" is also interpreted as boolean false.
Therefore, if a column is integer field, "#True" represents the count of "1" and "#False"
represents the count of "0" in the field.
Since some buggy data files sometimes include "0" as null accidentally, this feature may
help you to count up pseudo blank cells.


## Development

- Golang 1.5

### Setup and library dependency

`Makefile` includes setup target calling `go get -d`:

```bash
$ make setup
```

Libraries:

- github.com/Sirupsen/logrus
- gopkg.in/alecthomas/kingpin.v2
- github.com/asaskevich/govalidator
- github.com/tealeg/xlsx

### Build

```bash
$ go build -o cntblank

OR

$ make build
```

`build.sh` is a build script to generate binary files for multiple architecture
using docker container. And, `make` wraps whole processes to generate binaries.
