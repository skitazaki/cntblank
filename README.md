# Count Blank Cells

Count blank cells on text-based tabular data.

Short examples:

```bash
$ ./cntblank --output-meta examples/prefecture_jp.tsv
INFO[0000] start parsing with 2 columns.                 path=examples/prefecture_jp.tsv
INFO[0000] finish parsing 48 lines to get 47 records with 2 columns. 0 errors detected.  path=examples/prefecture_jp.tsv
# File          examples/prefecture_jp.tsv      prefecture_jp.tsv
# Field         2                               (has header)
# Record        47
seq     Name    #Blank  %Blank  MinLength       MaxLength       #Int    #Float  #Bool   Minimum Maximum MinimumF        MaximumF        #True   #False
1       都道府県コード  0       0.0000  2       2       47                      1       47
2       都道府県        0       0.0000  3       4
```

Since file argument is optional, it accepts standard input.

```bash
$ curl -sL https://raw.githubusercontent.com/datasets/language-codes/master/data/ietf-language-tags.csv |
      ./cntblank --input-delimiter=, 2>/dev/null | cut -f1-6 | expand -t 15
seq            Name           #Blank         %Blank         MinLength      MaxLength
1              lang           0              0.0000         2              14
2              langType       0              0.0000         2              4
3              territory      213            0.3096        2              3
4              revGenDate     0              0.0000         10             10
5              defs           0              0.0000         1              1
6              file           0              0.0000         6              18
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
| Minimum | Minimum integer value. |
| Maximum | Maximum integer value. |
| MinimumF | Minimum float value. |
| MaximumF | Maximum float value. |
| #True | Count of cells which should be treated as boolean true. |
| #False | Count of cells which should be treated as boolean false. |

Since type detecton is done order by `int`, `float`, and `bool`, there is no relation with minimum/maximum integer value and minimum/maximum float value.
And, "1" or "0" is not treated as boolean because it is treated as integer in the first detection logic.


## Development setup

- Golang 1.5

### Library Dependency

- github.com/Sirupsen/logrus
- gopkg.in/alecthomas/kingpin.v2

### Build

```bash
$ go build -o cntblank
```

`build.sh` is a build script to generate binary files for multiple architecture
using docker container.
