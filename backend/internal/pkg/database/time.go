package database

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

// Time returns a scanner/valuer that accepts SQLite TEXT timestamps and Go time.Time.
func Time(t *time.Time) interface {
	sql.Scanner
	driver.Valuer
} {
	return (*timeValue)(t)
}

// BindTime formats a time for SQL parameters (RFC3339Nano UTC). Prefer this over raw time.Time
// so SQLite TEXT values stay parseable without relying on Go's String() layout.
func BindTime(t time.Time) driver.Valuer {
	return timeValue(t)
}

// BindNullTime formats an optional time for SQL parameters.
func BindNullTime(t *time.Time) any {
	if t == nil {
		return nil
	}
	return BindTime(*t)
}

// NullTime scans optional timestamps (NULL → nil pointer).
func NullTime(t **time.Time) sql.Scanner {
	return &nullTimeValue{dest: t}
}

type timeValue time.Time

func (t *timeValue) Scan(src any) error {
	if t == nil {
		return fmt.Errorf("database.Time: nil destination")
	}
	parsed, err := parseTimeValue(src)
	if err != nil {
		return err
	}
	if parsed == nil {
		*t = timeValue{}
		return nil
	}
	*t = timeValue(*parsed)
	return nil
}

func (t timeValue) Value() (driver.Value, error) {
	tt := time.Time(t)
	if tt.IsZero() {
		return nil, nil
	}
	return tt.UTC().Format(time.RFC3339Nano), nil
}

type nullTimeValue struct {
	dest **time.Time
}

func (n *nullTimeValue) Scan(src any) error {
	if n == nil || n.dest == nil {
		return fmt.Errorf("database.NullTime: nil destination")
	}
	parsed, err := parseTimeValue(src)
	if err != nil {
		return err
	}
	*n.dest = parsed
	return nil
}

func parseTimeValue(src any) (*time.Time, error) {
	switch v := src.(type) {
	case nil:
		return nil, nil
	case time.Time:
		t := v.UTC()
		return &t, nil
	case string:
		return parseTimeString(v)
	case []byte:
		return parseTimeString(string(v))
	default:
		return nil, fmt.Errorf("database: cannot scan %T into time.Time", src)
	}
}

func parseTimeString(s string) (*time.Time, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, nil
	}
	// Strip Go monotonic clock suffix from time.Time.String().
	if i := strings.Index(s, " m="); i >= 0 {
		s = s[:i]
	}
	formats := []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02 15:04:05.999999999-07:00",
		"2006-01-02 15:04:05.999999999 -0700 MST",
		"2006-01-02 15:04:05 -0700 MST",
		"2006-01-02 15:04:05.999999999",
		"2006-01-02 15:04:05",
		"2006-01-02",
	}
	for _, layout := range formats {
		if parsed, err := time.Parse(layout, s); err == nil {
			t := parsed.UTC()
			return &t, nil
		}
	}
	return nil, fmt.Errorf("database: cannot parse time %q", s)
}
