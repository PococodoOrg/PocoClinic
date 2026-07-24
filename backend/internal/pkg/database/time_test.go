package database

import (
	"testing"
	"time"
)

func TestTimeScan_GoStringLayout(t *testing.T) {
	var got time.Time
	raw := "2026-07-23 11:08:20.1918763 -0500 CDT m=+0.126976501"
	if err := Time(&got).Scan(raw); err != nil {
		t.Fatalf("scan: %v", err)
	}
	if got.Year() != 2026 || got.Month() != time.July {
		t.Fatalf("unexpected: %v", got)
	}
}

func TestTimeScan_RFC3339Text(t *testing.T) {
	var got time.Time
	raw := "2026-07-22T16:04:05.123456789Z"
	if err := Time(&got).Scan(raw); err != nil {
		t.Fatalf("scan: %v", err)
	}
	if got.Year() != 2026 || got.Month() != time.July || got.Day() != 22 {
		t.Fatalf("unexpected time: %v", got)
	}
}

func TestTimeScan_DateOnly(t *testing.T) {
	var got time.Time
	if err := Time(&got).Scan("1990-01-15"); err != nil {
		t.Fatalf("scan: %v", err)
	}
	if got.Format("2006-01-02") != "1990-01-15" {
		t.Fatalf("got %v", got)
	}
}

func TestNullTimeScan_NullAndText(t *testing.T) {
	var dest *time.Time
	if err := NullTime(&dest).Scan(nil); err != nil {
		t.Fatalf("null: %v", err)
	}
	if dest != nil {
		t.Fatalf("expected nil, got %v", dest)
	}
	if err := NullTime(&dest).Scan("2026-01-02T03:04:05Z"); err != nil {
		t.Fatalf("text: %v", err)
	}
	if dest == nil || dest.Year() != 2026 {
		t.Fatalf("dest %#v", dest)
	}
}

func TestTimeValue_FormatsUTC(t *testing.T) {
	in := time.Date(2026, 7, 22, 12, 0, 0, 0, time.UTC)
	v, err := (*timeValue)(&in).Value()
	if err != nil {
		t.Fatalf("value: %v", err)
	}
	s, ok := v.(string)
	if !ok || s != "2026-07-22T12:00:00Z" {
		t.Fatalf("value %#v", v)
	}
}
