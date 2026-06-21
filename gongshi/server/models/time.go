package models

import (
	"database/sql/driver"
	"errors"
	"strings"
	"time"
)

const timeFormat = "2006-01-02 15:04:05"

type LocalTime time.Time

func (t LocalTime) MarshalJSON() ([]byte, error) {
	if time.Time(t).IsZero() {
		return []byte(`""`), nil
	}
	b := make([]byte, 0, len(timeFormat)+2)
	b = append(b, '"')
	b = time.Time(t).AppendFormat(b, timeFormat)
	b = append(b, '"')
	return b, nil
}

func (t *LocalTime) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), `"`)
	if s == "" || s == "null" {
		*t = LocalTime(time.Time{})
		return nil
	}
	parsed, err := time.ParseInLocation(timeFormat, s, time.Local)
	if err != nil {
		parsed, err = time.ParseInLocation(time.RFC3339Nano, s, time.Local)
	}
	if err != nil {
		return err
	}
	*t = LocalTime(parsed)
	return nil
}

func (t LocalTime) Value() (driver.Value, error) {
	if time.Time(t).IsZero() {
		return nil, nil
	}
	return time.Time(t).Format(timeFormat), nil
}

func (t *LocalTime) Scan(value interface{}) error {
	if value == nil {
		*t = LocalTime(time.Time{})
		return nil
	}
	var v time.Time
	var err error
	switch val := value.(type) {
	case time.Time:
		v = val
	case string:
		v, err = time.ParseInLocation(timeFormat, val, time.Local)
		if err != nil {
			v, err = time.ParseInLocation(time.RFC3339Nano, val, time.Local)
		}
	case []byte:
		v, err = time.ParseInLocation(timeFormat, string(val), time.Local)
		if err != nil {
			v, err = time.ParseInLocation(time.RFC3339Nano, string(val), time.Local)
		}
	default:
		return errors.New("unsupported LocalTime scan value")
	}
	if err != nil {
		return err
	}
	*t = LocalTime(v)
	return nil
}

func (t LocalTime) String() string {
	return time.Time(t).Format(timeFormat)
}

type LocalTimePtr *LocalTime
