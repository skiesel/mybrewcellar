package models

import (
	"time"
)

type Date struct {
	Date time.Time
}

func Now() *Date {
	now := time.Now()
	utc, err := time.LoadLocation("UTC")
	if err != nil {
		return nil
	}
	t := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, utc)
	return &Date{
		Date: t,
	}
}

func ParseDate(date string) *Date {
	const layout = "Jan 2, 2006"
	t, err := time.Parse(layout, date)
	if err != nil {
		return nil
	}
	utc, err := time.LoadLocation("UTC")
	if err != nil {
		return nil
	}
	t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, utc)
	return &Date{
		Date: t,
	}
}

func (date *Date) ToString() string {
	const layout = "Jan 2, 2006"
	return date.Date.Format(layout)
}
