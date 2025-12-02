package tables

import (
	"database/sql"
	"strconv"
	"time"
)

func ToString(v interface{}) (string, bool) {
	switch t := v.(type) {
	case string:
		return t, true
	default:
		return "", false
	}
}

func ToFloat64(v interface{}) (float64, bool) {
	switch t := v.(type) {
	case float64:
		return t, true
	case float32:
		return float64(t), true
	case int:
		return float64(t), true
	case int64:
		return float64(t), true
	case int32:
		return float64(t), true
	default:
		return 0, false
	}
}

func ToTime(v interface{}) (time.Time, bool) {
	switch t := v.(type) {
	case time.Time:
		return t, true
	case string:
		if ts, err := time.Parse(time.RFC3339Nano, t); err == nil {
			return ts, true
		}
		if ts, err := time.Parse(time.RFC3339, t); err == nil {
			return ts, true
		}
		if ts, err := time.Parse("2006-01-02", t); err == nil {
			return ts, true
		}
		return time.Time{}, false
	default:
		return time.Time{}, false
	}
}

func VString(v sql.NullString) string {
	if v.Valid {
		return v.String
	}
	return ""
}

func VTimeString(v sql.NullString) string {
	if !v.Valid || v.String == "" {
		return ""
	}
	s := v.String
	if ts, err := time.Parse(time.RFC3339Nano, s); err == nil {
		return ts.UTC().Format(time.RFC3339Nano)
	}
	if ts, err := time.Parse(time.RFC3339, s); err == nil {
		return ts.UTC().Format(time.RFC3339Nano)
	}
	if secs, err := strconv.ParseInt(s, 10, 64); err == nil {
		return time.Unix(secs, 0).UTC().Format(time.RFC3339Nano)
	}
	if ts, err := time.Parse("2006-01-02", s); err == nil {
		return ts.UTC().Format(time.RFC3339Nano)
	}
	return s
}

func VFloat(v sql.NullFloat64) float64 {
	if v.Valid {
		return v.Float64
	}
	return 0
}

// MoscowDateMidnightUTC converts YYYY-MM-DD (Moscow local date) to the UTC instant of Moscow midnight
func MoscowDateMidnightUTC(s string) time.Time {
	if s == "" {
		return time.Time{}
	}
	ts, err := time.Parse("2006-01-02", s)
	if err != nil {
		return time.Time{}
	}
	return ts.Add(-3 * time.Hour).UTC()
}
