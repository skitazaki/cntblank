package csvhelper

import (
	"reflect"
	"strconv"
)

func min(x, y int) int {
	if x > y {
		return y
	}
	return x
}

// Unmarshal unmarshals list of strings into given interface struct.
// see http://stackoverflow.com/questions/20768511/unmarshal-csv-record-into-struct-in-go
func Unmarshal(record []string, v interface{}, strict bool) error {
	s := reflect.ValueOf(v).Elem()
	if strict && s.NumField() != len(record) {
		return &FieldMismatch{s.NumField(), len(record)}
	}
	for i := 0; i < min(s.NumField(), len(record)); i++ {
		f := s.Field(i)
		switch f.Type().String() {
		case "string":
			f.SetString(record[i])
		case "int":
			ival, err := strconv.ParseInt(record[i], 10, 0)
			if err != nil {
				return err
			}
			f.SetInt(ival)
		case "sql.NullString":
			val := ToNullString(record[i])
			f.Set(reflect.ValueOf(val))
		case "sql.NullBool":
			val := ToNullBool(record[i])
			f.Set(reflect.ValueOf(val))
		case "sql.NullInt64":
			val := ToNullInt64(record[i])
			f.Set(reflect.ValueOf(val))
		case "sql.NullFloat64":
			val := ToNullFloat64(record[i])
			f.Set(reflect.ValueOf(val))
		case "null.String":
			val := ToString(record[i])
			f.Set(reflect.ValueOf(val))
		case "null.Bool":
			val := ToBool(record[i])
			f.Set(reflect.ValueOf(val))
		case "null.Int":
			val := ToInt(record[i])
			f.Set(reflect.ValueOf(val))
		case "null.Float":
			val := ToFloat(record[i])
			f.Set(reflect.ValueOf(val))
		default:
			return &UnsupportedType{f.Type().String()}
		}
	}
	return nil
}

type FieldMismatch struct {
	expected, found int
}

func (e *FieldMismatch) Error() string {
	return "CSV line fields mismatch. Expected " + strconv.Itoa(e.expected) + " found " + strconv.Itoa(e.found)
}

type UnsupportedType struct {
	Type string
}

func (e *UnsupportedType) Error() string {
	return "Unsupported type: " + e.Type
}
