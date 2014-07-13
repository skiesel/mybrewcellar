package models

import (
	"strconv"
	"strings"
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
	time := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, utc)
	return &Date{
		Date: time,
	}
}

func ParseDate(date string) *Date {
	tokens := strings.Split(date, "-")
	if len(tokens) != 3 {
		return nil
	}
	year, err := strconv.Atoi(tokens[0])
	if err != nil {
		return nil
	}
	month, err := strconv.Atoi(tokens[1])
	if err != nil {
		return nil
	}
	day, err := strconv.Atoi(tokens[2])
	if err != nil {
		return nil
	}
	utc, err := time.LoadLocation("UTC")
	if err != nil {
		return nil
	}
	time := time.Date(year, time.Month(month), day, 0, 0, 0, 0, utc)
	return &Date{
		Date: time,
	}
}

func (date *Date) ToString() string {
	const layout = "Jan 2, 2006"
	return date.Date.Format(layout)
}

func (date *Date) ToDSString() string {
	const layout = "2006-01-02"
	return date.Date.Format(layout)
}
