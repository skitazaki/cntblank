[![Build Status](https://secure.travis-ci.org/skitazaki/cntblank.png?branch=master)](http://travis-ci.org/skitazaki/cntblank/tree/master)
[![Build Status](https://secure.travis-ci.org/skitazaki/cntblank.png?branch=develop)](http://travis-ci.org/skitazaki/cntblank)

# Count Blank Cells

Count blank cells on text-based tabular data.

Short examples:

It can parse tab-delimited data file, which contains only 2 columns.

```bash
$ ./cntblank --output-meta testdata/prefecture_jp.tsv
INFO[0000] start parsing with 2 columns.                 path=testdata/prefecture_jp.tsv
INFO[0000] finish parsing 48 lines to get 47 records with 2 columns. 0 errors detected.  path=testdata/prefecture_jp.tsv
# File          testdata/prefecture_jp.tsv      prefecture_jp.tsv    4884d04103df0fd8a9e792866ca0b870
# Field         2                               (has header)
# Record        47
seq     Name    #Blank  %Blank  MinLength       MaxLength       #Int    #Float  #Bool   #Time    Minimum Maximum #True   #False
1       都道府県コード  0       0.0000  2       2       47       47               1       47
2       都道府県        0       0.0000  3       4
```

Since it accepts standard input when no file arguments are given,
you can pipe another output such as downloaded contents.

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

It also accepts Microsoft Excel file whose extension is ".xlsx".

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

## Development

Requirements:

- Golang 1.7
- `gb` for build tool

### Setup and library dependency

`Makefile` includes *setup* target calling `gb vendor restore`:

```bash
$ make setup
```

To see libraries:

```bash
$ gb vendor list
```

### Build

*build* target calls `go fmt`, `go vet`, `goimports`, and `gb build`.

```bash
$ make build
```

*test* target calls `gb test`.

```bash
$ make test
```

*local* target runs program against under *testdata/* directory after building binary.

```bash
$ make local
```

To generate binary files for multiple architecture,
simply run `make dist`.

## Changes

v0.7:

- Support Excel writer
- Expand files walking directories.
- Introduce *null* package on dependency to separate *csvhelper* package.
- Move to Go 1.7.

v0.6:

- Support JSON and HTML writer.
- Add `--output-format` option on command line arguments.

v0.5:

- Support Excel reader
- Migrate build tool to `gb`.
- Update encoding package to handle cp932 charset range.

v0.4:

- Change output tabular layout.
- Parse several date formats.
- Accept multiple files on command line arguments.
- Introduce *govalidate* package on dependency.
- Add `Makefile` to setup and build.
- Integrate Travis-CI.

v0.3:

- Change the output format of blank ratio for Excel pasting.
- Trim blank characters.

v0.2:

- Add length and range of values on output report.
- Move to Go 1.5.

v0.1:

- Just first release.


## Memorandum about development

Packages to think about vendoring or reference:

- [ImJasonH/csvstruct: Decode/encode CSV data into/from structs using reflection.](https://github.com/ImJasonH/csvstruct)
- [lukasmartinelli/pgclimb: Export data from PostgreSQL into different data formats](https://github.com/lukasmartinelli/pgclimb)

Tasks:

- [x] Write results as Excel format
- [ ] Parse whole sheets on Excel
- [ ] Guess field data type for SQL CREATE statement (CHAR, VARCHAR, NUMERIC, DATE, TIMESTAMP)