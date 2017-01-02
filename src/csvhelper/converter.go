package csvhelper

import (
	"database/sql"
	"strconv"

	null "gopkg.in/guregu/null.v3"
)

// ToNullString is a helper function
func ToNullString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: s != ""}
}

// ToNullInt64 is a helper function
func ToNullInt64(s string) sql.NullInt64 {
	i, err := strconv.Atoi(s)
	return sql.NullInt64{Int64: int64(i), Valid: err == nil}
}

// ToNullFloat64 is a helper function
func ToNullFloat64(s string) sql.NullFloat64 {
	f, err := strconv.ParseFloat(s, 64)
	return sql.NullFloat64{Float64: f, Valid: err == nil}
}

// ToNullBool is a helper function
func ToNullBool(s string) sql.NullBool {
	b, err := strconv.ParseBool(s)
	return sql.NullBool{Bool: b, Valid: err == nil}
}

// ToString is a helper function
func ToString(s string) null.String {
	return null.NewString(s, s != "")
}

// ToInt is a helper function
func ToInt(s string) null.Int {
	i, err := strconv.Atoi(s)
	return null.NewInt(int64(i), err == nil)
}

// ToFloat is a helper function
func ToFloat(s string) null.Float {
	f, err := strconv.ParseFloat(s, 64)
	return null.NewFloat(f, err == nil)
}

// ToBool is a helper function
func ToBool(s string) null.Bool {
	b, err := strconv.ParseBool(s)
	return null.NewBool(b, err == nil)
}
